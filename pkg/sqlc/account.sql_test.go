package sqlc

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    gofakeit.Name(),
		Balance:  int64(gofakeit.Number(0, 1000)),
		Currency: gofakeit.CurrencyShort(),
	}

	id, err := testQueries.CreateAccount(context.Background(), arg)
	assert.NoError(t, err)
	got, err := testQueries.GetAccount(context.Background(), id)
	assert.NoError(t, err)
	assert.NotEmpty(t, got)

	assert.Equal(t, arg.Owner, got.Owner)
	assert.Equal(t, arg.Balance, got.Balance)
	assert.Equal(t, arg.Currency, got.Currency)
	assert.NotZero(t, got.ID)
	assert.NotZero(t, got.CreatedAt)

	return got
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: int64(gofakeit.Number(0, 1000)),
	}

	err := testQueries.UpdateAccount(context.Background(), arg)
	assert.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, account2)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, arg.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	assert.NoError(t, err)
	assert.Len(t, accounts, 5)

	for _, account := range accounts {
		assert.NotEmpty(t, account)
	}
}
