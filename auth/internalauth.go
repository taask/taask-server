package auth

import (
	"encoding/binary"
	"fmt"

	"github.com/cohix/simplcrypto"
	"github.com/pkg/errors"
)

// InternalAuthManager manages the auth of members, with in-memory storage
type InternalAuthManager struct {
	// the join code that runners need to sign in order to auth successfully
	memberGroups map[string]*MemberGroup

	// maps member UUIDs to memberAuths
	authedMembers map[string]*memberAuth

	// the master keypair
	masterKeyPair *simplcrypto.KeyPair
}

// NewInternalAuthManager returns a new InternalAuthManager
func NewInternalAuthManager() (*InternalAuthManager, error) {
	masterPair, err := simplcrypto.GenerateMasterKeyPair()
	if err != nil {
		return nil, errors.Wrap(err, "failed to GenerateMasterKeyPair")
	}

	manager := &InternalAuthManager{
		memberGroups:  make(map[string]*MemberGroup),
		authedMembers: make(map[string]*memberAuth),
		masterKeyPair: masterPair,
	}

	return manager, nil
}

// AttemptAuth checks the auth request for a member
// TODO: consider including the group UUID in the crypto meat grinder?
func (am *InternalAuthManager) AttemptAuth(attempt *Attempt) (*EncMemberSession, error) {
	group, ok := am.memberGroups[attempt.GroupUUID]
	if !ok {
		return nil, fmt.Errorf("failed to find member group with uuid %s", attempt.GroupUUID)
	}

	if !timestampIsValid(attempt.Timestamp) {
		return nil, errors.New("auth timestamp not within valid range")
	}

	nonce := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonce, uint64(attempt.Timestamp))
	hashWithNonce := append(group.AuthHash, nonce...)

	pubkey, err := simplcrypto.KeyPairFromSerializedPubKey(attempt.PubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to KeyPairFromSerializedPubKey")
	}

	if err := pubkey.Verify(hashWithNonce, attempt.AuthHashSig); err != nil {
		return nil, errors.Wrap(err, "failed to Verify")
	}

	challenge := newChallengeBytes()
	if challenge == nil {
		return nil, errors.New("failed to newChallengeBytes")
	}

	encChallenge, err := pubkey.Encrypt(challenge)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Encrypt challenge")
	}

	memberAuth := &memberAuth{
		UUID:             attempt.MemberUUID,
		GroupUUID:        attempt.GroupUUID,
		SessionChallenge: challenge,
		PubKey:           pubkey,
	}

	am.authedMembers[attempt.MemberUUID] = memberAuth

	encMemberChal := &EncMemberSession{
		EncSessionChallenge: encChallenge,
	}

	return encMemberChal, nil
}

// CheckAuth verifies a challenge signature is legit
func (am *InternalAuthManager) CheckAuth(session *Session) error {
	auth, ok := am.authedMembers[session.MemberUUID]
	if !ok {
		return errors.New(fmt.Sprintf("auth for member %s does not exist", session.MemberUUID))
	}

	if session.MemberUUID != auth.UUID {
		return errors.New("session member UUID does match auth UUID")
	}

	if session.GroupUUID != auth.GroupUUID {
		return errors.New("session GroupUUID does match auth GroupUUID")
	}

	if err := auth.PubKey.Verify(auth.SessionChallenge, session.SessionChallengeSig); err != nil {
		return errors.Wrap(err, "failed to Verify")
	}

	auth.SessionChallenge = nil

	return nil
}

// CheckAuthEnsureAdmin checks the auth against the admin group
func (am *InternalAuthManager) CheckAuthEnsureAdmin(session *Session) error {
	if session.GroupUUID != AdminGroupUUID {
		return errors.New("session not authenticated with admin group")
	}

	return am.CheckAuth(session)
}

// DeleteMemberAuth deletes a member's pubkey after it's been unregistered
func (am *InternalAuthManager) DeleteMemberAuth(memberUUID string) error {
	_, ok := am.authedMembers[memberUUID]
	if !ok {
		return errors.New(fmt.Sprintf("member %s key does not exist", memberUUID))
	}

	delete(am.authedMembers, memberUUID)

	return nil
}

// AddGroup adds a member group
func (am *InternalAuthManager) AddGroup(group *MemberGroup) error {
	_, existing := am.memberGroups[group.UUID]
	if existing {
		return fmt.Errorf("group with UUID %s already exists", group.UUID)
	}

	am.memberGroups[group.UUID] = group

	return nil
}

// MasterPubKey returns the pubkey from the master keypair
func (am *InternalAuthManager) MasterPubKey() *simplcrypto.SerializablePubKey {
	return am.masterKeyPair.SerializablePubKey()
}

// ReEncryptTaskKey re-encrypts a task key using a member's pubkey
func (am *InternalAuthManager) ReEncryptTaskKey(memberUUID string, encTaskKey *simplcrypto.Message) (*simplcrypto.Message, error) {
	memberAuth, ok := am.authedMembers[memberUUID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("member %s key does not exist", memberUUID))
	}

	decKeyJSON, err := am.masterKeyPair.Decrypt(encTaskKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Decrypt task key")
	}

	reEncKey, err := memberAuth.PubKey.Encrypt(decKeyJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Encrypt task key")
	}

	return reEncKey, nil
}
