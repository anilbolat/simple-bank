package db

import (
	"context"
	"github.com/anilbolat/simple-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomTransfer(accountFromExpected Account, accountToExpected Account, t *testing.T) Transfer {
	arg := CreateTransferParams{
		FromAccountID: accountFromExpected.ID,
		ToAccountID:   accountToExpected.ID,
		Amount:        util.RandomMoney(),
	}

	transferActual, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transferActual)

	require.Equal(t, arg.FromAccountID, transferActual.FromAccountID)
	require.Equal(t, arg.ToAccountID, transferActual.ToAccountID)
	require.Equal(t, arg.Amount, transferActual.Amount)

	require.NotZero(t, transferActual.CreatedAt)
	require.NotZero(t, transferActual.ID)

	return transferActual
}

func TestQueries_CreateTransfer(t *testing.T) {
	accountFromExpected := createRandomAccount(t)
	accountToExpected := createRandomAccount(t)
	createRandomTransfer(accountFromExpected, accountToExpected, t)
}

func TestQueries_GetTransfer(t *testing.T) {
	accountFromExpected := createRandomAccount(t)
	accountToExpected := createRandomAccount(t)
	transferExpected := createRandomTransfer(accountFromExpected, accountToExpected, t)

	transferActual, err := testQueries.GetTransfer(context.Background(), transferExpected.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transferActual)

	require.Equal(t, transferExpected.ID, transferActual.ID)
	require.Equal(t, transferExpected.FromAccountID, transferActual.FromAccountID)
	require.Equal(t, transferExpected.ToAccountID, transferActual.ToAccountID)
	require.Equal(t, transferExpected.Amount, transferActual.Amount)
	require.Equal(t, transferExpected.CreatedAt, transferActual.CreatedAt)
}

func TestQueries_ListTransfers(t *testing.T) {
	accountFromExpected := createRandomAccount(t)
	accountToExpected := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(accountFromExpected, accountToExpected, t)
		createRandomTransfer(accountToExpected, accountFromExpected, t)
	}

	arg := ListTransfersParams{
		FromAccountID: accountFromExpected.ID,
		ToAccountID:   accountFromExpected.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == accountFromExpected.ID || transfer.ToAccountID == accountFromExpected.ID)
	}
}
