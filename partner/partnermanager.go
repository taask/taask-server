package partner

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cohix/simplcrypto"
	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/config"
	"github.com/taask/taask-server/keyservice"
	"github.com/taask/taask-server/model"
	"github.com/taask/taask-server/update"
)

const overridePartnerHostEnvKey = "TAASK_PARTNER_HOST"

// Manager controls partner updating and health checking
type Manager struct {
	// keyservice gives us access to functions of the node keypair
	keyservice keyservice.KeyService

	// our UUID, for auth purposes
	UUID string

	// the auth manager responsible for incoming partners
	Auth auth.Manager

	// the partner we are syncing with
	partner     *Partner
	partnerLock *sync.Mutex

	// the auth info for the partner cluster
	config *config.ClientAuthConfig

	// the apply func is provided by the brain as a "callback" for sending updates
	applyUpdateFunc func(update.PartnerUpdate)
}

type activeSession struct {
	*auth.Session
	Keypair      *simplcrypto.KeyPair
	MasterPubKey *simplcrypto.KeyPair
}

// NewManager creates a new partner manager
func NewManager(config *config.ClientAuthConfig, keyservice keyservice.KeyService) (*Manager, error) {
	if config.Service == nil {
		return nil, nil
	}

	host := config.Service.Host
	if envHost, useEnv := os.LookupEnv(overridePartnerHostEnvKey); useEnv && envHost != "" {
		host = envHost
		log.LogInfo(fmt.Sprintf("overriding partner host from env: %s", host))
	}

	partner := &Partner{
		Update:     update.NewPartnerUpdate(),
		host:       host,
		port:       config.Service.Port,
		updateLock: &sync.Mutex{},
	}

	uuid := model.NewPartnerUUID()

	authMan, err := auth.NewInternalAuthManager(keyservice)
	if err != nil {
		return nil, errors.Wrap(err, "failed to NewInternalAuthManager")
	}

	if config.MemberGroup.Name != "partner" {
		return nil, fmt.Errorf("partner auth config with group name %s not allowed", config.MemberGroup.Name)
	}

	if config.MemberGroup.UUID != auth.PartnerGroupUUID {
		return nil, fmt.Errorf("partner auth config with group uuid %s not allowed", config.MemberGroup.UUID)
	}

	if err := authMan.AddGroup(&config.MemberGroup); err != nil {
		return nil, errors.Wrap(err, "failed to AddGroup")
	}

	manager := &Manager{
		keyservice:  keyservice,
		UUID:        uuid,
		partner:     partner,
		partnerLock: &sync.Mutex{},
		config:      config,
		Auth:        authMan,
	}

	return manager, nil
}

// SetApplyUpdateFunc sets the applyUpdate callback func
func (m *Manager) SetApplyUpdateFunc(applyFunc func(update.PartnerUpdate)) {
	m.applyUpdateFunc = applyFunc
}

// SetPartnerUUID allows the brain to set our partner's uuid
func (m *Manager) SetPartnerUUID(uuid string) {
	m.partner.UUID = uuid
}

// HealthyPartner returns a UUID and pubkey of a healthy partner
func (m *Manager) HealthyPartner() (string, *simplcrypto.KeyPair) {
	if m == nil {
		return "", nil
	}

	if m.partner == nil {
		return "", nil
	}

	if m.partner.HealthChecker != nil {
		if m.partner.HealthChecker.IsHealthy {
			if m.partner.ActiveSession != nil {
				return m.partner.UUID, m.partner.ActiveSession.MasterPubKey
			} else {
				pubkey, err := m.Auth.MemberPubkey(m.partner.UUID)
				if err != nil {
					// TODO: handle this more gracefully
					return "", nil
				}

				return m.partner.UUID, pubkey
			}
		}
	}

	return "", nil
}

