package partner

import (
	"github.com/cohix/simplcrypto"
	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/update"
)

// RunWithServer is analagous to Start(), but using a gRPC server instead of a client
func (m *Manager) RunWithServer(server PartnerService_StreamUpdatesServer) error {
	log.LogInfo("PartnerManager RunWithServer")

	log.LogInfo("waiting for session validation from partner")
	if err := m.receiveAndCheckSessionFromPartner(server); err != nil {
		return errors.Wrap(err, "PartnerManager failed to receiveAndCheckSessionFromPartner")
	}

	log.LogInfo("partner session validated")
	log.LogInfo("sending data key to partner")
	if err := m.generateAndSendDataKeyToPartner(m.partner, server); err != nil {
		return errors.Wrap(err, "PartnerManager failed to generateAndSendDataKeyToPartner")
	}

	sendChan, recvChan := m.serverSendRecvChans(m.partner, server)

	m.partner.HealthChecker = newHealthChecker()
	go m.partner.HealthChecker.startHealthCheckingWithServer(server)

	log.LogInfo("update stream starting")
	return m.streamUpdates(sendChan, recvChan, m.partner.HealthChecker.UnhealthyChan)
}

// serverSendRecvChans creates a channel to send updates over, and then spawns a goroutine to
// constantly receive from the update server, decrypt the updates, and send them down the chan
func (m *Manager) serverSendRecvChans(partner *Partner, server PartnerService_StreamUpdatesServer) (chan update.PartnerUpdate, chan update.PartnerUpdate) {
	sendChan := make(chan update.PartnerUpdate)
	recvChan := make(chan update.PartnerUpdate)

	go func(server PartnerService_StreamUpdatesServer, sendChan chan update.PartnerUpdate, recvChan chan update.PartnerUpdate) {
		for {
			select {
			case update := <-sendChan:
				updateReq, err := m.encryptAndSignUpdateForPartner(partner, &update)
				if err != nil {
					log.LogWarn(errors.Wrap(err, "clientSendRecvChans failed to encryptAndSignUpdateForPartner").Error())
					continue
				}

				if err := server.Send(updateReq); err != nil {
					log.LogWarn(errors.Wrap(err, "streamClientRecvChan failed to Recv, terminating partner stream...").Error())
					break
				}
			default:
				updateReq, err := server.Recv()
				if err != nil {
					log.LogWarn(errors.Wrap(err, "streamClientRecvChan failed to Recv, terminating partner stream...").Error())
					break
				}

				update, err := m.decryptAndVerifyUpdateFromPartner(partner, updateReq)
				if err != nil {
					log.LogWarn(errors.Wrap(err, "failed to decryptAndVerifyUpdateFromPartner").Error())
					continue
				} else if update == nil {
					log.LogInfo("received health check, discarding...")
					continue
				}

				recvChan <- *update
			}
		}
	}(server, sendChan, recvChan)

	return sendChan, recvChan
}

func (m *Manager) receiveAndCheckSessionFromPartner(server PartnerService_StreamUpdatesServer) error {
	updateReq, err := server.Recv()
	if err != nil {
		return errors.Wrap(err, "failed to Recv data key update")
	}

	if updateReq.Session == nil {
		return errors.New("session check session missing")
	}

	return m.Auth.CheckAuth(updateReq.Session)
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
