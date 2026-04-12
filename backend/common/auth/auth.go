package auth

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type auth struct{}

func NewAuth() *auth {
	return &auth{}
}

func (a *auth) CompareHashAndPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}

	return nil
}

// GeneratePasswordHash generates a bcrypt hash of the given password using the specified cost.
func (a *auth) GeneratePasswordHash(password string) (string, error) {
	cost := bcrypt.DefaultCost
	envCost := os.Getenv("BCRYPT_GENERATE_FROM_PASSWORD_COST")
	if len(envCost) > 0 {
		costInt, err := strconv.Atoi(envCost)
		if err != nil {
			return "", fmt.Errorf("invalid cost value: %v", err)
		}

		cost = costInt
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	return string(hash), nil
}
