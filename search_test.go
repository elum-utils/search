package search

import (
	"testing"
	"time"
)

func TestDuplicate(t *testing.T) {

	SetTimeout(5 * time.Second)

	t.Run("successful search after creation", func(t *testing.T) {
		Create(12, "en", 18, 30, 0, 25, 1, true, "music")

		user := Search(13, "en", 18, 30, 1, 25, 0, "music")
		if user == nil {
			t.Error("Should find user after creation")
		}
	})

	t.Run("duplicate creation prevents search", func(t *testing.T) {
		Create(12, "en", 18, 30, 0, 25, 1, true, "music")

		user := Search(13, "en", 18, 30, 1, 25, 0, "music")
		if user != nil {
			t.Error("Should not find user after duplicate creation")
		}
	})

	t.Run("successful search after timeout", func(t *testing.T) {
		time.Sleep(6 * time.Second)

		Create(12, "en", 18, 30, 0, 25, 1, true, "music")

		user := Search(13, "en", 18, 30, 1, 25, 0, "music")
		if user == nil {
			t.Error("Should find user after timeout expiration")
		}
	})

}

func TestCreateAndDelete(t *testing.T) {

	SetTimeout(0)
	// Create a record
	Create(10, "ru", 20, 40, 1, 30, 0, false, "music", "art")
	if _, exists := store.entries[10]; !exists {
		t.Error("Record was not created")
	}

	// Delete the record
	Delete(10)
	if _, exists := store.entries[10]; exists {
		t.Error("Record was not deleted")
	}
}

