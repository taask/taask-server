package partner

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/cohix/simplcrypto"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/config"
	"github.com/taask/taask-server/model"
	"github.com/taask/taask-server/service"
)

// Manager controls partner updating and health checking
type Manager struct {
	// our UUID, for auth purposes
	UUID string

	// the auth manager responsible for partners
	Auth auth.Manager

	// this is the same keypair that our authManager is created with,
	// allowing us to decrypt messages sent to us using that pubkey
	masterKeypair *simplcrypto.KeyPair

	// the partner we are syncing with
	partner *Partner

	// the auth info for the partner cluster
	config *config.ClientAuthConfig

	// the apply func is provided by the brain as a "callback" for sending updates
	applyUpdateFunc func(*Update)
}

// Partner describes a partner server
type Partner struct {
	UUID   string
	Client service.PartnerServiceClient

	host string
	port string

	// healthChecker sends heartbeats to the partner
	healthChecker healthChecker

	// The session that we (as the outgoing partner) hold for the incoming partner
	ActiveSession *activeSession

	// the sym key used to encrypt data between partners
	DataKey *simplcrypto.SymKey

	// the update that is being "buffered" for this partner
	Update     *Update
	updateLock *sync.Mutex
}

type activeSession struct {
	*auth.Session
	Keypair      *simplcrypto.KeyPair
	MasterPubKey *simplcrypto.KeyPair
}

// NewManager creates a new partner manager
func NewManager(config *config.ClientAuthConfig, masterKeypair *simplcrypto.KeyPair) *Manager {
	if config.Service == nil {
		return nil
	}

	partner := &Partner{
		host:       config.Service.Host,
		port:       config.Service.Port,
		updateLock: &sync.Mutex{},
	}

	uuid := model.NewPartnerUUID()

	return &Manager{
		UUID:    uuid,
		partner: partner,
		config:  config,
	}
}

func (m *Manager) streamUpdates(recvChan chan *Update) error {
	// the inner loop does partner sync (flushes the queued updates, receives updates)
	timeChan := make(chan time.Time)

	for {
		select {
		case update := <-recvChan:
			go m.applyUpdate(update)
		case <-timeChan:
			// TODO: determine if flushupdates should be allowed to set the next time or not
			go m.flushUpdates(timeChan)
		}
	}
}

func (m *Manager) applyUpdate(update *Update) {
	// TODO: determine if we should use a channel for this
	m.applyUpdateFunc(update)
}

func (m *Manager) flushUpdates(timeChan chan time.Time) {

	// timeChan <- time.After(time.Duration(time.Second * 5))
}

func (p *Partner) lockUnlock() func() {
	if p.updateLock == nil {
		p.updateLock = &sync.Mutex{} // TODO: this is an awful idea, right?
	}

	p.updateLock.Lock()

	return func() {
		p.updateLock.Unlock()
	}
}

func (m *Manager) decryptAndVerifyUpdateFromPartner(partner *Partner, updateReq *service.UpdateRequest) (*Update, error) {
	if partner.DataKey == nil {
		return nil, errors.New("missing data key for partner")
	}

	if updateReq.UpdateSignature == nil {
		return nil, errors.New("update request from partner missing signature")
	}

	updateJSON, err := partner.DataKey.Decrypt(updateReq.EncUpdate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Decrypt update from partner")
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

	update := Update{}
	if err := json.Unmarshal(updateJSON, &update); err != nil {
		return nil, errors.Wrap(err, "failed to Unmarshal update JSON from partner")
	}

	return &update, nil
}

// SetPartner allows the brain to set our partner (TODO: figure out if this is gross and if it will bite us later)
func (m *Manager) SetPartner(partner *Partner) {
	m.partner = partner
}
