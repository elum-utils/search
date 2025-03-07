package search

import (
	"database/sql"
	"fmt"
	"strings"
)

type SearchResult struct {
	ID     uint64
	UserID uint64
}

func Search(
	Language string,
	YourStart int,
	YourEnd int,
	YourSex int,
	MySex int,
	MyAge int,
	interests ...string,
) (*SearchResult, error) {

	// Check if core is initialized
	if core == nil || core.sql == nil {
		return nil, ErrNotInitialize
	}

	// Use strings.Builder for constructing the query
	var queryBuilder strings.Builder

	queryBuilder.WriteString(`SELECT id, user FROM search WHERE 
      language = ?
      AND ? BETWEEN your_start AND your_end
      AND my_age BETWEEN ? AND ?
      AND (your_sex = ? OR your_sex = 2)`)

	if MySex != 2 {
		queryBuilder.WriteString("AND (my_sex = ?)")
	}

	// Populate arguments for the SQL query
	args := []any{
		Language,
		MyAge,
		YourStart,
		YourEnd,
		MySex,
		YourSex,
		MySex,
	}

	// Process interests and filter them through the map of valid interests
	if len(interests) > 0 {
		queryBuilder.WriteString(" AND (")
		first := true
		for _, interest := range interests {
			if core.interests[interest] { // Validate against the map
				if !first {
					queryBuilder.WriteString(" OR ")
				}
				fmt.Fprintf(&queryBuilder, "(%s = 1)", interest)
				first = false
			}
		}
		queryBuilder.WriteString(")")
	}

	// Finalize the SQL query with sorting and limit
	queryBuilder.WriteString(" ORDER BY priority DESC LIMIT 1")

	// Use the passed context to ensure consistent timeout management
	return query(Params{
		Query: queryBuilder.String(),
		Args:  args,
	}, func(rows *sql.Rows) (*SearchResult, error) {

		// Process the SQL result
		if rows.Next() {
			item := new(SearchResult)
			err := rows.Scan(
				&item.ID,
				&item.UserID,
			)
			if err != nil {
				return nil, err
			}
			return item, nil
		}

		// Return nil if no records are found
		return nil, nil
	})
}
