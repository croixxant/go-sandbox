package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/golang/mock/mockgen/model" // https://github.com/golang/mock#debugging-errors
)

type Store interface {
	Querier
	LikeTx(ctx context.Context, arg LikeTxParams) error
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type LikeTxParams struct {
	UserID     int64
	PostID     int64
	LikesCount int32
}

func (store *SQLStore) LikeTx(ctx context.Context, arg LikeTxParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		err = q.AddPostLikesCount(ctx, AddPostLikesCountParams{
			Likes: arg.LikesCount,
			ID:    arg.PostID,
		})
		if err != nil {
			return err
		}

		err = q.AddUserLikesCount(ctx, AddUserLikesCountParams{
			Likes: arg.LikesCount,
			ID:    arg.UserID,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