func (m *Manager) handleUpdates(sendChan, recvChan chan update.PartnerUpdate, unhealthyChan chan error) error {
	// the inner loop does partner sync (flushes the queued updates, receives updates)
	timeChan := make(chan bool, 1)
	timeChan <- true

	for {
		select {
		case update := <-recvChan:
			go m.applyUpdate(update)
		case <-timeChan:
			// TODO: determine if flushupdates should be allowed to set the next time or not
			go m.flushUpdates(sendChan, timeChan)
		case err := <-unhealthyChan:
			return errors.Wrap(err, "handleUpdates detected unhealthy partner, terminating update handler")
		}
	}
}

func (m *Manager) applyUpdate(update update.PartnerUpdate) {
	// TODO: determine if we should use a channel for this
	log.LogInfo("applying update from partner")

	m.applyUpdateFunc(update)
}

func (m *Manager) flushUpdates(sendChan chan update.PartnerUpdate, timeChan chan bool) {
	defer m.partner.lockUnlockUpdate()()
	log.LogInfo("flushing updates to partner")

	updateToSend := *m.partner.Update

	if len(updateToSend.Tasks) == 0 && len(updateToSend.Groups) == 0 && len(updateToSend.Sessions) == 0 {
		log.LogInfo("no updates queued for partner")
	} else {
		sendChan <- updateToSend

		m.partner.Update = update.NewPartnerUpdate()
	}

	<-time.After(time.Duration(time.Second * 5))

	timeChan <- true
}

func (m *Manager) decryptAndVerifyUpdateFromPartner(partner *Partner, updateReq *UpdateRequest) (*update.PartnerUpdate, error) {
	if updateReq.IsHealthCheck {
		return nil, nil
	}

	if partner.DataKey == nil {
		return nil, errors.New("missing data key for partner")
	}

	updateJSON, err := partner.DataKey.Decrypt(updateReq.EncUpdate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Decrypt update from partner")
	}

	if updateReq.UpdateSignature == nil {
		return nil, errors.New("update request from partner missing signature")
	}

	if partner.ActiveSession != nil {
		// if we are the outgoing partner, verify using the pubkey the server sent us
		if err := partner.ActiveSession.MasterPubKey.Verify(updateJSON, updateReq.UpdateSignature); err != nil {
			return nil, errors.Wrap(err, "failed to Verify update from incoming partner")
		}
	} else {
		// if we are the incoming partner, verify using the pubkey associated with the outgoing partner's auth
		if err := m.Auth.VerifySignatureFromMember(partner.UUID, updateJSON, updateReq.UpdateSignature); err != nil {
			return nil, errors.Wrap(err, "failed to Verify update from outgoing partner")
		}
	}

	update := update.PartnerUpdate{}
	if err := json.Unmarshal(updateJSON, &update); err != nil {
		return nil, errors.Wrap(err, "failed to Unmarshal update JSON from partner")
	}

	return &update, nil
}

func (m *Manager) encryptAndSignUpdateForPartner(partner *Partner, update *update.PartnerUpdate) (*UpdateRequest, error) {
	updateReq := &UpdateRequest{}

	if partner.DataKey == nil {
		return nil, errors.New("missing data key for partner")
	}

	updateJSON, err := json.Marshal(update)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal update")
	}

	encUpdate, err := partner.DataKey.Encrypt(updateJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Encrypt update")
	}

	var updateSig *simplcrypto.Signature
	var signErr error

	if partner.ActiveSession != nil {
		// if we are the outgoing partner, sign with our session keypair
		updateSig, signErr = partner.ActiveSession.Keypair.Sign(updateJSON)
		if signErr != nil {
			return nil, errors.Wrap(signErr, "failed to Sign update with activeSession.KeyPair")
		}

		updateReq.Session = partner.ActiveSession.Session
	} else {
		// if we are the incoming partner, sign with our node keypair
		updateSig, signErr = m.keyservice.Sign(updateJSON)
		if signErr != nil {
			return nil, errors.Wrap(signErr, "failed to Sign update with masterKeypair")
		}
	}

	updateReq.EncUpdate = encUpdate
	updateReq.UpdateSignature = updateSig

	return updateReq, nil
}
