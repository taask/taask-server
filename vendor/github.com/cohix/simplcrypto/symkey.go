package simplcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
)

const (
	symKeySize = 32
)

// SymKey respresents a symmetric key
type SymKey struct {
	Key string `json:"key"`
	KID string `json:"kid"`
}

// GenerateSymKey generates a new symmetric key
func GenerateSymKey() (*SymKey, error) {
	rawKey := make([]byte, symKeySize)
	if _, err := rand.Read(rawKey); err != nil {
		return nil, err
	}

	keyString := Base64URLEncode(rawKey)

	kid := generateNewKID()

	symKey := &SymKey{
		Key: keyString,
		KID: kid,
	}

	return symKey, nil
}

// Encrypt encrypts data into a Message
func (sk *SymKey) Encrypt(src []byte) (*Message, error) {
	rawKey, err := sk.rawKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(rawKey)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aead.NonceSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	ivString := Base64URLEncode(iv)

	encData := aead.Seal(nil, iv, src, nil)

	m := &Message{
		Data:    encData,
		KID:     sk.KID,
		KeyType: KeyTypeSymmetric,
		IV:      ivString,
	}

	return m, nil
}

// Decrypt returns decrypted data from a Message
func (sk *SymKey) Decrypt(src *Message) ([]byte, error) {
	if src.KID != sk.KID {
		return nil, fmt.Errorf("attempting to decrypt message with KID %q with symKey %q", src.KID, sk.KID)
	}

	if src.KeyType != KeyTypeSymmetric {
		return nil, fmt.Errorf("attempting to decrypt message encrypted with %q with key of type %q", src.KeyType, KeyTypeSymmetric)
	}

	rawKey, err := sk.rawKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(rawKey)
	if err != nil {
		return nil, err
	}

	iv, err := Base64URLDecode(src.IV)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decData, err := aead.Open(nil, iv, src.Data, nil)
	if err != nil {
		return nil, err
	}

	return decData, nil
}

// JSON returns the JSON representation of the symKey
func (sk *SymKey) JSON() []byte {
	keyJSON, _ := json.Marshal(sk)

	return keyJSON
}

// SymKeyFromJSON deserializes a key from JSON
func SymKeyFromJSON(keyJSON []byte) (*SymKey, error) {
	key := &SymKey{}

	if err := json.Unmarshal(keyJSON, key); err != nil {
		return nil, err
	}

	return key, nil
}

func (sk *SymKey) rawKey() ([]byte, error) {
	return Base64URLDecode(sk.Key)
}
