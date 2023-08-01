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
		transfer := result.Transfer
		require.NotEmpty(t, result)
		require.Equal(t, accountFrom.ID, transfer.FromAccountID)
		require.Equal(t, accountTo.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		entryFrom := result.FromEntry
		require.NotEmpty(t, entryFrom)
		require.Equal(t, accountFrom.ID, entryFrom.AccountID)
		require.Equal(t, -amount, entryFrom.Amount)
		require.NotZero(t, entryFrom.ID)
		require.NotZero(t, entryFrom.CreatedAt)

		_, err = store.GetEntry(context.Background(), entryFrom.ID)
		require.NoError(t, err)

		entryTo := result.ToEntry
		require.NotEmpty(t, entryTo)
		require.Equal(t, accountTo.ID, entryTo.AccountID)
		require.Equal(t, amount, entryTo.Amount)
		require.NotZero(t, entryTo.ID)
		require.NotZero(t, entryTo.CreatedAt)

		_, err = store.GetEntry(context.Background(), entryTo.ID)
		require.NoError(t, err)

		// TODO: check accounts' balances

	}
}
