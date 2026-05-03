package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"WebCrawler/internal/fetcher"
	"WebCrawler/internal/parser"
	"WebCrawler/internal/storage"
)

func main() {
	seedURL := "https://habr.com/ru/companies/yadro/articles/895084/"
	numWorkers := 10
	rateLimit := time.Tick(100 * time.Millisecond)

	// Context for Graceful Shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var pagesParsed atomic.Int32
	maxPages := int32(100)

	store := storage.NewURLStorage()
	queue := make(chan string, 1000)
	var wg sync.WaitGroup

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\n[EMERGENCY] Recovered from panic: %v\n", r)
		}
		fmt.Println("[System] Executing final data save...")
		if err := store.ExportToJSON("tree_emergency.json"); err != nil {
			fmt.Printf("Error during emergency save %v\n", err)
		}
	}()

	for i := 1; i <= numWorkers; i++ {
		go worker(ctx, cancel, i, queue, store, &wg, rateLimit, &pagesParsed, maxPages)
	}

	store.Add("root", seedURL)
	wg.Add(1)
	queue <- seedURL

	wg.Wait()
	close(queue)

	fmt.Println("Crawling completed")
	fmt.Println("Saving URL tree to tree.json...")
	if err := store.ExportToJSON("tree.json"); err != nil {
		fmt.Printf("Error saving JSON: %v\n", err)
	} else {
		fmt.Println("Tree saved successfully!")
	}
}

func worker(
	ctx context.Context,
	cancel context.CancelFunc,
	id int,
	queue chan string,
	store *storage.URLStorage,
	wg *sync.WaitGroup,
	rateLimit <-chan time.Time,
	pagesParsed *atomic.Int32,
	maxPages int32,
) {
	for targetURL := range queue {
		// Queue pop: if pulled cancel signal, pop WaitGroup and grab new URL
		// without work. Save from Deadlock
		if ctx.Err() != nil {
			wg.Done()
			continue
		}

		<-rateLimit

		fmt.Printf("[Worker %d] Fetching: %s\n", id, targetURL)
		body, err := fetcher.Fetch(targetURL)
		if err != nil {
			fmt.Printf("[Worker %d] Error fetching %s: %v\n", id, targetURL, err)
			wg.Done()
			continue
		}

		currentCount := pagesParsed.Add(1)
		if currentCount >= maxPages {
			fmt.Println("[Limit reached] Stopping crawler...")
			cancel() // Call cancel context for all goroutines
		}

		links := parser.ExtractLinks(body, targetURL)
		if err := body.Close(); err != nil {
			fmt.Printf("Body closing error: %v\n", err)
		}

		for _, link := range links {
			if ok := store.Add(targetURL, link); ok {
				wg.Add(1)
				go func(l string) {
					select {
					case queue <- l:
					case <-ctx.Done():
						wg.Done()
					}
				}(link)
			}
		}
		wg.Done()
	}
}
