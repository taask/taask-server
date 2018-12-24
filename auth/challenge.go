package auth

import (
	"crypto/rand"
	"time"

	"github.com/cohix/simplcrypto"
)

const (
	challengeLen = 24
	joinCodeLen  = 32
)

// GenerateJoinCode generates a runner join code
func GenerateJoinCode() string {
	codeBytes := make([]byte, joinCodeLen)
	if _, err := rand.Read(codeBytes); err != nil {
		return ""
	}

	return simplcrypto.Base64URLEncode(codeBytes)
}

func timestampIsValid(timestamp int64) bool {
	realTime := time.Unix(timestamp, 0)
	now := time.Now()

	if now.Sub(realTime).Seconds() > 5 || now.Sub(realTime).Seconds() < 0 {
		return false
	}

	return true
}

func newChallengeBytes() []byte {
	chal := make([]byte, challengeLen)
	if _, err := rand.Read(chal); err != nil {
		return nil
	}

	return chal
}
