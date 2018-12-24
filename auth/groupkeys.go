package auth

import (
	"crypto/sha256"

	"github.com/cohix/simplcrypto"
	"github.com/pkg/errors"

	"golang.org/x/crypto/pbkdf2"
)

// GroupKeyIterations and others are the consts for member auth
const (
	GroupKeyIterations = 100000
	GroupKeySaltString = "to there and back again" // TODO: store the salt with the group
)

// GroupAuthHash generates the auth hash from a Join Code and a passphrase
func GroupAuthHash(joinCode, passphrase string) []byte {
	sha := sha256.New()

	codeBytes := []byte(joinCode)
	passBytes := []byte(passphrase)

	hashIn := append(codeBytes, passBytes...)

	sha.Write(hashIn)

	return sha.Sum(nil)
}

// GroupDerivedKey derives a symmetric key from a passphrase
func GroupDerivedKey(passphrase string) (*simplcrypto.SymKey, error) {
	keyBytes := pbkdf2.Key([]byte(passphrase), []byte(GroupKeySaltString), 100000, sha256.Size, sha256.New)

	key, err := simplcrypto.GenerateSymKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to GenerateSymKey")
	}

	key.Key = simplcrypto.Base64URLEncode(keyBytes)

	return key, nil
}
