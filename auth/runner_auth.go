package auth

import (
	"bytes"
	"fmt"

	"github.com/cohix/simplcrypto"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/model"
)

// RunnerAuthManager manages the auth of runners
type RunnerAuthManager struct {
	// the join code that runners need to sign in order to auth successfully
	JoinCode string

	// tracks signatures already used for auth to prevent replay auths
	usedSignatures []*simplcrypto.Signature

	// maps
	authedRunnerChallenges map[string]*runnerChallenge
}

// EncRunnerChallenge is sent back to the runner as an auth challenge
type EncRunnerChallenge struct {
	EncChallenge    *simplcrypto.Message
	EncChallengeKey *simplcrypto.Message
}

type runnerChallenge struct {
	Challenge []byte
	PubKey    *simplcrypto.KeyPair
}

// NewRunnerAuthManager returns a new RunnerAuthManager
func NewRunnerAuthManager(joinCode string) *RunnerAuthManager {
	return &RunnerAuthManager{
		JoinCode:               joinCode,
		usedSignatures:         []*simplcrypto.Signature{},
		authedRunnerChallenges: make(map[string]*runnerChallenge),
	}
}

// AttemptAuth checks the auth request for a runner
func (am *RunnerAuthManager) AttemptAuth(pubKey *simplcrypto.SerializablePubKey, codeSig *simplcrypto.Signature) (*model.AuthRunnerResponse, error) {
	for _, sig := range am.usedSignatures {
		if sig.KID == codeSig.KID {
			return nil, errors.New("runner pubKey previously used, generate new pubKey and try again")
		}

		if bytes.Compare(sig.Signature, codeSig.Signature) == 0 {
			return nil, errors.New("auth signature previously used")
		}
	}

	runnerKey, err := simplcrypto.KeyPairFromSerializedPubKey(pubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to KeyPairFromSerializedPubKey")
	}

	if err := runnerKey.Verify([]byte(am.JoinCode), codeSig); err != nil {
		return nil, errors.Wrap(err, "failed to Verify")
	}

	// at this point, the auth succeeds, so we generate and encrypt the challenge
	am.usedSignatures = append(am.usedSignatures, codeSig)

	challenge := newChallengeBytes()
	if challenge == nil {
		return nil, errors.New("failed to newChallengeBytes")
	}

	challengeSymKey, err := simplcrypto.GenerateSymKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to GenerateSymKey")
	}

	encChallenge, err := challengeSymKey.Encrypt(challenge)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Encrypt challenge")
	}

	encChallengeKey, err := runnerKey.Encrypt(challengeSymKey.JSON())
	if err != nil {
		return nil, errors.Wrap(err, "failed to Encrypt challengeSymKey")
	}

	runnerChal := &runnerChallenge{
		Challenge: challenge,
		PubKey:    runnerKey,
	}

	am.authedRunnerChallenges[runnerKey.KID] = runnerChal

	encRunnerChal := &model.AuthRunnerResponse{
		EncChallenge:    encChallenge,
		EncChallengeKey: encChallengeKey,
	}

	return encRunnerChal, nil
}

// CheckRunnerChallenge verifies a challenge signature is legit and then deletes it from existence
func (am *RunnerAuthManager) CheckRunnerChallenge(chalSig *simplcrypto.Signature) error {
	challenge, ok := am.authedRunnerChallenges[chalSig.KID]
	if !ok {
		return errors.New(fmt.Sprintf("runner challenge for KID %s does not exist", chalSig.KID))
	}

	if err := challenge.PubKey.Verify(challenge.Challenge, chalSig); err != nil {
		return errors.Wrap(err, "failed to Verify")
	}

	delete(am.authedRunnerChallenges, chalSig.KID)

	return nil
}
