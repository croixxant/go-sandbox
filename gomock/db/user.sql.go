// Code generated by sqlc. DO NOT EDIT.
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const addUserLikesCount = `-- name: AddUserLikesCount :exec
UPDATE users
SET likes_count = likes_count + ?
WHERE id = ?
`

type AddUserLikesCountParams struct {
	Likes int32
	ID    int64
}

func (q *Queries) AddUserLikesCount(ctx context.Context, arg AddUserLikesCountParams) error {
	_, err := q.db.ExecContext(ctx, addUserLikesCount, arg.Likes, arg.ID)
	return err
}

const createUser = `-- name: CreateUser :execresult
INSERT INTO users (
  email, hashed_password
) VALUES (
  ?, ?
)
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser, arg.Email, arg.HashedPassword)
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, email, hashed_password, confirmed_at, likes_count FROM users
WHERE id = ? LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.ConfirmedAt,
		&i.LikesCount,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, email, hashed_password, confirmed_at, likes_count FROM users
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.HashedPassword,
			&i.ConfirmedAt,
			&i.LikesCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
