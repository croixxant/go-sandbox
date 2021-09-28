package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/croixxant/golang-examples/gomock/db"
)

func TestStore_LikeTx(t *testing.T) {
	store := db.NewStore(testDB)

	createUserResult := createRandomUser(t)
	userID, _ := createUserResult.LastInsertId()
	createPostResult := createRandomPost(t, userID)
	postID, _ := createPostResult.LastInsertId()
	createUser2Result := createRandomUser(t)
	likeUserID, _ := createUser2Result.LastInsertId()

	n := 10
	likesCount := 5

	errs := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			err := store.LikeTx(context.Background(), db.LikeTxParams{
				UserID:     likeUserID,
				PostID:     postID,
				LikesCount: int32(likesCount),
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)
	}

	user, _ := testQueries.GetUser(context.Background(), likeUserID)
	assert.Equal(t, int32(likesCount*n), user.LikesCount)
	post, _ := testQueries.GetPost(context.Background(), postID)
	assert.Equal(t, int32(likesCount*n), post.LikesCount)
}
