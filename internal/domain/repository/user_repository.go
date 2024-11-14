package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rozhnof/stakewolle-auth-service/internal/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) (*time.Time, error)
}

type UserCache interface {
	Get(ctx context.Context, key string) (models.User, error)
	Set(ctx context.Context, key string, value models.User, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
