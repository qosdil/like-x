package repository

import (
	"context"
	"fmt"
	user "likexuser/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	likexService "github.com/qosdil/like-x/backend/common/service"
)

type queryRowConn interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type db struct {
	conn queryRowConn
}

// Create inserts a new user into the database and returns the generated internal ID.
func (r *db) Create(ctx context.Context, input CreateInput) (id user.ID, err error) {
	sql := "INSERT INTO users (public_id, full_name, password_hash) VALUES ($1, $2, $3) RETURNING id"
	err = r.conn.QueryRow(ctx, sql, input.PublicID, input.FullName, input.PasswordHash).Scan(&id)
	if err != nil {
		return user.ID(0), fmt.Errorf("failed to create user: %v", err)
	}

	return id, nil
}

// FirstIDByPublicID retrieves the internal numeric user ID for the given public ID.
func (r *db) FirstIDByPublicID(ctx context.Context, publicID user.PublicID) (id user.ID, err error) {
	sql := "SELECT id FROM users WHERE public_id = $1"
	err = r.conn.QueryRow(ctx, sql, string(publicID)).Scan(&id)
	if err == nil {
		return
	}

	if err == pgx.ErrNoRows {
		err = likexService.ErrNotFound
		return
	}

	err = fmt.Errorf("failed to retrieve user id: %v", err)
	return
}

// FirstPasswordHashByPublicID retrieves the password hash for a user by their public ID.
func (r *db) FirstPasswordHashByPublicID(ctx context.Context, publicID user.PublicID) (passwordHash string, err error) {
	sql := "SELECT password_hash FROM users WHERE public_id = $1"
	err = r.conn.QueryRow(ctx, sql, string(publicID)).Scan(&passwordHash)
	if err == nil {
		return
	}

	if err == pgx.ErrNoRows {
		err = likexService.ErrNotFound
		return
	}

	err = fmt.Errorf("failed to retrieve password hash: %v", err)
	return
}

// NewPgx creates a Repository implementation backed by a pgx connection pool.
func NewPgx(conn *pgxpool.Pool) Repository {
	return &db{conn: conn}
}
