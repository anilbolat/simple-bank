package db

import (
	"context"
	"database/sql"
)

// Store provides all funcs to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, queryFn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(tx)
	err = queryFn(queries)
	if err != nil {
		if errRb := tx.Rollback(); errRb != nil {
			return errRb
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, an entry record, and update accounts' balances within a single db tx
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(queries *Queries) error {
		var err error

		// create transfer
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		// create entry for 'the from account'
		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// create entry for 'the to account'
		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// update balance for 'the from account'
		accountFrom, err := queries.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		result.FromAccount, err = queries.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: accountFrom.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		// update balance for 'the to account'
		accountTo, err := queries.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		result.ToAccount, err = queries.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: accountTo.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return result, err
	}

	return result, nil
}
