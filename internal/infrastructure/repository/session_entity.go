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
		// RefreshToken: models.RefreshToken{
		// 	Token:     session.RefreshToken,
		// 	ExpiredAt: session.ExpiredAt,
		// 	IsRevoked: session.IsRevoked,
		// },
	}
}

func sessionsToModel(sessionEntityList []SessionEntity) []models.Session {
	sessionList := make([]models.Session, 0, len(sessionEntityList))
	for _, sessionEntity := range sessionEntityList {
		sessionList = append(sessionList, *sessionToModel(&sessionEntity))
	}

	return sessionList
}

func sessionFromModel(session *models.Session) *SessionEntity {
	return &SessionEntity{
		ID:     session.ID,
		UserID: session.UserID,
		// RefreshToken: session.RefreshToken.Token,
		// ExpiredAt:    session.RefreshToken.ExpiredAt,
		// IsRevoked:    session.RefreshToken.IsRevoked,
	}
}

func sessionsFromModel(sessionList []models.Session) []SessionEntity {
	sessionEntityList := make([]SessionEntity, 0, len(sessionList))
	for _, session := range sessionList {
		sessionEntityList = append(sessionEntityList, *sessionFromModel(&session))
	}

	return sessionEntityList
}
