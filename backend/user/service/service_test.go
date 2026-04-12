package service

import (
	"context"
	user "likexuser/model"
	"likexuser/repository"
	"os"
	"testing"

	"github.com/qosdil/like-x/backend/common/http/auth"
	likexService "github.com/qosdil/like-x/backend/common/service"
	"golang.org/x/crypto/bcrypt"
)

type mockRepository struct {
	createOutput user.ID
	createErr    error
	firstHash    string
	firstErr     error
	firstID      user.ID
	firstIDErr   error
}

func (m *mockRepository) Create(ctx context.Context, input repository.CreateInput) (user.ID, error) {
	return m.createOutput, m.createErr
}

func (m *mockRepository) FirstPasswordHashByPublicID(ctx context.Context, publicID user.PublicID) (string, error) {
	return m.firstHash, m.firstErr
}

func (m *mockRepository) FirstIDByPublicID(ctx context.Context, publicID user.PublicID) (user.ID, error) {
	return m.firstID, m.firstIDErr
}

type fakeAuthenticator struct {
	token string
	err   error
}

func (f fakeAuthenticator) CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (f fakeAuthenticator) GeneratePasswordHash(password string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return "hash", nil
}

func (f fakeAuthenticator) GenerateToken(_ string) (string, error) {
	return f.token, f.err
}

func TestSignUp_Valid(t *testing.T) {
	m := &mockRepository{createOutput: user.ID(1)}
	svc := NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, m)

	out, err := svc.SignUp(context.Background(), user.CreateInput{FullName: "John Doe", Password: "secret123"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if out.PublicID == "" {
		t.Fatal("expected non-empty public id")
	}
	if out.ID != user.ID(1) {
		t.Fatalf("expected ID=1, got %d", out.ID)
	}
}

func TestSignUp_InvalidFullName(t *testing.T) {
	m := &mockRepository{}
	svc := NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, m)

	_, err := svc.SignUp(context.Background(), user.CreateInput{FullName: "Jan", Password: "secret123"})
	if err != likexService.ErrBadRequest {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestSignUp_InvalidPassword(t *testing.T) {
	m := &mockRepository{}
	svc := NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, m)

	_, err := svc.SignUp(context.Background(), user.CreateInput{FullName: "John Doe", Password: "123"})
	if err != likexService.ErrBadRequest {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestSignUp_RepositoryError(t *testing.T) {
	m := &mockRepository{createErr: likexService.ErrInternal}
	svc := NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, m)

	_, err := svc.SignUp(context.Background(), user.CreateInput{FullName: "John Doe", Password: "secret123"})
	if err != likexService.ErrInternal {
		t.Fatalf("expected ErrInternal, got %v", err)
	}
}

func TestAuthenticate_Success(t *testing.T) {
	pwHash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	m := &mockRepository{firstHash: string(pwHash)}
	svc := NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, m)

	out, err := svc.Authenticate(context.Background(), user.AuthInput{PublicID: "pub-1", Password: "secret123"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if out.Token != "token" {
		t.Fatalf("expected token 'token', got '%s'", out.Token)
	}
}

func TestAuthenticate_NotFound(t *testing.T) {
	m := &mockRepository{firstErr: likexService.ErrNotFound}
	svc := NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, m)

	_, err := svc.Authenticate(context.Background(), user.AuthInput{PublicID: "pub-1", Password: "secret123"})
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthenticate_InvalidPassword(t *testing.T) {
	pwHash, _ := bcrypt.GenerateFromPassword([]byte("otherpass"), bcrypt.DefaultCost)
	m := &mockRepository{firstHash: string(pwHash)}
	svc := NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, m)

	_, err := svc.Authenticate(context.Background(), user.AuthInput{PublicID: "pub-1", Password: "secret123"})
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthenticateInternal_Success(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "supersecretkeythatisatleast32characterslong")
	defer os.Unsetenv("JWT_SECRET_KEY")

	token, err := auth.GenerateJWT("pub-1")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	m := &mockRepository{firstID: 42}
	svc := NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, m)

	out, err := svc.AuthenticateInternal(context.Background(), token)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if out.ID != 42 {
		t.Fatalf("expected ID 42, got %d", out.ID)
	}
}

func TestAuthenticateInternal_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "supersecretkeythatisatleast32characterslong")
	defer os.Unsetenv("JWT_SECRET_KEY")

	m := &mockRepository{firstID: 42}
	svc := NewService(fakeAuthenticator{token: "token"}, fakeAuthenticator{token: "token"}, m)

	_, err := svc.AuthenticateInternal(context.Background(), "bad.token.value")
	if err != auth.ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}
