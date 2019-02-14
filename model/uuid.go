package model

import (
	"math/rand"
	"time"
)

// UUIDLength and others are consts for UUIDs
const (
	UUIDLength        = 26
	lowercaseAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	uppercaseAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewTaskUUID returns a new task UUID
func NewTaskUUID() string {
	return lowercaseUUID()
}

// NewRunnerUUID returns a new runner UUID
func NewRunnerUUID() string {
	return uppercaseUUID()
}

// NewPartnerUUID returns a new partner UUID
func NewPartnerUUID() string {
	return uppercaseUUID()
}

// NewMemberGroupUUID returns a new task UUID
func NewMemberGroupUUID() string {
	return uppercaseUUID()
}

func uppercaseUUID() string {
	uuid := ""

	for i := 0; i < UUIDLength; i++ {
		index := rand.Intn(len(uppercaseAlphabet))

		uuid += string(uppercaseAlphabet[index])
	}

	return uuid
}

func lowercaseUUID() string {
	uuid := ""

	for i := 0; i < UUIDLength; i++ {
		index := rand.Intn(len(lowercaseAlphabet))

		uuid += string(lowercaseAlphabet[index])
	}

	return uuid
}
