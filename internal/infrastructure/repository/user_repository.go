package pgrepo

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rozhnof/stakewolle-auth-service/internal/domain/models"
	"github.com/rozhnof/stakewolle-auth-service/internal/domain/repository"
	"github.com/rozhnof/stakewolle-auth-service/internal/infrastructure/database/postgres"
	db_queries "github.com/rozhnof/stakewolle-auth-service/internal/infrastructure/repository/queries"
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

func (s *UserRepository) GetUserIDByReferralCode(ctx context.Context, referralCode string) (*uuid.UUID, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.GetUserIDByReferralCode")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)
	querier := db_queries.New(db)

	userID, err := querier.GetUserIDByReferralCode(ctx, referralCode)
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

func (s *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.Create")
	defer span.End()

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := db_queries.New(db)

		userArgs := db_queries.CreateUserParams{
			Username:     user.Username,
			ReferrerID:   user.ReferrerID,
			HashPassword: user.HashPassword,
		}

		userRow, err := querier.CreateUser(ctx, userArgs)
		if err != nil {
			return err
		}

		user.ID = &userRow.ID

		return nil
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.GetByID")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)
	querier := db_queries.New(db)

	userRows, err := querier.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           &userRows.UserID,
		ReferrerID:   userRows.ReferrerID,
		Username:     userRows.Username,
		HashPassword: userRows.HashPassword,
	}

	if userRows.ReferralCodeID != nil {
		user.ReferralCode = &models.ReferralCode{
			Code:      *userRows.ReferralCode,
			ExpiredAt: *userRows.ReferralCodeExpiredAt,
		}
	}

	if userRows.SessionID != nil {
		user.Session = &models.Session{
			RefreshToken: models.RefreshToken{
				Token:     *userRows.RefreshToken,
				ExpiredAt: *userRows.SessionExpiredAt,
			},
		}
	}

	return user, nil
}

func (s *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.GetByUsername")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)
	querier := db_queries.New(db)

	userRows, err := querier.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           &userRows.UserID,
		ReferrerID:   userRows.ReferrerID,
		Username:     userRows.Username,
		HashPassword: userRows.HashPassword,
	}

	if userRows.ReferralCodeID != nil {
		user.ReferralCode = &models.ReferralCode{
			Code:      *userRows.ReferralCode,
			ExpiredAt: *userRows.ReferralCodeExpiredAt,
		}
	}

	if userRows.SessionID != nil {
		user.Session = &models.Session{
			RefreshToken: models.RefreshToken{
				Token:     *userRows.RefreshToken,
				ExpiredAt: *userRows.SessionExpiredAt,
			},
		}
	}

	return user, nil
}

func (s *UserRepository) List(ctx context.Context, filters repository.UserFilters, pagination repository.Pagination) ([]models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.List")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)
	querier := db_queries.New(db)

	args := db_queries.ListParams{
		Limit:       pagination.Limit,
		Offset:      pagination.Offset,
		UserIds:     filters.UserIDs,
		ReferrerIds: filters.ReferrerIDs,
	}

	userRows, err := querier.List(ctx, args)
	if err != nil {
		return nil, err
	}

	userList := make([]models.User, 0, len(userRows))

	for _, userRow := range userRows {
		user := models.User{
			ID:           &userRow.UserID,
			ReferrerID:   userRow.ReferrerID,
			Username:     userRow.Username,
			HashPassword: userRow.HashPassword,
		}

		if userRow.ReferralCodeID != nil {
			user.ReferralCode = &models.ReferralCode{
				Code:      *userRow.ReferralCode,
				ExpiredAt: *userRow.ReferralCodeExpiredAt,
			}
		}

		if userRow.SessionID != nil {
			user.Session = &models.Session{
				RefreshToken: models.RefreshToken{
					Token:     *userRow.RefreshToken,
					ExpiredAt: *userRow.SessionExpiredAt,
				},
			}
		}

		userList = append(userList, user)
	}

	return userList, nil
}

