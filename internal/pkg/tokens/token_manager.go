package tokens

import (
	"time"

	"github.com/rozhnof/stakewolle-auth-service/internal/domain/models"
	"github.com/rozhnof/stakewolle-auth-service/internal/pkg/secrets"
)

type TokenManager struct {
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	secretManager   secrets.SecretManager
}

func NewTokenManager(accessTokenTTL time.Duration, refreshTokenTTL time.Duration, secretManager secrets.SecretManager) TokenManager {
	return TokenManager{
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		secretManager:   secretManager,
	}
}

func (t TokenManager) NewAccessToken() (models.AccessToken, error) {
	return models.NewAccessToken(t.accessTokenTTL, t.secretManager.SecretKey())
}

func (t TokenManager) NewRefreshToken() (models.RefreshToken, error) {
	return models.NewRefreshToken(t.refreshTokenTTL)
}
