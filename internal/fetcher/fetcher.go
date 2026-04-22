package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func Fetch(targetURL string) (io.ReadCloser, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
