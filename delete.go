package search

import (
	"context"
	"time"
)

// Delete removes a record from the search table by user ID
func Delete(
	UserID uint64,
) error {

	// Check if core is initialized
	if core == nil || core.sql == nil {
		return ErrNotInitialize
	}

	// Construct the SQL query for deletion
	query := "DELETE FROM search WHERE user = ?"

	// Use the passed context to ensure consistent timeout management
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	// Execute the SQL deletion
	_, err := core.sql.ExecContext(ctx, query, UserID)
	if err != nil {
		return err
	}

	return nil
}
