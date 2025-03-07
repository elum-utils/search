package search

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	core *SearchSystem

	ErrNotInitialize = errors.New("the search system has not been initialized")
)

type SearchSystem struct {
	sql       *sql.DB
	interests map[string]bool
}

func New(config ...Config) error {

	cfg := configDefault(config...)

	// Remove existing database file, if any

	if cfg.Reset && cfg.LocalFile != "file::memory:" {
		err := os.Remove(cfg.LocalFile)
		if err != nil && !os.IsNotExist(err) {
			return err // Return an error if it's not "file does not exist" error
		}
	}

	// Initialize SQLite database connection
	db, err := sql.Open("sqlite3", cfg.LocalFile)
	if err != nil {
		return err
	}

	// Ping the database to ensure the connection is valid.
	err = db.Ping()
	if err != nil {
		return err // Terminate the program if the connection is invalid.
	}

	if err := create(db, cfg.Interests); err != nil {
		return err
	}

	interestsMap := make(map[string]bool, len(cfg.Interests))
	for _, interest := range cfg.Interests {
		interestsMap[interest] = true
	}

	core = &SearchSystem{
		sql:       db,
		interests: interestsMap,
	}

	return nil

}

type Params struct {
	Query string // SQL query string
	Args  []any  // Arguments for the SQL query
}

func query[T any](
	params Params,
	callback func(rows *sql.Rows) (*T, error),
) (*T, error) {

	// Create a context with a timeout for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel() // Cancel the context after the query execution

	// Execute the query with the provided arguments
	rows, err := core.sql.QueryContext(ctx, params.Query, params.Args...)
	if err != nil {
		// Return the SQL error if it is any other error
		return nil, err
	}
	defer rows.Close() // Close the rows after finishing the query

	// Call the callback function to process the rows and extract the result
	clbRes, clbErr := callback(rows)

	// Return the result and any potential MySQL error from the callback
	return clbRes, clbErr

}

func Close() error {
	return core.sql.Close()
}
