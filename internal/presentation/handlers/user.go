package http_handlers

import (
	"github.com/google/uuid"
	"github.com/rozhnof/stakewolle-auth-service/internal/domain/models"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func UserToModel(u User) models.User {
	return models.User{
		ID:       u.ID,
		Username: u.Username,
	}
}

func UserToDTO(u models.User) User {
	return User{
		ID:       u.ID,
		Username: u.Username,
	}
}
