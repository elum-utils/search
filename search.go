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
	MyID uint64,
	Language string,
	YourStart int,
	YourEnd int,
	YourSex int,
	MyAge int,
	MySex int,
	interests ...string,
) (*SearchResult, error) {
	if core == nil || core.sql == nil {
		return nil, ErrNotInitialize
	}

	var queryBuilder strings.Builder
	var args []interface{}

	queryBuilder.WriteString(`
		SELECT id, user
		FROM search
		WHERE language = ?
		AND ? BETWEEN your_start AND your_end        
		AND (? = 2 OR your_sex = ? OR your_sex = 2)    
		AND my_age BETWEEN ? AND ?                     
		AND (? = 2 OR my_sex = ? OR my_sex = 2)       
		AND user != ?                                
	`)

	args = append(args,
		Language,
		MyAge,
		MySex, MySex,
		YourStart, YourEnd,
		YourSex, YourSex,
		MyID,
	)

	// Interests
	if len(interests) > 0 {
		var validInterests []string
		for _, interest := range interests {
			if core.interests[interest] {
				validInterests = append(validInterests, fmt.Sprintf("(%s = 1)", interest))
			}
		}

		if len(validInterests) > 0 {
			queryBuilder.WriteString(" AND (")
			queryBuilder.WriteString(strings.Join(validInterests, " OR "))
			queryBuilder.WriteString(")")
		} else {
			return nil, nil
		}
	}

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

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// rows, err := core.sql.QueryContext(ctx, queryBuilder.String(), args...)
	// if err != nil {
	// 	return nil, fmt.Errorf("query error: %w", err)
	// }
	// defer rows.Close()

	// if rows.Next() {
	// 	var item SearchResult
	// 	if err := rows.Scan(&item.ID, &item.UserID); err != nil {
	// 		return nil, fmt.Errorf("scan error: %w", err)
	// 	}
	// 	return &item, nil
	// }

	return nil, nil
}
