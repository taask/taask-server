package brain

import "github.com/cohix/simplcrypto"

// GetMasterPartnerPubKey returns the master runner pubkey
func (m *Manager) GetMasterPartnerPubKey() *simplcrypto.SerializablePubKey {
	return m.PartnerManager.Auth.MasterPubKey()
}
