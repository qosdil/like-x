package service

import (
	"context"
	user "likexuser/model"
	"likexuser/repository"
	"log"

	likexService "github.com/qosdil/like-x/backend/common/service"
)

// Service defines business logic for user-related operations.
type Service struct {
	repo repository.Repository
}

// NewService constructs a new Service with the provided user repository.
func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

// SignUp validates user input and creates a new user via the repository.
func (s *Service) SignUp(ctx context.Context, input user.CreateInput) (user.CreateOutput, error) {
	// Validate full_name
	if len(input.FullName) < user.FullNameMinLength || len(input.FullName) > user.FullNameMaxlength {
		return user.CreateOutput{}, likexService.ErrBadRequest
	}

	// Validate password
	if len(input.Password) < user.PasswordMinLength || len(input.Password) > user.PasswordMaxLength {
		return user.CreateOutput{}, likexService.ErrBadRequest
	}

	// Create the user in the repository and handle any errors.
	signUp, err := s.repo.Create(ctx, input)
	if err != nil {
		log.Printf("failed to sign up a user: %v", err)
		return user.CreateOutput{}, likexService.ErrInternal
	}

	return signUp, nil
}
