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

type ReferralCodeRepository struct {
	db        postgres.Database
	txManager postgres.TransactionManager
	log       *slog.Logger
	tracer    trace.Tracer
}

func NewReferralCodeRepository(db postgres.Database, txManager postgres.TransactionManager, log *slog.Logger, tracer trace.Tracer) *ReferralCodeRepository {
	return &ReferralCodeRepository{
		db:        db,
		txManager: txManager,
		log:       log,
		tracer:    tracer,
	}
}

func (s *ReferralCodeRepository) Create(ctx context.Context, referralCode *models.ReferralCode) (*models.ReferralCode, error) {
	ctx, span := s.tracer.Start(ctx, "ReferralCodeRepository.Create")
	defer span.End()

	referralCodeEntity := referralCodeFromModel(referralCode)

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, referralCodeQueryCreate, referralCodeEntity.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	createdReferralCodeEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[ReferralCodeEntity])
	if err != nil {
		return nil, err
	}

	return referralCodeToModel(&createdReferralCodeEntity), nil
}

func (s *ReferralCodeRepository) GetByUsername(ctx context.Context, username string) (*models.ReferralCode, error) {
	ctx, span := s.tracer.Start(ctx, "ReferralCodeRepository.GetByUsername")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, referralCodeQueryGetByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	referralCodeEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[ReferralCodeEntity])
	if err != nil {
		return nil, err
	}

	return referralCodeToModel(&referralCodeEntity), nil
}

func (s *ReferralCodeRepository) Delete(ctx context.Context, referralCodeID uuid.UUID) (*time.Time, error) {
	ctx, span := s.tracer.Start(ctx, "ReferralCodeRepository.Delete")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, referralCodeQueryDelete, referralCodeID)
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
