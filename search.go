package search

import (
	"slices"
	"sync"
)

type SearchEntry struct {
	ID        uint64
	UserID    uint64
	Language  string
	YourStart int
	YourEnd   int
	YourSex   int
	MyAge     int
	MySex     int
	Priority  bool
	Interests []string
}

type memoryStore struct {
	sync.RWMutex
	entries map[uint64]*SearchEntry
	index   map[string]map[int]map[uint64]*SearchEntry
}

var store *memoryStore

func New() error {
	store = &memoryStore{
		entries: make(map[uint64]*SearchEntry),
		index:   make(map[string]map[int]map[uint64]*SearchEntry),
	}
	return nil
}

func Close() {
	store.Lock()
	defer store.Unlock()
	store.entries = make(map[uint64]*SearchEntry)
	store.index = make(map[string]map[int]map[uint64]*SearchEntry)
}

func Create(
	UserID uint64,
	Language string,
	YourStart int,
	YourEnd int,
	YourSex int,
	MyAge int,
	MySex int,
	Priority bool,
	interests ...string,
) {

	entry := &SearchEntry{
		ID:        UserID,
		UserID:    UserID,
		Language:  Language,
		YourStart: YourStart,
		YourEnd:   YourEnd,
		YourSex:   YourSex,
		MyAge:     MyAge,
		MySex:     MySex,
		Priority:  Priority,
		Interests: interests,
	}

	store.Lock()
	defer store.Unlock()

	if _, exists := store.index[Language]; !exists {
		store.index[Language] = make(map[int]map[uint64]*SearchEntry)
	}

	if _, exists := store.index[Language][MySex]; !exists {
		store.index[Language][MySex] = make(map[uint64]*SearchEntry)
	}

	delete(store.index[Language][MySex], UserID)
	delete(store.entries, UserID)

	store.index[Language][MySex][UserID] = entry
	store.entries[UserID] = entry

}

func Delete(userID uint64) {

	store.Lock()
	defer store.Unlock()

	if user, exists := store.entries[userID]; exists {
		delete(store.index[user.Language][user.MySex], userID)
		delete(store.entries, userID)
	}

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
) *SearchEntry {
	store.RLock()

	var best *SearchEntry
	maxScore := -1
	var bestID uint64

	sexes := []int{YourSex}
	if YourSex == 2 {
		sexes = []int{0, 1}
	}

	iterateUsers(Language, sexes, func(user *SearchEntry) {
		if user.UserID == MyID {
			return
		}
		if YourSex != 2 && YourSex != user.MySex {
			return
		}
		if user.YourSex != 2 && user.YourSex != MySex {
			return
		}

		score := 0

		score += ageScore(MyAge, user.YourStart, user.YourEnd)
		score += ageScore(user.MyAge, YourStart, YourEnd)

		if user.Priority {
			score += 3
		}

		for _, myInterest := range interests {
			if slices.Contains(user.Interests, myInterest) {
				score++
			}
		}

		if score > maxScore {
			best = user
			maxScore = score
			bestID = user.UserID
		}
	})

	store.RUnlock()

	if best != nil {
		Delete(bestID)
	}

	return best
}

func iterateUsers(
	language string,
	sexes []int,
	process func(user *SearchEntry),
) {
	for _, sex := range sexes {
		if entriesByPriority, ok := store.index[language][sex]; ok {
			for _, user := range entriesByPriority {
				process(user)
			}
		}
	}
}

func ageScore(age, min, max int) int {
	if age >= min && age <= max {
		return 5
	}
	var diff int
	if age < min {
		diff = min - age
	} else {
		diff = age - max
	}

	return diff % 5

}
