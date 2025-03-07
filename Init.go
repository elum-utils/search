package search

import (
	"database/sql"
	"fmt"
	"strings"
)

// Init creates the required table and indexes in the database
func create(db *sql.DB, interests []string) error {
    // Define the base fields of the table
    baseFields := `
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user BIGINT UNIQUE,
        priority INT DEFAULT 0,
        language TEXT,
        your_start INT,
        your_end INT,
        your_sex TINYINT,
        my_age INT,
        my_sex TINYINT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    `

    // Prepare interest fields if provided
    var fields string
    if len(interests) > 0 {
        var interestFields []string
        for _, interest := range interests {
            interestFields = append(interestFields, fmt.Sprintf("%s TINYINT DEFAULT 0", interest))
        }
        fields = fmt.Sprintf("%s, %s", baseFields, strings.Join(interestFields, ", "))
    } else {
        fields = baseFields
    }

    createTableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS search (%s);", fields)

    // Execute table creation query
    if _, err := db.Exec(createTableQuery); err != nil {
        return err
    }

    // Create indexes for each interest field, if any
    for _, interest := range interests {
        indexQuery := fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s ON search (%s);", interest, interest)
        if _, err := db.Exec(indexQuery); err != nil {
            return err
        }
    }

    // Create additional indexes
    // Multi-column index
    indexQuery := `
        CREATE INDEX IF NOT EXISTS search_language_your_sex_your_start_your_end_created_at_index
        ON search (language, your_sex, your_start, your_end, created_at);
    `
    if _, err := db.Exec(indexQuery); err != nil {
        return err
    }

    // Index on priority
    indexQuery = `
        CREATE INDEX IF NOT EXISTS search_priority_index
        ON search (priority DESC);
    `
    if _, err := db.Exec(indexQuery); err != nil {
        return err
    }

    // Index on user and created_at
    indexQuery = `
        CREATE INDEX IF NOT EXISTS search_user_created_at_index
        ON search (user, created_at);
    `
    if _, err := db.Exec(indexQuery); err != nil {
        return err
    }

    return nil
}