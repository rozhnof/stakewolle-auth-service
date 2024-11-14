package pgrepo

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rozhnof/stakewolle-auth-service/internal/domain/models"
	"github.com/rozhnof/stakewolle-auth-service/internal/infrastructure/database/postgres"
	"go.opentelemetry.io/otel/trace"
)

type SessionRepository struct {
	db        postgres.Database
	txManager postgres.TransactionManager
	log       *slog.Logger
	tracer    trace.Tracer
}

func NewSessionRepository(db postgres.Database, txManager postgres.TransactionManager, log *slog.Logger, tracer trace.Tracer) *SessionRepository {
	return &SessionRepository{
		db:        db,
		txManager: txManager,
		log:       log,
		tracer:    tracer,
	}
}

func (s *SessionRepository) Create(ctx context.Context, session *models.Session) (*models.Session, error) {
	ctx, span := s.tracer.Start(ctx, "SessionRepository.Create")
	defer span.End()

	sessionEntity := sessionFromModel(session)

	args := []any{
		sessionEntity.UserID,
		sessionEntity.RefreshToken,
		sessionEntity.ExpiredAt,
		sessionEntity.IsRevoked,
	}

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, sessionQueryCreate, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	createdSessionEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[SessionEntity])
	if err != nil {
		return nil, err
	}

	return sessionToModel(&createdSessionEntity), nil
}

func (s *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	ctx, span := s.tracer.Start(ctx, "SessionRepository.GetByRefreshToken")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, sessionQueryGetByRefreshToken, refreshToken)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessionEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[SessionEntity])
	if err != nil {
		return nil, err
	}

	return sessionToModel(&sessionEntity), nil
}

func (s *SessionRepository) GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	ctx, span := s.tracer.Start(ctx, "SessionRepository.GetByID")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, sessionQueryGetByID, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessionEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[SessionEntity])
	if err != nil {
		return nil, err
	}

	return sessionToModel(&sessionEntity), nil
}

func (s *SessionRepository) Update(ctx context.Context, session *models.Session) (*models.Session, error) {
	ctx, span := s.tracer.Start(ctx, "SessionRepository.Update")
	defer span.End()

	sessionEntity := sessionFromModel(session)

	args := []any{
		sessionEntity.ID,
		sessionEntity.UserID,
		sessionEntity.RefreshToken,
		sessionEntity.ExpiredAt,
		sessionEntity.IsRevoked,
	}

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, sessionQueryUpdate, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	updatedSessionEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[SessionEntity])
	if err != nil {
		return nil, err
	}

	return sessionToModel(&updatedSessionEntity), nil
}

func (s *SessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) (*time.Time, error) {
	ctx, span := s.tracer.Start(ctx, "SessionRepository.Delete")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, sessionQueryDelete, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deletedAt, err := pgx.RowTo[time.Time](rows)
	if err != nil {
		return nil, err
	}

	return &deletedAt, nil
}

func (s *SessionRepository) RevokeByUserID(ctx context.Context, userID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "SessionRepository.RevokeByUserID")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	_, err := db.Exec(ctx, sessionQueryRevoke, userID)
	if err != nil {
		return err
	}

	return nil
}
