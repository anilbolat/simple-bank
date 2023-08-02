package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore_TransferTx(t *testing.T) {
	ctx := context.Background()
	store := NewStore(testDB)

	accountFromInit := createRandomAccount(t)
	accountToInit := createRandomAccount(t)
	amount := int64(10)
	transferTxParams := TransferTxParams{
		FromAccountID: accountFromInit.ID,
		ToAccountID:   accountToInit.ID,
		Amount:        amount,
	}

	errs := make(chan error)
	results := make(chan TransferTxResult)
	n := 5
	// run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		go func() {
			transferTxResult, err := store.TransferTx(ctx, transferTxParams)
			errs <- err
			results <- transferTxResult
		}()
	}

	// check transferTxResults
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		assertTransfer(t, result.Transfer, accountFromInit, accountToInit, amount)
		_, err = store.GetTransfer(ctx, result.Transfer.ID)
		require.NoError(t, err)

		// check entries
		assertEntry(t, result.FromEntry, accountFromInit, -amount)
		_, err = store.GetEntry(ctx, result.FromEntry.ID)
		require.NoError(t, err)

		assertEntry(t, result.ToEntry, accountToInit, amount)
		_, err = store.GetEntry(ctx, result.ToEntry.ID)
		require.NoError(t, err)

		// check accounts in transfer obj
		accountFrom := result.FromAccount
		require.NotEmpty(t, accountFrom)
		require.Equal(t, accountFromInit.ID, accountFrom.ID)

		accountTo := result.ToAccount
		require.NotEmpty(t, accountTo)
		require.Equal(t, accountToInit.ID, accountTo.ID)

		// check accounts' balances
		diff1 := accountFromInit.Balance - accountFrom.Balance
		diff2 := accountTo.Balance - accountToInit.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check updated balances
	updatedAccountFrom, err := store.GetAccount(ctx, accountFromInit.ID)
	require.NoError(t, err)
	updatedAccountTo, err := store.GetAccount(ctx, accountToInit.ID)
	require.NoError(t, err)

	require.Equal(t, updatedAccountFrom.Balance, accountFromInit.Balance-int64(n)*amount)
	require.Equal(t, updatedAccountTo.Balance, accountToInit.Balance+int64(n)*amount)
}

func assertEntry(t *testing.T, entry Entry, account Account, amount int64) {
	require.NotEmpty(t, entry)
	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func assertTransfer(t *testing.T, transfer Transfer, accountFrom Account, accountTo Account, amount int64) {
	require.NotEmpty(t, transfer)
	require.Equal(t, accountFrom.ID, transfer.FromAccountID)
	require.Equal(t, accountTo.ID, transfer.ToAccountID)
	require.Equal(t, amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}
