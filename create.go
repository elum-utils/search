package search

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Create inserts a new record into the search table
func Create(
	UserID uint64,
	Language string,
	YourStart int,
	YourEnd int,
	YourSex int,
	MyAge int,
	MySex int,
	interests ...string,
) error {
    // Check if core and its SQL client are initialized
    if core == nil || core.sql == nil {
        return ErrNotInitialize
    }

    // Base columns
    columns := []string{
        "user", "language", "your_start", "your_end", "your_sex", "my_age", "my_sex",
    }

    // Base values
    values := []interface{}{
        UserID, Language, YourStart, YourEnd, YourSex, MyAge, MySex,
    }

    // Process and validate interests
    for _, interest := range interests {
        if core.interests[interest] { // Ensure interest is valid
            columns = append(columns, interest)
            values = append(values, 1) // Set the presence of interest to 1
        }
    }

    // Create column string and value placeholders
    columnsStr := strings.Join(columns, ", ")
    placeholders := strings.TrimSuffix(strings.Repeat("?, ", len(values)), ", ")

    // Generate update set expressions
    var updateExpressions []string
    for _, column := range columns {
        if column != "user" { // Ignore updating "user" as it's a unique key
            updateExpressions = append(updateExpressions, fmt.Sprintf("%s = excluded.%s", column, column))
        }
    }

    // Construct the SQL query with ON CONFLICT clause
    query := fmt.Sprintf(
        "INSERT INTO search (%s) VALUES (%s) ON CONFLICT(user) DO UPDATE SET %s",
        columnsStr,
        placeholders,
        strings.Join(updateExpressions, ", "),
    )

    // Execute the SQL query with context to manage execution time
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := core.sql.ExecContext(ctx, query, values...)
    if err != nil {
        return err
    }

    return nil
}
