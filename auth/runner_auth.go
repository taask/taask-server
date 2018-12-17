package auth

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/cohix/simplcrypto"
	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
)

const (
	joinCodeWritePath = "/taask/config/joincode"
)

type runnerAuth struct {
	Challenge []byte
	PubKey    *simplcrypto.KeyPair
}

// RunnerAuthManager manages the auth of runners
type RunnerAuthManager struct {
	// the join code that runners need to sign in order to auth successfully
	JoinCode string

	// tracks signatures already used for auth to prevent replay auths
	// TODO: nuke this and use a nonce
	usedSignatures []*simplcrypto.Signature

	// maps pubkey KIDs to runnerAuths
	authedRunners map[string]*runnerAuth

	// maps runners to their pubkeys
	activeRunnerKeys map[string]*simplcrypto.KeyPair

	// the keypair used to transfer access to task data to runners
	runnerMasterKeyPair *simplcrypto.KeyPair
}

// EncRunnerAuth is sent back to the runner as an auth challenge
type EncRunnerAuth struct {
	EncChallenge    *simplcrypto.Message
	EncChallengeKey *simplcrypto.Message
}

// NewRunnerAuthManager returns a new RunnerAuthManager
func NewRunnerAuthManager(joinCode string) (*RunnerAuthManager, error) {
	defer writeJoinCode(joinCode)

	masterPair, err := simplcrypto.GenerateMasterKeyPair()
	if err != nil {
		return nil, errors.Wrap(err, "failed to GenerateMasterKeyPair")
	}

	manager := &RunnerAuthManager{
		JoinCode:            joinCode,
		usedSignatures:      []*simplcrypto.Signature{},
		authedRunners:       make(map[string]*runnerAuth),
		runnerMasterKeyPair: masterPair,
	}

	return manager, nil
}

// AttemptAuth checks the auth request for a runner
func (am *RunnerAuthManager) AttemptAuth(pubKey *simplcrypto.SerializablePubKey, codeSig *simplcrypto.Signature) (*EncRunnerAuth, error) {
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

	runnerChal := &runnerAuth{
		Challenge: challenge,
		PubKey:    runnerKey,
	}

	am.authedRunners[runnerKey.KID] = runnerChal

	encRunnerChal := &EncRunnerAuth{
		EncChallenge:    encChallenge,
		EncChallengeKey: encChallengeKey,
	}

	return encRunnerChal, nil
}

// CheckRunnerAuth verifies a challenge signature is legit and then deletes it from existence, retaining the pubkey
func (am *RunnerAuthManager) CheckRunnerAuth(runnerUUID string, chalSig *simplcrypto.Signature) error {
	auth, ok := am.authedRunners[chalSig.KID]
	if !ok {
		return errors.New(fmt.Sprintf("runner auth for KID %s does not exist", chalSig.KID))
	}

	if err := auth.PubKey.Verify(auth.Challenge, chalSig); err != nil {
		return errors.Wrap(err, "failed to Verify")
	}

	am.activeRunnerKeys[runnerUUID] = auth.PubKey

	delete(am.authedRunners, chalSig.KID)

	return nil
}

// DeleteRunnerKey deletes a runner's pubkey after it's been unregistered
func (am *RunnerAuthManager) DeleteRunnerKey(uuid string) error {
	_, ok := am.activeRunnerKeys[uuid]
	if !ok {
		return errors.New(fmt.Sprintf("runner %s key does not exist", uuid))
	}

	delete(am.activeRunnerKeys, uuid)

	return nil
}

// RunnerMasterPubKey returns the pubkey from the master keypair
func (am *RunnerAuthManager) RunnerMasterPubKey() *simplcrypto.SerializablePubKey {
	return am.runnerMasterKeyPair.SerializablePubKey()
}

// ReEncryptTaskKey re-encrypts a task key using a runner's pubkey
func (am *RunnerAuthManager) ReEncryptTaskKey(runnerUUID string, encTaskKey *simplcrypto.Message) (*simplcrypto.Message, error) {
	runnerPubKey, ok := am.activeRunnerKeys[runnerUUID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("runner %s key does not exist", runnerUUID))
	}

	decKeyJSON, err := am.runnerMasterKeyPair.Decrypt(encTaskKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Decrypt task key")
	}

	reEncKey, err := runnerPubKey.Encrypt(decKeyJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Encrypt task key")
	}

	return reEncKey, nil
}

func writeJoinCode(joinCode string) {
	if err := ioutil.WriteFile(joinCodeWritePath, []byte(joinCode), 0666); err != nil {
		log.LogError(errors.Wrap(err, "failed to WriteFile join code"))
	}
}
