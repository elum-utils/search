package search

import (
    "testing"
)

func TestSearchIgnoreSex(t *testing.T) {
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

    // Create potential matches
    Create(1, "en", 18, 30, 1, 25, 0, "music", "art") // User 1: Female
    Create(2, "en", 18, 30, 1, 25, 1, "music", "art") // User 2: Male

    // Search with any gender
    result, err := Search(3, "en", 18, 30, 2, 1, 25, "music")
    if err != nil {
        t.Fatalf("Error during search: %v", err)
    }

    // Since we allow any gender, we should get the result with the highest priority (which is arbitrary here)
    if result == nil || (result.UserID != 1 && result.UserID != 2) {
        t.Errorf("Expected to find a user of any gender, got %+v", result)
    }
}

func TestSearchIgnoreSexNoInterests(t *testing.T) {
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

    // Create potential matches
    Create(1, "en", 18, 30, 1, 25, 0) // User 1: Female
    Create(2, "en", 18, 30, 1, 25, 1) // User 2: Male

    // Search with any gender and no specific interests
    result, err := Search(3, "en", 18, 30, 2, 1, 25)
    if err != nil {
        t.Fatalf("Error during search: %v", err)
    }

    // Since we allow any gender, we should get the result with the highest priority (which is arbitrary here)
    if result == nil || (result.UserID != 1 && result.UserID != 2) {
        t.Errorf("Expected to find a user of any gender, got %+v", result)
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