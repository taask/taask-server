package simplcrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

const (
	// MasterKeyPairKID is the KID for the master node's keyPair
	MasterKeyPairKID = "simpl.key.masterkeypair"
)

// KeyPair stores a private/public pair to represent an AstroCache node
type KeyPair struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
	KID     string
}

// GenerateNewKeyPair generates a new KeyPair
func GenerateNewKeyPair() (*KeyPair, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	pair := &KeyPair{
		Private: priv,
		Public:  &priv.PublicKey,
		KID:     generateNewKID(),
	}

	return pair, nil
}

// GenerateMasterKeyPair generates a keyPair for the master node
func GenerateMasterKeyPair() (*KeyPair, error) {
	keyPair, err := GenerateNewKeyPair()
	if err != nil {
		return nil, err
	}

	keyPair.KID = MasterKeyPairKID

	return keyPair, nil
}

// Encrypt performs rsaOAEP on the input bytes and returns an encrypted message
func (kp *KeyPair) Encrypt(src []byte) (*Message, error) {
	encData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, kp.Public, src, nil)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		Data:    encData,
		KeyType: KeyTypePair,
		KID:     kp.KID,
	}

	return msg, nil
}

// Decrypt performs rsaOAEP decryption on the source message
func (kp *KeyPair) Decrypt(src *Message) ([]byte, error) {
	if kp.Private == nil {
		return nil, errors.New("attempted to decrypt message with nil private key")
	}

	if src.KeyType != KeyTypePair {
		return nil, fmt.Errorf("attempting to decrypt message encrypted with %q with key of type %q", src.KeyType, KeyTypePair)
	}

	if src.KID != kp.KID {
		return nil, fmt.Errorf("attempted to decrypt message with KID %q with keyPair %q", src.KID, kp.KID)
	}

	decMsg, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, kp.Private, src.Data, nil)
	if err != nil {
		return nil, err
	}

	return decMsg, nil
}

// Sign creates a DSS with a keyPair
func (kp *KeyPair) Sign(src []byte) (*Signature, error) {
	if kp.Private == nil {
		return nil, errors.New("attempting to sign data with nil private key")
	}

	h := crypto.SHA256
	hasher := h.New()
	hasher.Write(src)
	hashed := hasher.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, kp.Private, h, hashed)
	if err != nil {
		return nil, err
	}

	out := &Signature{
		Signature: sig,
		KID:       kp.KID,
	}

	return out, nil
}

// SigVerified and others represent results of a signature verification
const (
	SigVerified   = true
	SigUnverified = false
)

// Verify verifies src with kp's pubKey and the provided signature
func (kp *KeyPair) Verify(src []byte, sig *Signature) error {
	if sig.KID != kp.KID {
		return fmt.Errorf("signature KID %s does not match key KID %s", sig.KID, kp.KID)
	}

	h := crypto.SHA256
	hasher := h.New()
	hasher.Write(src)
	hashed := hasher.Sum(nil)

	if err := rsa.VerifyPKCS1v15(kp.Public, h, hashed, sig.Signature); err != nil {
		return err
	}

	return nil
}

// KeyPairFromPubKeyJSON unmarshals and de-serializes a serializablePubKey from JSON so it can be used to encrypt or validate signatures
func KeyPairFromPubKeyJSON(src []byte) (*KeyPair, error) {
	serialized := SerializablePubKey{}
	if err := json.Unmarshal(src, &serialized); err != nil {
		return nil, err
	}

	pubKey, err := serialized.deserialize()
	if err != nil {
		return nil, err
	}

	keyPair := &KeyPair{
		Public: pubKey,
		KID:    serialized.KID,
	}

	return keyPair, nil
}

// KeyPairFromSerializedPubKey de-serializes a serializablePubKey from JSON so it can be used to encrypt or validate signatures
func KeyPairFromSerializedPubKey(serialized *SerializablePubKey) (*KeyPair, error) {
	pubKey, err := serialized.deserialize()
	if err != nil {
		return nil, err
	}

	keyPair := &KeyPair{
		Public: pubKey,
		KID:    serialized.KID,
	}

	return keyPair, nil
}

// PubKeyJSON exports the KeyPair's pubKey to JSON using serializablePubKey
func (kp *KeyPair) PubKeyJSON() []byte {
	base64N := Base64URLEncode(kp.Public.N.Bytes())

	serializable := SerializablePubKey{
		N:   base64N,
		E:   int64(kp.Public.E),
		KID: kp.KID,
	}

	json, _ := json.Marshal(serializable)

	return json
}

// SerializablePubKey exports the KeyPair's pubKey to serializablePubKey
func (kp *KeyPair) SerializablePubKey() *SerializablePubKey {
	base64N := Base64URLEncode(kp.Public.N.Bytes())

	serializable := &SerializablePubKey{
		N:   base64N,
		E:   int64(kp.Public.E),
		KID: kp.KID,
	}

	return serializable
}

// SerializablePubKey is defined in keypair.pb.go
func (spk *SerializablePubKey) deserialize() (*rsa.PublicKey, error) {
	realN := &big.Int{}
	nBytes, err := Base64URLDecode(spk.N)
	if err != nil {
		return nil, err
	}

	realN.SetBytes(nBytes)

	pubKey := &rsa.PublicKey{
		N: realN,
		E: int(spk.E),
	}

	return pubKey, nil
}

func generateNewKID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)

	return Base64URLEncode(bytes)
}
