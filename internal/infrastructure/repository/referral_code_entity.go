package pgrepo

import (
	"time"

	"github.com/google/uuid"
	"github.com/rozhnof/stakewolle-auth-service/internal/domain/models"
)

type ReferralCodeEntity struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	ExpiredAt time.Time `db:"expired_at"`
}

func referralCodeToModel(referralCode *ReferralCodeEntity) *models.ReferralCode {
	return &models.ReferralCode{
		ID:        referralCode.ID,
		UserID:    referralCode.UserID,
		ExpiredAt: referralCode.ExpiredAt,
	}
}

func referralCodeFromModel(referralCode *models.ReferralCode) *ReferralCodeEntity {
	return &ReferralCodeEntity{
		ID:        referralCode.ID,
		UserID:    referralCode.UserID,
		ExpiredAt: referralCode.ExpiredAt,
	}
}
