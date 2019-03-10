package partner

import (
	"github.com/cohix/simplcrypto"
	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
)

// RunWithServer is analagous to Start(), but using a gRPC server instead of a client
func (m *Manager) RunWithServer(server PartnerService_StreamUpdatesServer) error {
	log.LogInfo("waiting for session validation from partner")
	if err := m.receiveAndCheckSessionFromPartner(server); err != nil {
		return errors.Wrap(err, "PartnerManager failed to receiveAndCheckSessionFromPartner")
	}

	log.LogInfo("partner session validated")
	log.LogInfo("sending data key to partner")
	if err := m.generateAndSendDataKeyToPartner(m.partner, server); err != nil {
		return errors.Wrap(err, "PartnerManager failed to generateAndSendDataKeyToPartner")
	}

	return m.syncWithPartner(m.partner, server)
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
