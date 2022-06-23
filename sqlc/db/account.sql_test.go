package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    gofakeit.Name(),
		Balance:  int64(gofakeit.Number(0, 1000)),
		Currency: gofakeit.CurrencyShort(),
	}

	id, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	got, err := testQueries.GetAccount(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.Equal(t, arg.Owner, got.Owner)
	require.Equal(t, arg.Balance, got.Balance)
	require.Equal(t, arg.Currency, got.Currency)
	require.NotZero(t, got.ID)
	require.NotZero(t, got.CreatedAt)

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
	require.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
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
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
