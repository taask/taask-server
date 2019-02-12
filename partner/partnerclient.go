package partner

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cohix/simplcrypto"
	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	"google.golang.org/grpc"
)

// Start starts the partner manager
func (m *Manager) Start() {
	for {
		if m.partner.HealthChecker != nil {
			if m.partner.HealthChecker.IsHealthy {
				if err := <-m.partner.HealthChecker.UnhealthyChan; err != nil {
					log.LogError(errors.Wrap(err, "partner healthChecker reported unhealthy, will retry in 5s"))
				}
			}
		}

		<-time.After(time.Duration(time.Second * 5))

		if err := m.partner.initClient(); err != nil {
			log.LogWarn(errors.Wrap(err, "PartnerManager failed to initClient, will retry in 5s").Error())
			continue
		}

		if err := m.authenticatePartner(m.partner); err != nil {
			log.LogWarn(errors.Wrap(err, "PartnerManager failed to authenticatePartner, will retry in 5s").Error())
			continue
		}

		client, err := m.partner.Client.StreamUpdates(context.Background(), nil)
		if err != nil {
			log.LogWarn(errors.Wrap(err, "PartnerManager failed to StreamUpdates, will retry in 5s").Error())
			continue
		}

		if err := m.receiveDataKeyFromPartner(m.partner, client); err != nil {
			log.LogWarn(errors.Wrap(err, "PartnerManager failed to receiveDataKeyFromPartner, will retry in 5s").Error())
			continue
		}

		recvChan := m.streamClientRecvChan(m.partner, client)

		m.partner.HealthChecker = newHealthChecker()
		m.partner.HealthChecker.startHealthCheckingWithClient(client)

		if err := m.streamUpdates(recvChan, m.partner.HealthChecker.UnhealthyChan); err != nil {
			log.LogWarn(errors.Wrap(err, "PartnerManager encountered streamUpdatesError, will retry in 5s").Error())
		}
	}
}

func (m *Manager) streamClientRecvChan(partner *Partner, client PartnerService_StreamUpdatesClient) chan *Update {
	recvChan := make(chan *Update)

	go func(client PartnerService_StreamUpdatesClient, recvChan chan *Update) {
		for {
			updateReq, err := client.Recv()
			if err != nil {
				log.LogWarn(errors.Wrap(err, "streamClientRecvChan failed to Recv, terminating partner stream...").Error())
				break
			}

			update, err := m.decryptAndVerifyUpdateFromPartner(partner, updateReq)
			if err != nil {
				log.LogWarn(errors.Wrap(err, "failed to decryptAndVerifyUpdateFromPartner").Error())
				continue
			}

			recvChan <- update
		}
	}(client, recvChan)

	return recvChan
}

func (p *Partner) initClient() error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", p.host, p.port), grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to Dial")
	}

	client := NewPartnerServiceClient(conn)

	p.Client = client

	return nil
}

func (m *Manager) authenticatePartner(partner *Partner) error {
	keypair, err := simplcrypto.GenerateNewKeyPair()
	if err != nil {
		return errors.Wrap(err, "failed to GenerateNewKeyPair")
	}

	timestamp := time.Now().Unix()

	nonce := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonce, uint64(timestamp))
	hashWithNonce := append(m.config.MemberGroup.AuthHash, nonce...)

	authHashSig, err := keypair.Sign(hashWithNonce)
	if err != nil {
		return errors.Wrap(err, "failed to Sign")
	}

	attempt := &auth.Attempt{
		MemberUUID:        m.UUID,
		GroupUUID:         m.config.MemberGroup.UUID,
		PubKey:            keypair.SerializablePubKey(),
		AuthHashSignature: authHashSig,
		Timestamp:         timestamp,
	}

	authResp, err := m.partner.Client.AuthPartner(context.Background(), attempt)
	if err != nil {
		return errors.Wrap(err, "failed to AuthClient")
	}

	challengeBytes, err := keypair.Decrypt(authResp.EncChallenge)
	if err != nil {
		return errors.Wrap(err, "failed to Decrypt challenge")
	}

	masterRunnerPubKey, err := simplcrypto.KeyPairFromSerializedPubKey(authResp.MasterPubKey)
	if err != nil {
		return errors.Wrap(err, "failed to KeyPairFromSerializablePubKey")
	}

	challengeSig, err := keypair.Sign(challengeBytes)
	if err != nil {
		return errors.Wrap(err, "failed to Sign challenge")
	}

	session := &activeSession{
		Session: &auth.Session{
			MemberUUID:          m.UUID,
			GroupUUID:           m.config.MemberGroup.UUID,
			SessionChallengeSig: challengeSig,
		},
		Keypair:      keypair,
		MasterPubKey: masterRunnerPubKey,
	}

	m.partner.ActiveSession = session

	return nil
}

func (m *Manager) receiveDataKeyFromPartner(partner *Partner, client PartnerService_StreamUpdatesClient) error {
	updateReq, err := client.Recv()
	if err != nil {
		return errors.Wrap(err, "failed to Recv data key update")
	}

	if updateReq.UpdateSignature == nil {
		return errors.New("data key update signature missing")
	}

	if updateReq.EncDataKey == nil {
		return errors.New("data key update encDataKey is empty")
	}

	dataKeyJSON, err := partner.ActiveSession.Keypair.Decrypt(updateReq.EncDataKey)
	if err != nil {
		return errors.Wrap(err, "failed to Decrypt data key update")
	}

	if err := partner.ActiveSession.MasterPubKey.Verify(dataKeyJSON, updateReq.UpdateSignature); err != nil {
		return errors.Wrap(err, "failed to Verify data key update")
	}

	dataKey := simplcrypto.SymKey{}
	if err := json.Unmarshal(dataKeyJSON, &dataKey); err != nil {
		return errors.Wrap(err, "failed to Unmarshal data key")
	}

	partner.UUID = updateReq.PartnerUUID
	partner.DataKey = &dataKey

	return nil
}
