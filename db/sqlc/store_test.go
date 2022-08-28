package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
	Deadlock reason:
		1. Concurrent requests change or retrieve the db records at the same time may get the same balance, not the true balance
			=> SELECT ... FOR UPDATE: allow only one transaction at a time to interact with the db record
	 	2. Multiple requests interact with the same record with foreign key constraints, including inserting, and updating other tables having columns are foreign keys
			=> FOR NO KEY UPDATE: tell the db that other transactions retrieve or change the record column which is not a key
*/
func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// Channels - share data between channels without explicit locking
	errs := make(chan error)
	results := make(chan TransferTxResult)

	// Run n concurrent transfer transactions
	n := 6
	amount := int64(10)
	for index := 0; index < n; index++ {
		// txName := fmt.Sprintf("tx %d", index+1)

		// In the fact, transactions occur concurrently => use goroutines
		go func() {
			// Pass a value into context
			// ctx := context.WithValue(context.Background(), txKey, txName)
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			// Return error will exit the program => use channels
			errs <- err
			results <- result
		}()
	}

	// Check the results
	existed := make(map[int]bool)
	for index := 0; index < n; index++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check the transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// Check the record
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check the entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// Check accounts's balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

/*
	Deadlock reason:
		1. Query order - Circular waiting for the ShareLock (The different order in which 2 concurrent transactions update the same account's balance)
			The transaction 1 updates account 1 before account 2 while the other transaction updates account 2 before account 1
*/
func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// Channels - share data between channels without explicit locking
	errs := make(chan error)

	// Run n concurrent transfer transactions
	n := 10
	amount := int64(10)
	for index := 0; index < n; index++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if index%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		// In the fact, transactions occur concurrently => use goroutines
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			// Return error will exit the program => use channels
			errs <- err
		}()
	}

	// Check the results
	for index := 0; index < n; index++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
