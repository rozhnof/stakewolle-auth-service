package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rozhnof/stakewolle-auth-service/internal/domain/models"
)

type ReferralCodeRepository interface {
	Create(ctx context.Context, referralCode *models.ReferralCode) (*models.ReferralCode, error)
	Delete(ctx context.Context, referralCodeID uuid.UUID) (*time.Time, error)
	GetByUsername(ctx context.Context, username string) (*models.ReferralCode, error)
}

type ReferralCodeCache interface {
	Get(ctx context.Context, key string) (models.ReferralCode, error)
	Set(ctx context.Context, key string, value models.ReferralCode, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
