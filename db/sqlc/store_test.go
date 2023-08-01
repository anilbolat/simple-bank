package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)
	amount := int64(10)
	transferTxParams := TransferTxParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Amount:        amount,
	}

	errs := make(chan error)
	results := make(chan TransferTxResult)
	n := 5
	// run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		go func() {

			transferTxResult, err := store.TransferTx(context.Background(), transferTxParams)
			errs <- err
			results <- transferTxResult
		}()
	}

	// check transferTxResults
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		assertTransfer(t, result.Transfer, accountFrom, accountTo, amount)
		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		// check entries
		assertEntry(t, result.FromEntry, accountFrom, amount)
		_, err = store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		assertEntry(t, result.ToEntry, accountTo, amount)
		_, err = store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		// TODO: check accounts' balances

	}
}

func assertEntry(t *testing.T, entry Entry, account Account, amount int64) {
	require.NotEmpty(t, entry)
	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, -amount, entry.Amount)
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