func TestSearchAllConditions(t *testing.T) {

	SetTimeout(0)

	testCases := []struct {
		name     string
		setup    func()
		params   []any
		expected uint64
	}{
		{
			name: "Does not match",
			setup: func() {
				Create(1, "en", 18, 24, 0, 25, 1, false)
			},
			params:   []any{2, "en", 45, 100, 1, 39, 0},
			expected: 0,
		},
		{
			name: "Language does not match",
			setup: func() {
				Create(1, "en", 18, 30, 0, 25, 1, false, "music")
			},
			params:   []any{2, "ru", 18, 30, 1, 25, 0, "music"},
			expected: 0,
		},
		{
			name: "Age is outside of preferred range",
			setup: func() {
				Create(1, "en", 18, 30, 0, 35, 1, false, "music")
				Create(2, "en", 18, 30, 0, 17, 1, false, "music")
			},
			params:   []any{3, "en", 18, 30, 1, 40, 0, "music"},
			expected: 2,
		},
		{
			name: "Sex does not match",
			setup: func() {
				Create(3, "en", 18, 30, 1, 25, 1, false, "music")
			},
			params:   []any{4, "en", 18, 30, 0, 25, 0, "music"},
			expected: 0,
		},
		{
			name: "Match by interests",
			setup: func() {
				Create(4, "en", 18, 30, 1, 25, 0, false, "music", "art")
			},
			params:   []any{5, "en", 18, 30, 0, 25, 1, "music"},
			expected: 4,
		},
		{
			name: "Best match by number of interests",
			setup: func() {
				Create(6, "en", 18, 30, 0, 25, 1, false, "music")
				Create(7, "en", 18, 30, 1, 25, 0, false, "music", "art")
			},
			params:   []any{8, "en", 18, 30, 0, 25, 1, "music", "art"},
			expected: 7,
		},
		{
			name: "Any sex (2) is accepted",
			setup: func() {
				Create(9, "en", 18, 30, 0, 25, 1, false, "music")
			},
			params:   []any{10, "en", 18, 30, 2, 25, 0, "music"},
			expected: 9,
		},
		{
			name: "Any sex (2) rejected because user requires specific sex",
			setup: func() {
				Create(10, "en", 18, 30, 1, 25, 1, false, "music")
			},
			params:   []any{11, "en", 18, 30, 2, 25, 0, "music"},
			expected: 0,
		},
		{
			name: "Priority affects score",
			setup: func() {
				Create(11, "en", 18, 30, 0, 25, 1, false, "music")
				Create(12, "en", 18, 30, 0, 25, 1, true, "music")
			},
			params:   []any{13, "en", 18, 30, 1, 25, 0, "music"},
			expected: 12,
		},
		{
			name: "User cannot match with self",
			setup: func() {
				Create(12, "en", 18, 30, 1, 25, 1, false, "music")
			},
			params:   []any{12, "en", 18, 30, 1, 25, 1, "music"},
			expected: 0,
		},
		{
			name: "Explicit mismatch: YourSex != user.MySex",
			setup: func() {
				// user.MySex = 0 (male)
				Create(201, "en", 18, 30, 1, 25, 0, false, "music")
			},
			params:   []any{202, "en", 18, 30, 1, 25, 1, "music"}, // expects female
			expected: 0,
		},
		{
			name: "Score < 0 leads to zero return",
			setup: func() {
				// user expects 18–20, seeker is 50 => diff = 30
				Create(203, "en", 18, 20, 1, 25, 0, false, "music")
			},
			params:   []any{204, "en", 50, 55, 1, 50, 1, "music"},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset store for each test
			Close()
			defer Close()

			if tc.setup != nil {
				tc.setup()
			}

			// Extract parameters
			myID := uint64(tc.params[0].(int))
			lang := tc.params[1].(string)
			yourStart := tc.params[2].(int)
			yourEnd := tc.params[3].(int)
			yourSex := tc.params[4].(int)
			myAge := tc.params[5].(int)
			mySex := tc.params[6].(int)

			var interests []string
			if len(tc.params) > 8 {
				for _, p := range tc.params[7:] {
					interests = append(interests, p.(string))
				}
			}

			result := Search(myID, lang, yourStart, yourEnd, yourSex, myAge, mySex, interests...)

			if tc.expected == 0 {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("Expected UserID %d, got nil", tc.expected)
				} else if result.UserID != tc.expected {
					t.Errorf("Expected UserID %d, got %d", tc.expected, result.UserID)
				}
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {

	SetTimeout(0)
	// User with no interests
	Create(11, "en", 18, 30, 1, 25, 0, false)
	if len(store.entries[11].Interests) != 0 {
		t.Error("Expected empty interests")
	}

	// Search with no interests
	result := Search(13, "en", 18, 30, 1, 25, 0)
	if result != nil {
		t.Error("Expected nil when searching with no interests")
	}
}

func TestRaceConditions(t *testing.T) {

	Close()
	defer Close()

	SetTimeout(0)
	defer SetTimeout(0)

	done := make(chan bool)

	// 10 goroutines creating records
	for i := 0; i < 10; i++ {
		go func(id int) {
			Create(uint64(id), "en", 20, 30, 1, 25, 0, false, "music")
			done <- true
		}(i)
	}

	// 10 goroutines deleting records
	for i := 0; i < 10; i++ {
		go func(id int) {
			Delete(uint64(id))
			done <- true
		}(i)
	}

	// 10 goroutines searching
	for i := 0; i < 10; i++ {
		go func(id int) {
			Search(uint64(100+id), "en", 18, 30, 1, 25, 0, "music")
			done <- true
		}(i)
	}

	// Wait for all routines to complete
	for i := 0; i < 30; i++ {
		<-done
	}
}

// Benchmark how fast we can create N users
func BenchmarkCreate(b *testing.B) {
	SetTimeout(0)
	defer Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Create(uint64(i), "en", 18, 30, i%3, 25, i%2, false, "music", "art", "games")
	}
}

// Benchmark search performance with only 1 user in the system
func BenchmarkSearchSingleUser(b *testing.B) {
	SetTimeout(0)
	defer Close()

	Create(1, "en", 18, 30, 1, 25, 0, false, "music", "art")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Search(2, "en", 18, 30, 1, 25, 0, "music", "art")
	}
}

// Benchmark search with many users (1000)
func BenchmarkSearchMultipleUsers(b *testing.B) {
	SetTimeout(0)
	defer Close()

	for i := 0; i < 1000; i++ {
		Create(uint64(i), "en", 18+i%10, 30+i%10, i%3, 20+i%10, i%2,
			false, "music", "art", "games", "science")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Search(1001, "en", 18, 30, 1, 25, 0, "music", "art")
	}
}

// Benchmark with many interests on both sides
func BenchmarkSearchWithManyInterests(b *testing.B) {
	SetTimeout(0)
	defer Close()

	for i := 0; i < 1000; i++ {
		interests := []string{"music", "art"}
		if i%2 == 0 {
			interests = append(interests, "games")
		}
		if i%3 == 0 {
			interests = append(interests, "science")
		}
		if i%5 == 0 {
			interests = append(interests, "sports")
		}

		Create(uint64(i), "en", 18+i%10, 30+i%10, i%3, 20+i%10, i%2, false, interests...)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Search(1001, "en", 18, 30, 1, 25, 0, "music", "art", "games", "science", "sports")
	}
}

// Benchmark performance when searching users with different languages
func BenchmarkSearchWithDifferentLanguages(b *testing.B) {
	SetTimeout(0)
	defer Close()

	for i := 0; i < 1000; i++ {
		lang := "en"
		if i%2 == 0 {
			lang = "fr"
		}
		if i%3 == 0 {
			lang = "de"
		}
		Create(uint64(i), lang, 18+i%10, 30+i%10, i%3, 20+i%10, i%2, false, "music", "art")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Search(1001, "en", 18, 30, 1, 25, 0, "music", "art")
	}
}
