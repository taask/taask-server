package auth

import "github.com/cohix/simplcrypto"

// DefaultGroupUUID and others are consts for the internal (memory) auth manager
const (
	DefaultGroupUUID = "defaultgroupuuid"
	AdminGroupUUID   = "admingroupuuid"
	PartnerGroupUUID = "partnergroupuuid"
)

// Manager describes the interface for things that are able to manage auth
type Manager interface {
	AttemptAuth(attempt *Attempt) (*EncMemberSession, error)
	CheckAuth(session *Session) error
	CheckAuthEnsureAdmin(session *Session) error
	DeleteMemberAuth(uuid string) error
	AddGroup(group *MemberGroup) error
	MasterPubKey() *simplcrypto.SerializablePubKey
	EncryptForMember(memberUUID string, msg []byte) (*simplcrypto.Message, error)
	ReEncryptTaskKey(memberUUID string, encTaskKey *simplcrypto.Message) (*simplcrypto.Message, error)
	VerifySignatureFromMember(memberUUID string, msg []byte, sig *simplcrypto.Signature) error
}

// EncMemberSession is sent back to the member as an auth challenge
type EncMemberSession struct {
	EncSessionChallenge *simplcrypto.Message
}

// MemberAuth represents an existing member session
type MemberAuth struct {
	UUID             string
	GroupUUID        string
	SessionChallenge []byte
	PubKey           *simplcrypto.KeyPair
}
