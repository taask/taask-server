package partner

import (
	"github.com/cohix/simplcrypto"
	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/service"
)

// StartWithServer is analagous to Start(), but using a gRPC server instead of a client
func (m *Manager) StartWithServer(server service.PartnerService_StreamUpdatesServer) {
	for {
		if err := m.generateAndSendDataKeyToPartner(m.partner, server); err != nil {
			log.LogWarn(errors.Wrap(err, "PartnerManager failed to generateAndSendDataKeyToPartner, will retry in 5s").Error())
			continue
		}

		recvChan := streamServerRecvChan(m.partner, server)

		if err := m.streamUpdates(recvChan); err != nil {
			log.LogWarn(errors.Wrap(err, "PartnerManager encountered streamUpdatesError, will retry in 5s").Error())
		}
	}
}

func streamServerRecvChan(partner *Partner, server service.PartnerService_StreamUpdatesServer) chan *Update {
	recvChan := make(chan *Update)

	go func(server service.PartnerService_StreamUpdatesServer, recvChan chan *Update) {
		for {
			updateReq, err := server.Recv()
			if err != nil {
				log.LogWarn(errors.Wrap(err, "streamServerRecvChan failed to Recv, terminating partner stream...").Error())
				break
			}

			update, err := decryptAndVerifyUpdateFromPartner(partner, updateReq)
			if err != nil {
				log.LogWarn(errors.Wrap(err, "failed to decryptAndVerifyUpdateFromPartner").Error())
				continue
			}

			recvChan <- update
		}
	}(server, recvChan)

	return recvChan
}

func (m *Manager) generateAndSendDataKeyToPartner(partner *Partner, server service.PartnerService_StreamUpdatesServer) error {
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

	updateReq := &service.UpdateRequest{
		PartnerUUID:     m.UUID,
		EncDataKey:      encKey,
		UpdateSignature: keySig,
	}

	if err := server.Send(updateReq); err != nil {
		return errors.Wrap(err, "failed to Send data key update")
	}

	return nil
}
