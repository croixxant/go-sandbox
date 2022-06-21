package db_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"

	"github.com/croixxant/go-sandbox/sqlc/db"
)

func createRandomPost(t *testing.T, userID int64) sql.Result {
	arg := db.CreatePostParams{
		UserID: userID,
		Body:   gofakeit.LoremIpsumWord(),
	}

	result, err := testQueries.CreatePost(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	return result
}

func TestQueries_CreatePost(t *testing.T) {
	createUserResult := createRandomUser(t)
	userID, _ := createUserResult.LastInsertId()

	createRandomPost(t, userID)
}

func TestQueries_AddPostLikesCount(t *testing.T) {
	createUserResult := createRandomUser(t)
	userID, _ := createUserResult.LastInsertId()

	result := createRandomPost(t, userID)
	id, _ := result.LastInsertId()

	arg := db.AddPostLikesCountParams{
		Likes: int32(gofakeit.Number(1, 10)),
		ID:    id,
	}

	err := testQueries.AddPostLikesCount(context.Background(), arg)
	assert.NoError(t, err)
}

func TestQueries_DeletePost(t *testing.T) {
	createUserResult := createRandomUser(t)
	userID, _ := createUserResult.LastInsertId()

	result := createRandomPost(t, userID)
	id, _ := result.LastInsertId()

	err := testQueries.DeletePost(context.Background(), id)
	assert.NoError(t, err)

	post, err := testQueries.GetPost(context.Background(), id)
	assert.Error(t, err)
	assert.Empty(t, post)

}

func TestQueries_GetPost(t *testing.T) {
	createUserResult := createRandomUser(t)
	userID, _ := createUserResult.LastInsertId()

	result := createRandomPost(t, userID)
	id, _ := result.LastInsertId()

	post, err := testQueries.GetPost(context.Background(), id)
	assert.NoError(t, err)
	assert.NotEmpty(t, post)
	assert.True(t, id == post.ID)
}

func TestQueries_GetUserPosts(t *testing.T) {
	createUserResult := createRandomUser(t)
	userID, _ := createUserResult.LastInsertId()

	result := createRandomPost(t, userID)
	id, _ := result.LastInsertId()

	posts, err := testQueries.GetUserPosts(context.Background(), userID)
	assert.NoError(t, err)
	var exist bool
	for _, post := range posts {
		if id == post.ID {
			exist = true
			break
		}
	}
	assert.True(t, exist)
}
