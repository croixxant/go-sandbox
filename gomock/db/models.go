// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
)

type Post struct {
	ID         int64
	UserID     int64
	Body       string
	LikesCount int32
}

type User struct {
	ID             int64
	Email          string
	HashedPassword string
	ConfirmedAt    sql.NullTime
	LikesCount     int32
}