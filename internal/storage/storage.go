package storage

import (
	"encoding/json"
	"fmt"
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

func (u *URLStorage) PrintTree(root string, indent string, isLast bool) {
	u.mu.Lock()
	nodes := u.Graph[root]
	u.mu.Unlock()

	marker := "├── "
	if isLast {
		marker = "└── "
	}

	// Print current node (except root)
	if root != "root" {
		fmt.Print(indent)
		fmt.Print(marker)
		fmt.Print(root)
	}

	// Calculate new marker for children
	newIndent := indent
	if root != "root" {
		if isLast {
			newIndent += "    "
		} else {
			newIndent += "│   "
		}
	}

	// Recursive traversal
	for i, child := range nodes {
		lastChild := i == len(nodes)-1
		u.PrintTree(child, newIndent, lastChild)
	}
}
