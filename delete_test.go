package search

import (
	"testing"
)

// Test case for deleting a user record by ID
func TestDeleteUserByID(t *testing.T) {
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

	// Create a record to be deleted
	err = Create(1, "en", 18, 30, 1, 25, 1, "music", "travel", "art")
	if err != nil {
		t.Fatalf("Error during record creation: %v", err)
	}

	// Delete the record by user ID
	err = Delete(1)
	if err != nil {
		t.Fatalf("Error during record deletion: %v", err)
	}

	// Verify the record no longer exists
	result, err := Search("en", 18, 30, 1, 1, 25)
	if err != nil {
		t.Fatalf("Error during search: %v", err)
	}
	if result != nil {
		t.Errorf("Expected no results, but found: %+v", result)
	}
}

// Benchmark for the DeleteUserByID function
func BenchmarkDeleteUserByID(b *testing.B) {
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

	for i := 0; i < 1000; i++ { // Insert 50,000 records
		err = Create(uint64(i), "en", 18, 30, 1, 25, 1, "music", "travel", "art")
		if err != nil {
			b.Fatalf("Error during record creation: %v", err)
		}
	}

	b.ResetTimer() // Reset the timer to exclude setup time

	// Benchmark the delete function
	for i := 0; i < b.N; i++ {
		err := Delete(uint64(i % 1000)) // Use record IDs in range of inserted records
		if err != nil {
			b.Fatalf("Error during delete: %v", err)
		}
	}
}
