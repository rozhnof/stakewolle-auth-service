package models

import (
	"time"

	"github.com/google/uuid"
)

type ReferralCode struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ExpiredAt time.Time
}

func (t *ReferralCode) Valid() bool {
	return t.ExpiredAt.After(time.Now())
}
