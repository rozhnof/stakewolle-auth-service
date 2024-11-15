package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID           *uuid.UUID
	ReferrerID   *uuid.UUID
	Username     string
	HashPassword string

	Session      *Session
	ReferralCode *ReferralCode
}

func NewUser(username string, hashPassword string, referrerID *uuid.UUID) User {
	return User{
		ReferrerID:   referrerID,
		Username:     username,
		HashPassword: hashPassword,
	}
}
