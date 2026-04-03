package repository

import (
	"context"
	"errors"
	user "likexuser/model"
	"testing"

	"github.com/jackc/pgx/v5"
	likexService "github.com/qosdil/like-x/backend/common/service"
)

type fakeRow struct {
	id       user.ID
	password string
	err      error
}

func (r fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != 1 {
		return errors.New("expected one destination")
	}
	if p, ok := dest[0].(*user.ID); ok {
		*p = r.id
		return nil
	}
	if p, ok := dest[0].(*string); ok {
		*p = r.password
		return nil
	}
	return errors.New("invalid destination type")
}

type fakeConn struct {
	row fakeRow
}

func (c fakeConn) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.row
}

// NewDbForTest allows injecting a mock query-row connection for unit tests.
func NewDbForTest(conn queryRowConn) Repository {
	return &db{conn: conn}
}

func TestPgx_Create_Success(t *testing.T) {
	conn := fakeConn{row: fakeRow{id: 42}}
	repo := NewDbForTest(conn)

	out, err := repo.Create(context.Background(), user.CreateInput{FullName: "John Doe", Password: "secret123"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if out.ID != 42 {
		t.Fatalf("expected ID=42, got %v", out.ID)
	}
	if out.PublicID == "" {
		t.Fatal("expected non-empty PublicID")
	}
}

func TestPgx_Create_QueryError(t *testing.T) {
	conn := fakeConn{row: fakeRow{err: pgx.ErrNoRows}}
	repo := NewDbForTest(conn)

	_, err := repo.Create(context.Background(), user.CreateInput{FullName: "John Doe", Password: "secret123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPgx_FirstPasswordHashByPublicID_Success(t *testing.T) {
	conn := fakeConn{row: fakeRow{password: "hash123"}}
	repo := NewDbForTest(conn)

	got, err := repo.FirstPasswordHashByPublicID(context.Background(), "pub-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != "hash123" {
		t.Fatalf("expected hash123, got %s", got)
	}
}

func TestPgx_FirstPasswordHashByPublicID_NotFound(t *testing.T) {
	conn := fakeConn{row: fakeRow{err: pgx.ErrNoRows}}
	repo := NewDbForTest(conn)

	_, err := repo.FirstPasswordHashByPublicID(context.Background(), "pub-1")
	if err != likexService.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPgx_FirstPasswordHashByPublicID_Error(t *testing.T) {
	conn := fakeConn{row: fakeRow{err: errors.New("db fail")}}
	repo := NewDbForTest(conn)

	_, err := repo.FirstPasswordHashByPublicID(context.Background(), "pub-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
