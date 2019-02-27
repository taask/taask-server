package partner

import (
	"github.com/cohix/simplcrypto"
	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
)

// RunWithServer is analagous to Start(), but using a gRPC server instead of a client
func (m *Manager) RunWithServer(server PartnerService_StreamUpdatesServer) error {
	log.LogInfo("PartnerManager RunWithServer")

	if err := m.generateAndSendDataKeyToPartner(m.partner, server); err != nil {
		return errors.Wrap(err, "PartnerManager failed to generateAndSendDataKeyToPartner")
	}

	recvChan := m.streamServerRecvChan(m.partner, server)

	m.partner.HealthChecker = newHealthChecker()
	m.partner.HealthChecker.startHealthCheckingWithServer(server)

	return m.streamUpdates(recvChan, m.partner.HealthChecker.UnhealthyChan)
}

func (m *Manager) streamServerRecvChan(partner *Partner, server PartnerService_StreamUpdatesServer) chan *Update {
	recvChan := make(chan *Update)

	go func(server PartnerService_StreamUpdatesServer, recvChan chan *Update) {
		for {
			updateReq, err := server.Recv()
			if err != nil {
				log.LogWarn(errors.Wrap(err, "streamServerRecvChan failed to Recv, terminating partner stream...").Error())
				break
			}

			update, err := m.decryptAndVerifyUpdateFromPartner(partner, updateReq)
			if err != nil {
				log.LogWarn(errors.Wrap(err, "failed to decryptAndVerifyUpdateFromPartner").Error())
				continue
			}

			recvChan <- update
		}
	}(server, recvChan)

	return recvChan
}

func (m *Manager) generateAndSendDataKeyToPartner(partner *Partner, server PartnerService_StreamUpdatesServer) error {
	newKey, err := simplcrypto.GenerateSymKey()
	if err != nil {
		return errors.Wrap(err, "failed to GenerateSymKey")
	}

	partner.DataKey = newKey

	encKey, err := m.Auth.EncryptForMember(partner.UUID, newKey.JSON())
	if err != nil {
		return errors.Wrap(err, "failed to Encrypt data key")
	}

	keySig, err := m.masterKeypair.Sign(newKey.JSON())
	if err != nil {
		return errors.Wrap(err, "failed to Sign data key")
	}

	updateReq := &UpdateRequest{
		PartnerUUID:     m.UUID,
		EncDataKey:      encKey,
		UpdateSignature: keySig,
	}

	if err := server.Send(updateReq); err != nil {
		return errors.Wrap(err, "failed to Send data key update")
	}

	return nil
}
