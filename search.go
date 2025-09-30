package search

import (
	"slices"
	"sync"
	"time"
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
	history map[uint64]map[uint64]time.Time
	timeout time.Duration
	delay   time.Duration
}

var (
	cleanupOnce sync.Once
	store       *memoryStore = &memoryStore{
		entries: make(map[uint64]*SearchEntry),
		index:   make(map[string]map[int]map[uint64]*SearchEntry),
		history: make(map[uint64]map[uint64]time.Time),
		timeout: 0,
		delay:   time.Minute,
	}
)

func init() {
	cleanupOnce.Do(func() {
		go func() {
			for range time.Tick(store.delay) {
				cleanupExpiredHistory()
			}
		}()
	})
}

func cleanupExpiredHistory() {
	store.Lock()
	defer store.Unlock()

	now := time.Now()
	for userID, historyMap := range store.history {
		for otherUserID, expireTime := range historyMap {
			if now.After(expireTime) {
				delete(historyMap, otherUserID)
			}
		}
		if len(historyMap) == 0 {
			delete(store.history, userID)
		}
	}
}

func SetTimeout(d time.Duration) {
	store.Lock()
	defer store.Unlock()
	store.timeout = d
}

func SetDelay(d time.Duration) {
	store.Lock()
	defer store.Unlock()
	store.delay = d
}

func Close() {
	store.Lock()
	defer store.Unlock()

	store.entries = make(map[uint64]*SearchEntry)
	store.index = make(map[string]map[int]map[uint64]*SearchEntry)
	store.history = make(map[uint64]map[uint64]time.Time)
	store.timeout = 5 * time.Minute
	store.delay = time.Minute

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

		if isBlocked(MyID, user.UserID) {
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

		if score > 0 && score > maxScore {
			best = user
			maxScore = score
			bestID = user.UserID
		}
	})

	store.RUnlock()

	if best != nil {
		addToHistory(MyID, bestID)
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

func ageScore(age, mi, ma int) int {
	if age >= mi && age <= ma {
		return 5
	}

	var diff int
	if age < mi {
		diff = mi - age
	} else {
		diff = age - ma
	}

	if diff > 15 {
		return 0
	}

	score := 5 - (diff / 3)
	if score < 0 {
		return 0
	}

	return score
}

func isBlocked(a, b uint64) bool {
	store.RLock()
	defer store.RUnlock()

	if store.timeout == 0 {
		return false
	}

	if m, ok := store.history[a]; ok {
		if t, ok2 := m[b]; ok2 {
			if time.Now().Before(t) {
				return true
			}
		}
	}
	return false
}

func addToHistory(a, b uint64) {
	store.Lock()
	defer store.Unlock()

	if store.timeout == 0 {
		return
	}

	expire := time.Now().Add(store.timeout)

	if _, ok := store.history[a]; !ok {
		store.history[a] = make(map[uint64]time.Time)
	}
	if _, ok := store.history[b]; !ok {
		store.history[b] = make(map[uint64]time.Time)
	}

	store.history[a][b] = expire
	store.history[b][a] = expire
}
