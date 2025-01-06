// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const changePassEmail = `-- name: ChangePassEmail :one
UPDATE users
SET email = $1, hashed_passwords = $2, updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, email, hashed_passwords, is_chirpy_red
`

type ChangePassEmailParams struct {
	Email           string
	HashedPasswords string
	ID              uuid.UUID
}

func (q *Queries) ChangePassEmail(ctx context.Context, arg ChangePassEmailParams) (User, error) {
	row := q.db.QueryRowContext(ctx, changePassEmail, arg.Email, arg.HashedPasswords, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPasswords,
		&i.IsChirpyRed,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_passwords)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, hashed_passwords, is_chirpy_red
`

type CreateUserParams struct {
	Email           string
	HashedPasswords string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPasswords)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPasswords,
		&i.IsChirpyRed,
	)
	return i, err
}

const deleteUsers = `-- name: DeleteUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteUsers)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, created_at, updated_at, email, hashed_passwords, is_chirpy_red FROM users
WHERE email = $1
`

func (q *Queries) GetUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPasswords,
		&i.IsChirpyRed,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, created_at, updated_at, email, hashed_passwords, is_chirpy_red FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPasswords,
		&i.IsChirpyRed,
	)
	return i, err
}

const setChirpyRed = `-- name: SetChirpyRed :one
UPDATE users
SET is_chirpy_red = true
WHERE id = $1
RETURNING id, created_at, updated_at, email, hashed_passwords, is_chirpy_red
`

func (q *Queries) SetChirpyRed(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, setChirpyRed, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPasswords,
		&i.IsChirpyRed,
	)
	return i, err
}
