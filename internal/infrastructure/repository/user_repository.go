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

type UserRepository struct {
	db        postgres.Database
	txManager postgres.TransactionManager
	log       *slog.Logger
	tracer    trace.Tracer
}

func NewUserRepository(db postgres.Database, txManager postgres.TransactionManager, log *slog.Logger, tracer trace.Tracer) *UserRepository {
	return &UserRepository{
		db:        db,
		txManager: txManager,
		log:       log,
		tracer:    tracer,
	}
}

func (s *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.Create")
	defer span.End()

	userEntity := userFromModel(user)

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, userQueryCreate, userEntity.Username, userEntity.HashPassword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	createdUserEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UserEntity])
	if err != nil {
		return nil, err
	}

	return userToModel(&createdUserEntity), nil
}

func (s *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.GetByID")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, userQueryGetByID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UserEntity])
	if err != nil {
		return nil, err
	}

	return userToModel(&userEntity), nil
}

func (s *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.GetByUsername")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, userQueryGetByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UserEntity])
	if err != nil {
		return nil, err
	}

	return userToModel(&userEntity), nil
}

func (s *UserRepository) List(ctx context.Context) ([]models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.List")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, userQueryList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userEntityList, err := pgx.CollectRows(rows, pgx.RowTo[UserEntity])
	if err != nil {
		return nil, err
	}

	return usersToModel(userEntityList), nil
}

func (s *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.Update")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, userQueryUpdate, user.Username, user.HashPassword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	updatedUserEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UserEntity])
	if err != nil {
		return nil, err
	}

	return userToModel(&updatedUserEntity), nil
}

func (s *UserRepository) Delete(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.Delete")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, userQueryDelete, userID)
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
