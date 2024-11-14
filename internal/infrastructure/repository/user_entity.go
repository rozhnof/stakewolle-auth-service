package pgrepo

import (
	"github.com/google/uuid"
	"github.com/rozhnof/stakewolle-auth-service/internal/domain/models"
)

type UserEntity struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	HashPassword string    `db:"hash_password"`
}

func userToModel(user *UserEntity) *models.User {
	return &models.User{}
}

func usersToModel(userEntityList []UserEntity) []models.User {
	userList := make([]models.User, 0, len(userEntityList))
	for _, userEntity := range userEntityList {
		userList = append(userList, *userToModel(&userEntity))
	}

	return userList
}

func userFromModel(user *models.User) *UserEntity {
	return &UserEntity{}
}

func usersFromModel(userList []models.User) []UserEntity {
	userEntityList := make([]UserEntity, 0, len(userList))
	for _, user := range userList {
		userEntityList = append(userEntityList, *userFromModel(&user))
	}

	return userEntityList
}
