package storage

import (
	"encoding/json"
	"os"
	"sync"
)

type URLStorage struct {
	mu sync.Mutex
	// Field of visited URLs
	data map[string]struct{}
	// Add graph (tree) of linkings. Key - parent URL, value - list of finded URLs
	Graph map[string][]string
}

func NewURLStorage() *URLStorage {
	return &URLStorage{
		data:  make(map[string]struct{}),
		Graph: make(map[string][]string),
	}
}

func (u *URLStorage) Add(source, url string) bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.Graph[source] = append(u.Graph[source], url)

	if _, exists := u.data[url]; exists {
		return false
	}

	u.data[url] = struct{}{}
	return true
}

func (u *URLStorage) ExportToJSON(filename string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	data, err := json.MarshalIndent(u.Graph, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