func (s *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.Update")
	defer span.End()

	if user.ID == nil {
		return nil, errors.Wrap(repository.ObjectNotFound, "user id is nil")
	}

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := db_queries.New(db)

		oldUserRows, err := querier.GetUserByID(ctx, *user.ID)

		userArgs := db_queries.UpdateUserParams{
			ID:           *user.ID,
			Username:     user.Username,
			HashPassword: user.HashPassword,
		}

		userRow, err := querier.UpdateUser(ctx, userArgs)
		if err != nil {
			return err
		}

		user.ReferrerID = userRow.ReferrerID
		user.Username = userRow.Username
		user.HashPassword = userRow.HashPassword

		if err := updateReferralCode(ctx, querier, user, &oldUserRows); err != nil {
			return err
		}

		if err := updateSession(ctx, querier, user, &oldUserRows); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func updateReferralCode(ctx context.Context, querier *db_queries.Queries, user *models.User, oldUserRows *db_queries.GetUserByIDRow) error {
	if user.ReferralCode == nil {
		if _, err := querier.DeleteSessionByUserID(ctx, *user.ID); err != nil {
			return err
		}

		return nil
	}

	if user.ReferralCode != nil && oldUserRows.ReferralCodeID == nil {
		if err := createSession(ctx, querier, user); err != nil {
			return err
		}

		return nil
	}

	// new user and old user have referral code
	if user.ReferralCode != nil && oldUserRows.ReferralCodeID != nil {
		// referral codes are different, revoke old and create new
		if user.ReferralCode.Code != *oldUserRows.ReferralCode {
			if _, err := querier.DeleteSessionByUserID(ctx, *user.ID); err != nil {
				return err
			}

			if err := createSession(ctx, querier, user); err != nil {
				return err
			}
		}
	}

	return nil
}

func updateSession(ctx context.Context, querier *db_queries.Queries, user *models.User, oldUserRows *db_queries.GetUserByIDRow) error {
	if user.Session == nil {
		if _, err := querier.DeleteSessionByUserID(ctx, *user.ID); err != nil {
			return err
		}

		return nil
	}

	if user.Session != nil && oldUserRows.SessionID == nil {
		if err := createSession(ctx, querier, user); err != nil {
			return err
		}

		return nil
	}

	// new user and old user have session
	if user.Session != nil && oldUserRows.SessionID != nil {
		// sessions are different, revoke old and create new
		if user.Session.RefreshToken.Token != *oldUserRows.RefreshToken {
			if _, err := querier.DeleteSessionByUserID(ctx, *user.ID); err != nil {
				return err
			}

			if err := createSession(ctx, querier, user); err != nil {
				return err
			}
		}
	}

	return nil
}

func createReferralCode(ctx context.Context, querier *db_queries.Queries, user *models.User) error {
	referralCodeArgs := db_queries.CreateReferralCodeParams{
		UserID:    *user.ID,
		Code:      user.ReferralCode.Code,
		ExpiredAt: user.ReferralCode.ExpiredAt,
	}

	referralCodeRow, err := querier.CreateReferralCode(ctx, referralCodeArgs)
	if err != nil {
		return err
	}

	user.ReferralCode.Code = referralCodeRow.Code
	user.ReferralCode.ExpiredAt = referralCodeRow.ExpiredAt

	return nil
}

func createSession(ctx context.Context, querier *db_queries.Queries, user *models.User) error {
	sessionArgs := db_queries.CreateSessionParams{
		UserID:       *user.ID,
		RefreshToken: user.Session.RefreshToken.Token,
		ExpiredAt:    user.Session.RefreshToken.ExpiredAt,
	}

	sessionRow, err := querier.CreateSession(ctx, sessionArgs)
	if err != nil {
		return err
	}

	user.Session.RefreshToken.Token = sessionRow.RefreshToken
	user.Session.RefreshToken.ExpiredAt = sessionRow.ExpiredAt

	return nil
}

func (s *UserRepository) Delete(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.Delete")
	defer span.End()

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := db_queries.New(db)

		if _, err := querier.DeleteUser(ctx, userID); err != nil {
			return err
		}

		if _, err := querier.DeleteReferralCodeByUserID(ctx, userID); err != nil {
			return err
		}

		if _, err := querier.DeleteSessionByUserID(ctx, userID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
}
