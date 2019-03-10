package partner

import (
	"sync"

	"github.com/cohix/simplcrypto"
	"github.com/taask/taask-server/update"
)

// Partner describes a partner server
type Partner struct {
	UUID   string
	Client PartnerServiceClient

	host string
	port string

	// healthChecker maintains the health of the partner
	HealthChecker *healthChecker

	// The session that we (as the outgoing partner) hold for the incoming partner
	ActiveSession *activeSession

	// the sym key used to encrypt data between partners
	DataKey *simplcrypto.SymKey

	// the update that is being "buffered" for this partner
	Update *update.PartnerUpdate

	// we need to be able to lock the update such that attempts
	// to flush it don't conflict with attempts to update it
	updateLock *sync.Mutex
}

func (p *Partner) lockUnlockUpdate() func() {
	p.updateLock.Lock()

	return func() {
		p.updateLock.Unlock()
	}
}
