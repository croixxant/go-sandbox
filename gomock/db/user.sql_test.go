package db_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/croixxant/golang-examples/gomock/db"
)

func createRandomUser(t *testing.T) sql.Result {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	arg := db.CreateUserParams{
		Email:          gofakeit.Email(),
		HashedPassword: string(hashed),
	}

	result, err := testQueries.CreateUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	return result
}

func TestQueries_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestQueries_DeleteUser(t *testing.T) {
	result := createRandomUser(t)
	id, _ := result.LastInsertId()

	err := testQueries.DeleteUser(context.Background(), id)
	assert.NoError(t, err)

	user, err := testQueries.GetUser(context.Background(), id)
	assert.Error(t, err)
	assert.Empty(t, user)
}

func TestQueries_GetUser(t *testing.T) {
	result := createRandomUser(t)
	id, _ := result.LastInsertId()

	user, err := testQueries.GetUser(context.Background(), id)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)
	assert.True(t, id == user.ID)
}

func TestQueries_ListUsers(t *testing.T) {
	result := createRandomUser(t)
	id, _ := result.LastInsertId()

	users, err := testQueries.ListUsers(context.Background())
	assert.NoError(t, err)
	var exist bool
	for _, user := range users {
		if id == user.ID {
			exist = true
			break
		}
	}
	assert.True(t, exist)
}
