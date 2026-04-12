package repository

import (
	"context"
	"likexuser/model"
)

type CreateInput struct {
	model.PublicID
	model.FullName
	PasswordHash string
}

// Repository defines the data access contract for user persistence operations.
type Repository interface {
	Create(context.Context, CreateInput) (model.ID, error)
	FirstIDByPublicID(context.Context, model.PublicID) (model.ID, error)
	FirstPasswordHashByPublicID(context.Context, model.PublicID) (string, error)
}
