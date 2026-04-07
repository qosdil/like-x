package repository

import (
	"context"
	"likexuser/model"
)

// Repository defines the data access contract for user persistence operations.
type Repository interface {
	Create(context.Context, model.CreateInput) (model.CreateOutput, error)
	FirstIDByPublicID(context.Context, model.PublicID) (model.ID, error)
	FirstPasswordHashByPublicID(context.Context, model.PublicID) (string, error)
}
