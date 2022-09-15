package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier	// allow access to all the methods which use *Queries
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	// Composition - way to extend struct functionality in Go instead of inheritance
	// All the individual query functions provided by Queries will be available to Store
	*Queries

	// To manage DB transaction
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function with a taken generic db transaction
// Start a new db transaction then creating a new query with that transaction
// Call the callback function with created query
// Commit or rollback transaction based on the returned error by that function
// Private function: provide an exported function for each specific transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// To get the transaction name from the input context of the TransferTx() function
// Context background key should be type struct => use empty struct
// var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, update account's balance within a single database transaction\
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Get value in context with the key is txKey
		// txName := ctx.Value(txKey)

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, result.ToEntry, err = createEntries(ctx, q, arg.FromAccountID, arg.ToAccountID, arg.Amount)
		if err != nil {
			return err
		}

		// Get account => update its balance

		/*
			//! Use 2 queries for getting and updating is not too good
			//! => Use only 1 query for changing the account balance
				account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountId)
				if err != nil {
					return err
				}

				result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
					ID:      arg.FromAccountId,
					Balance: account1.Balance - arg.Amount,
				})
				if err != nil {
					return err
				}
		*/

		// Avoid DB deadlock : query order matter
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func createEntries(
	ctx context.Context,
	q *Queries,
	fromAccountID int64,
	toAccountID int64,
	amount int64,
) (fromEntry, toEntry Entry, err error) {
	fromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
		AccountID: fromAccountID,
		Amount:    -amount,
	})
	if err != nil {
		return
	}

	toEntry, err = q.CreateEntry(ctx, CreateEntryParams{
		AccountID: toAccountID,
		Amount:    amount,
	})
	
	return
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		// the same as: return account1, account2, err
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}
