package repository

import "github.com/google/uuid"

type UserFilters struct {
	UserIDs     []uuid.UUID
	ReferrerIDs []uuid.UUID
}
