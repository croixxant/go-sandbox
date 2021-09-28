// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	AddPostLikesCount(ctx context.Context, arg AddPostLikesCountParams) error
	AddUserLikesCount(ctx context.Context, arg AddUserLikesCountParams) error
	CreatePost(ctx context.Context, arg CreatePostParams) (sql.Result, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
	DeletePost(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetPost(ctx context.Context, id int64) (Post, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserPosts(ctx context.Context, userID int64) ([]Post, error)
	ListUsers(ctx context.Context) ([]User, error)
}

var _ Querier = (*Queries)(nil)
