package storage

import (
	"testing"
)

func TestAdd(t *testing.T) {
	u := NewURLStorage()

	// Test-case 1: Add new URL
	ok := u.Add("root", "https://example.com")
	if !ok {
		t.Error("Expected Add to return true for a new URL")
	}

	// Test-case 2: Graph check
	if len(u.Graph["root"]) != 1 || u.Graph["root"][0] != "https://example.com" {
		t.Errorf("Graph not updated correctly, got %v", u.Graph["root"])
	}

	// Test-case 3: Second add of visited URL
	ok = u.Add("other", "https://example.com")
	if ok {
		t.Error("Expected Add to return false for an already visited URL")
	}
}
