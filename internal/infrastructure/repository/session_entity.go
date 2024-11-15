package pgrepo

import (
	"time"

	"github.com/google/uuid"
	"github.com/rozhnof/stakewolle-auth-service/internal/domain/models"
)

type SessionEntity struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiredAt    time.Time `db:"expired_at"`
	IsRevoked    bool      `db:"is_revoked"`
}

func sessionToModel(session *SessionEntity) *models.Session {
	return &models.Session{
		ID:     session.ID,
		UserID: session.UserID,
		RefreshToken: models.RefreshToken{
			Token:     session.RefreshToken,
			ExpiredAt: session.ExpiredAt,
			IsRevoked: session.IsRevoked,
		},
	}
}

func sessionFromModel(session *models.Session) *SessionEntity {
	return &SessionEntity{
		ID:           session.ID,
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken.Token,
		ExpiredAt:    session.RefreshToken.ExpiredAt,
		IsRevoked:    session.RefreshToken.IsRevoked,
	}
}
