package keyservice

import (
	"github.com/cohix/simplcrypto"
	"github.com/pkg/errors"
)

// Manager defines an in-memory keyservice
type Manager struct {
	nodeKeypair *simplcrypto.KeyPair
}

// NewManager creates a new keyservice manager
func NewManager() (*Manager, error) {
	keypair, err := simplcrypto.GenerateMasterKeyPair()
	if err != nil {
		return nil, errors.Wrap(err, "failed to GenerateMasterKeyPair")
	}

	manager := &Manager{
		nodeKeypair: keypair,
	}

	return manager, nil
}

// PubKey returns the public key for the node
func (m *Manager) PubKey() *simplcrypto.SerializablePubKey {
	return m.nodeKeypair.SerializablePubKey()
}

// ReEncryptTaskKey re-encrypts a task key using a pubkey
func (m *Manager) ReEncryptTaskKey(taskKey *simplcrypto.Message, pubKey *simplcrypto.KeyPair) (*simplcrypto.Message, error) {
	decKeyJSON, err := m.nodeKeypair.Decrypt(taskKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Decrypt task key")
	}

	reEncKey, err := pubKey.Encrypt(decKeyJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Encrypt task key")
	}

	return reEncKey, nil
}

// Encrypt encrypts a message using the node keypair
func (m *Manager) Encrypt(body []byte) (*simplcrypto.Message, error) {
	msg, err := m.nodeKeypair.Encrypt(body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Encrypt body")
	}

	return msg, nil
}

// Sign signs something with the node keypair
func (m *Manager) Sign(body []byte) (*simplcrypto.Signature, error) {
	sig, err := m.nodeKeypair.Sign(body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Sign body")
	}

	return sig, nil
}

// Verify verifies a signature with the node keypair
func (m *Manager) Verify(body []byte, sig *simplcrypto.Signature) error {
	return m.nodeKeypair.Verify(body, sig)
}
