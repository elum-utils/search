package search

import "testing"

// Test case for inserting a new record.
func TestCreate(t *testing.T) {
	err := New(Config{
		Reset: true,
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

	err = Create(1, "en", 18, 30, 1, 25, 1, "music", "travel", "art")
	if err != nil {
		t.Fatalf("Error during record creation: %v", err)
	}
}

// Benchmarks for the Create function
func BenchmarkCreate(b *testing.B) {
	err := New(Config{
		Reset: true,
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

	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		err := Create(1, "en", 18, 30, 1, 25, 1, "music", "travel", "art")
		if err != nil {
			b.Fatalf("Error during record creation: %v", err)
		}
	}
}
