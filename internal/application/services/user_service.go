package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rozhnof/stakewolle-auth-service/internal/domain/models"
	"github.com/rozhnof/stakewolle-auth-service/internal/domain/repository"
	"go.opentelemetry.io/otel/trace"
)

type TokenManager interface {
	NewAccessToken() (models.AccessToken, error)
	NewRefreshToken() (models.RefreshToken, error)
}

type PasswordManager interface {
	HashPassword(password string) (string, error)
	CheckPassword(password string, hashPassword string) bool
}

const redisTTL = time.Hour * 60
const referralCodeTTL = time.Hour * 240

type Dependencies struct {
	UserRepository     repository.UserRepository
	UserCache          repository.UserCache
	TransactionManager repository.TransactionManager
	TokenManager       TokenManager
	PasswordManager    PasswordManager
}

func (d Dependencies) Valid() error {
	if d.UserRepository == nil {
		return errors.New("missing user repository")
	}

	if d.TokenManager == nil {
		return errors.New("missing token manager")
	}

	if d.TransactionManager == nil {
		return errors.New("missing transaction manager")
	}

	if d.PasswordManager == nil {
		return errors.New("missing password manager")
	}

	if d.UserCache == nil {
		return errors.New("missing user cache")
	}

	return nil
}

type UserService struct {
	Dependencies
	log    *slog.Logger
	tracer trace.Tracer
}

func NewUserService(d Dependencies, log *slog.Logger, tracer trace.Tracer) (*UserService, error) {
	if err := d.Valid(); err != nil {
		return nil, errors.Wrap(err, "missing required dependency")
	}

	return &UserService{
		Dependencies: d,
		log:          log,
		tracer:       tracer,
	}, nil
}

func (s *UserService) Register(ctx context.Context, username string, password string, referralCode string) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserService.Register")
	defer span.End()

	hashPassword, err := s.PasswordManager.HashPassword(password)
	if err != nil {
		return nil, err
	}

	var referrerID *uuid.UUID

	if referralCode != "" {
		userID, err := s.UserRepository.GetUserIDByReferralCode(ctx, referralCode)
		if err != nil {
			return nil, err
		}

		referrerID = userID
	}

	user := models.NewUser(username, hashPassword, referrerID)

	createdUser, err := s.UserRepository.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	s.UserCache.Set(ctx, username, *createdUser, redisTTL)

	return createdUser, nil
}

func (s *UserService) Login(ctx context.Context, username string, password string) (at string, rt string, err error) {
	ctx, span := s.tracer.Start(ctx, "UserService.Login")
	defer span.End()

	getUser := func(username string) (*models.User, error) {
		if cacheUser, err := s.UserCache.Get(ctx, username); err == nil {
			return &cacheUser, nil
		}

		dbUser, err := s.UserRepository.GetByUsername(ctx, username)
		if err != nil {
			return nil, err
		}

		return dbUser, nil
	}

	user, err := getUser(username)
	if err != nil {
		return "", "", err
	}

	if !s.PasswordManager.CheckPassword(password, user.HashPassword) {
		return "", "", ErrInvalidPassword
	}

	accessToken, err := s.TokenManager.NewAccessToken()
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.TokenManager.NewRefreshToken()
	if err != nil {
		return "", "", err
	}

	session := &models.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
	}

	if err := s.SessionRepository.RevokeByUserID(ctx, user.ID); err != nil {
		return "", "", err
	}

	if _, err := s.SessionRepository.Create(ctx, session); err != nil {
		return "", "", err
	}

	at = accessToken.Token
	rt = refreshToken.Token

	s.UserCache.Set(ctx, username, *user, redisTTL)

	return at, rt, nil
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (at string, rt string, err error) {
	ctx, span := s.tracer.Start(ctx, "UserService.Refresh")
	defer span.End()

	newAccessToken, err := s.TokenManager.NewAccessToken()
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.TokenManager.NewRefreshToken()
	if err != nil {
		return "", "", err
	}

	at = newAccessToken.Token
	rt = newRefreshToken.Token

	if err := s.TransactionManager.WithTransaction(ctx, func(ctx context.Context) error {
		session, err := s.SessionRepository.GetByRefreshToken(ctx, refreshToken)
		if err != nil {
			return ErrUnauthorizedRefresh
		}

		if !session.Valid() {
			return ErrUnauthorizedRefresh
		}

		session.RefreshToken = newRefreshToken

		if _, err := s.SessionRepository.Update(ctx, session); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", "", err
	}

	return at, rt, nil
}
