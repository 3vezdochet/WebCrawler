package storage

import "sync"

type URLStorage struct {
	mu   sync.Mutex
	data map[string]struct{}
}

func NewURLStorage() *URLStorage {
	return &URLStorage{
		data: make(map[string]struct{}),
	}
}

func (u *URLStorage) Add(url string) bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, exists := u.data[url]; exists {
		return false
	}

	u.data[url] = struct{}{}
	return true
}
