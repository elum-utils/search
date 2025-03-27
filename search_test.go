package search

import (
	"testing"
)

func TestDefault(t *testing.T) {

	err := New(Config{
		// LocalFile: "_search",
		Interests: []string{
			"music", "travel", "sport", "art", "cooking", "movies", "games",
			"reading", "tech", "animals", "nature", "photography", "dance",
			"space", "science", "history", "fashion", "yoga", "psychology",
			"volunteering", "flirt", "crypto", "anime", "lgbt",
		},
	})
	if err != nil {
		t.Fatalf("Error during initialization: %v", err)
	}
	defer Close()

	Create(1, "en", 18, 30, 1, 25, 1, "music", "travel", "art")
	result, err := Search(1, "en", 18, 30, 1, 1, 25)
	if err != nil {
		t.Fatalf("Error during search: %v", err)
	}

	if result != nil {
		t.Fatalf("Error during search: %v", err)
	}

	Create(2, "en", 18, 30, 1, 25, 1, "movies", "science", "tech")
	Create(3, "en", 18, 30, 1, 25, 1, "science", "photography")

	result, err = Search(4, "en", 18, 30, 1, 1, 25, "music", "art")
	if err != nil {
		t.Fatalf("Error during search: %v", err)
	}
	if result == nil || result.UserID != 1 {
		t.Errorf("Expected to find user 1 that matches the interests 'music' and 'art', got %+v", result)
	}

}

// Test case for searching a suitable interlocutor without specific interests.
func TestSearch(t *testing.T) {
	err := New(Config{
		LocalFile: "_search",
		Interests: []string{
			"music", "travel", "sport", "art", "cooking", "movies", "games",
			"reading", "tech", "animals", "nature", "photography", "dance",
			"space", "science", "history", "fashion", "yoga", "psychology",
			"volunteering", "flirt", "crypto", "anime", "lgbt",
		},
	})
	if err != nil {
		t.Fatalf("Error during initialization: %v", err)
	}
	defer Close()

	// Assume that some entries were created previously
	Create(1, "en", 18, 30, 1, 25, 1, "music", "travel", "art")
	Create(2, "en", 18, 30, 1, 25, 1, "movies", "science", "tech")

	result, err := Search(3, "en", 18, 30, 1, 1, 25)
	if err != nil {
		t.Fatalf("Error during search: %v", err)
	}
	if result == nil || result.UserID != 1 {
		t.Errorf("Expected to find user 1, got %+v", result)
	}
}

// Test case for searching a suitable interlocutor with specific interests.
func TestSearchWithInterests(t *testing.T) {
	err := New(Config{
		Interests: []string{
			"music", "travel", "sport", "art", "cooking", "movies", "games",
			"reading", "tech", "animals", "nature", "photography", "dance",
			"space", "science", "history", "fashion", "yoga", "psychology",
			"volunteering", "flirt", "crypto", "anime", "lgbt",
		},
	})
	if err != nil {
		t.Fatalf("Error during initialization: %v", err)
	}
	defer Close()

	// Assume that some entries were created previously
	Create(1, "en", 18, 30, 1, 25, 1, "music", "travel", "art")
	Create(2, "en", 18, 30, 1, 25, 1, "movies", "science", "tech")
	Create(3, "en", 18, 30, 1, 25, 1, "science", "photography")

	result, err := Search(4, "en", 18, 30, 1, 1, 25, "music", "art")
	if err != nil {
		t.Fatalf("Error during search: %v", err)
	}
	if result == nil || result.UserID != 1 {
		t.Errorf("Expected to find user 1 that matches the interests 'music' and 'art', got %+v", result)
	}
}

// Benchmarks for the Search function with specific interests
func BenchmarkSearchWithInterests(b *testing.B) {
	err := New(Config{
		Reset: false,
		Interests: []string{
			"music", "travel", "sport", "art", "cooking", "movies", "games",
			"reading", "tech", "animals", "nature", "photography", "dance",
			"space", "science", "history", "fashion", "yoga", "psychology",
			"volunteering", "flirt", "crypto", "anime", "lgbt",
		},
	})
	if err != nil {
		b.Fatalf("Error during initialization: %v", err)
	}
	defer Close()

	// Insert records to search for
	err = Create(1, "en", 18, 30, 1, 25, 1, "music", "travel", "art")
	if err != nil {
		b.Fatalf("Error during record creation for benchmark: %v", err)
	}

	err = Create(2, "en", 18, 30, 1, 25, 1, "movies", "science", "tech")
	if err != nil {
		b.Fatalf("Error during record creation for benchmark: %v", err)
	}

	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		_, err := Search(3, "en", 18, 30, 1, 1, 25, "music", "art")
		if err != nil {
			b.Fatalf("Error during search with interests: %v", err)
		}
	}
}

// Benchmarks for the Search function without specific interests
func BenchmarkSearch(b *testing.B) {
	err := New(Config{
		Reset: false,
		Interests: []string{
			"music", "travel", "sport", "art", "cooking", "movies", "games",
			"reading", "tech", "animals", "nature", "photography", "dance",
			"space", "science", "history", "fashion", "yoga", "psychology",
			"volunteering", "flirt", "crypto", "anime", "lgbt",
		},
	})
	if err != nil {
		b.Fatalf("Error during initialization: %v", err)
	}
	defer Close()

	// Insert a record to search for
	err = Create(1, "en", 18, 30, 1, 25, 1, "music", "travel", "art")
	if err != nil {
		b.Fatalf("Error during record creation for benchmark: %v", err)
	}

	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		_, err := Search(2, "en", 18, 30, 1, 1, 25)
		if err != nil {
			b.Fatalf("Error during search: %v", err)
		}
	}
}
