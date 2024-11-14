package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rozhnof/stakewolle-auth-service/internal/domain/models"
	"github.com/rozhnof/stakewolle-auth-service/internal/infrastructure/database/redis"
)

type UserCache struct {
	redis.Database
}

func NewUserCache(db redis.Database) *UserCache {
	return &UserCache{
		Database: db,
	}
}

func (r *UserCache) Get(ctx context.Context, username string) (models.User, error) {
	key := createUserKey(username)

	bytes, err := r.Client.Get(ctx, key).Bytes()
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	if err := json.Unmarshal(bytes, &user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserCache) Set(ctx context.Context, username string, user models.User, ttl time.Duration) error {
	key := createUserKey(username)

	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, key, bytes, ttl).Err()
}

func (r *UserCache) Delete(ctx context.Context, username string) error {
	key := createUserKey(username)

	return r.Client.Del(ctx, key).Err()
}

func createUserKey(username string) string {
	return fmt.Sprintf("user:%s", username)
}
