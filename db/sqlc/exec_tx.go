package db

import (
	"context"
	"fmt"
)

// execTx executes a function with a taken generic db transaction
// Start a new db transaction then creating a new query with that transaction
// Call the callback function with created query
// Commit or rollback transaction based on the returned error by that function
// Private function: provide an exported function for each specific transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit(ctx)
}
