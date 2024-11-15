package models

import (
	"time"
)

type ReferralCode struct {
	Code      string
	ExpiredAt time.Time
}

func (t *ReferralCode) Valid() bool {
	return t.ExpiredAt.After(time.Now())
}
