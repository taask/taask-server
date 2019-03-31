package keyservice

import "github.com/cohix/simplcrypto"

// KeyService defines the interface for a keyservice
type KeyService interface {
	PubKey() *simplcrypto.SerializablePubKey

	Encrypt(body []byte) (*simplcrypto.Message, error)
	ReEncryptTaskKey(taskKey *simplcrypto.Message, pubKey *simplcrypto.KeyPair) (*simplcrypto.Message, error)

	Sign(body []byte) (*simplcrypto.Signature, error)
	Verify(body []byte, sig *simplcrypto.Signature) error
}
