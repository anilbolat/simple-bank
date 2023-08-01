package db

import (
	"context"
	"github.com/anilbolat/simple-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomEntry(account Account, t *testing.T) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entryActual, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entryActual)
	require.Equal(t, account.ID, entryActual.AccountID)
	require.Equal(t, arg.Amount, entryActual.Amount)
	require.WithinDuration(t, account.CreatedAt, entryActual.CreatedAt, time.Second)

	return entryActual
}

func TestQueries_CreateEntry(t *testing.T) {
	accountExpected := createRandomAccount(t)
	createRandomEntry(accountExpected, t)
}

func TestQueries_GetEntry(t *testing.T) {
	accountExpected := createRandomAccount(t)
	entryExpected := createRandomEntry(accountExpected, t)

	entryActual, err := testQueries.GetEntry(context.Background(), entryExpected.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryActual)

	require.Equal(t, entryExpected.AccountID, entryActual.AccountID)
	require.Equal(t, entryExpected.Amount, entryActual.Amount)
	require.Equal(t, entryExpected.ID, entryActual.ID)
	require.WithinDuration(t, entryExpected.CreatedAt, entryActual.CreatedAt, time.Second)
}

func TestQueries_ListEntries(t *testing.T) {

	// given
	accountExpected := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(accountExpected, t)
	}

	arg := ListEntriesParams{
		AccountID: accountExpected.ID,
		Limit:     5,
		Offset:    5,
	}

	// test
	entryArr, err := testQueries.ListEntries(context.Background(), arg)

	// assert
	require.NoError(t, err)
	require.NotEmpty(t, entryArr)
	require.Len(t, entryArr, 5)

	for _, entry := range entryArr {
		require.NotEmpty(t, entry)
		require.Equal(t, accountExpected.ID, entry.AccountID)
	}

}
