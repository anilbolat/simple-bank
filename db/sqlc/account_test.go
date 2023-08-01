package db

import (
	"context"
	"database/sql"
	"github.com/anilbolat/simple-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestQueries_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestQueries_GetAccount(t *testing.T) {
	accountExpected := createRandomAccount(t)
	accountActual, err := testQueries.GetAccount(context.Background(), accountExpected.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountActual)

	require.Equal(t, accountExpected.ID, accountActual.ID)
	require.Equal(t, accountExpected.Owner, accountActual.Owner)
	require.Equal(t, accountExpected.Balance, accountActual.Balance)
	require.Equal(t, accountExpected.Currency, accountActual.Currency)
	require.Equal(t, accountExpected.CreatedAt, accountActual.CreatedAt)
	require.WithinDuration(t, accountExpected.CreatedAt, accountActual.CreatedAt, time.Second)
}

func TestQueries_UpdateAccount(t *testing.T) {
	accountExpected := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      accountExpected.ID,
		Balance: util.RandomMoney(),
	}

	accountActual, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accountActual)

	require.Equal(t, accountExpected.ID, accountActual.ID)
	require.Equal(t, accountExpected.Owner, accountActual.Owner)
	require.Equal(t, arg.Balance, accountActual.Balance)
	require.Equal(t, accountExpected.Currency, accountActual.Currency)
	require.Equal(t, accountExpected.CreatedAt, accountActual.CreatedAt)
	require.WithinDuration(t, accountExpected.CreatedAt, accountActual.CreatedAt, time.Second)
}

func TestQueries_DeleteAccount(t *testing.T) {
	accountExpected := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), accountExpected.ID)
	require.NoError(t, err)

	accountActual, err := testQueries.GetAccount(context.Background(), accountExpected.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountActual)
}

func TestQueries_ListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accountsActual, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accountsActual, 5)

	for _, account := range accountsActual {
		require.NotEmpty(t, account)
	}

}
