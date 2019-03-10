package partner

import (
	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/update"
)

type syncSource interface {
	Send(*UpdateRequest) error
	Recv() (*UpdateRequest, error)
}

func (m *Manager) syncWithPartner(partner *Partner, syncer syncSource) error {
	log.LogInfo("starting partner syncer")
	sendChan, recvChan := m.startSyncerForPartner(m.partner, syncer)

	m.partner.HealthChecker = newHealthChecker()
	go m.partner.HealthChecker.startHealthChecking(syncer)

	log.LogInfo("ready to handle partner updates")
	return m.handleUpdates(sendChan, recvChan, m.partner.HealthChecker.UnhealthyChan)
}

// clientSendRecvChans creates channels to send and receive from the partner, and continuously reads and writes
func (m *Manager) startSyncerForPartner(partner *Partner, syncer syncSource) (chan update.PartnerUpdate, chan update.PartnerUpdate) {
	// buffer updates so that we're not the bottleneck for anything
	sendChan := make(chan update.PartnerUpdate, 64)
	recvChan := make(chan update.PartnerUpdate, 64)

	go func(syncer syncSource, sendChan chan update.PartnerUpdate, recvChan chan update.PartnerUpdate) {
		incomingChan := make(chan *UpdateRequest, 64)
		incomingErrChan := syncerStartReceiving(syncer, incomingChan)

		for {
			shouldTerminate := false

			select {
			case update := <-sendChan:
				updateReq, err := m.encryptAndSignUpdateForPartner(partner, &update)
				if err != nil {
					log.LogWarn(errors.Wrap(err, "clientSendRecvChans failed to encryptAndSignUpdateForPartner").Error())
					continue
				}

				if err := syncer.Send(updateReq); err != nil {
					log.LogWarn(errors.Wrap(err, "streamClientRecvChan failed to Recv, terminating partner stream...").Error())
					shouldTerminate = true
				}
			case updateReq := <-incomingChan:
				update, err := m.decryptAndVerifyUpdateFromPartner(partner, updateReq)
				if err != nil {
					log.LogWarn(errors.Wrap(err, "failed to decryptAndVerifyUpdateFromPartner").Error())
					continue
				} else if update == nil {
					log.LogInfo("received health check, discarding...")
					continue
				}

				recvChan <- *update
			case err := <-incomingErrChan:
				log.LogWarn(errors.Wrap(err, "clientSendRecvChans received incoming recv error, terminating...").Error())
				shouldTerminate = true
			}

			if shouldTerminate {
				break
			}
		}
	}(syncer, sendChan, recvChan)

	return sendChan, recvChan
}

func syncerStartReceiving(syncer syncSource, incomingChan chan *UpdateRequest) chan error {
	incomingErrChan := make(chan error)

	go func(syncer syncSource, incomingChan chan *UpdateRequest) {
		for {
			updateReq, err := syncer.Recv()
			if err != nil {
				incomingErrChan <- errors.Wrap(err, "clientStartReceiving failed to Recv, terminating...")
				break
			}

			incomingChan <- updateReq
		}
	}(syncer, incomingChan)

	return incomingErrChan
}
