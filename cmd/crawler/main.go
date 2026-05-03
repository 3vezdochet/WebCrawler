package main

import (
	"fmt"
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

	var pagesParsed atomic.Int32
	maxPages := int32(100)

	store := storage.NewURLStorage()
	queue := make(chan string, 1000)
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		go worker(i, queue, store, &wg, rateLimit, &pagesParsed, maxPages)
	}

	store.Add("root", seedURL)
	wg.Add(1)
	queue <- seedURL

	wg.Wait()
	close(queue)

	fmt.Println("Crawling completed")

	fmt.Println("Structure of found URLs:")
	store.PrintTree("root", "", true)
}

func worker(
	id int,
	queue chan string,
	store *storage.URLStorage,
	wg *sync.WaitGroup,
	rateLimit <-chan time.Time,
	pagesParsed *atomic.Int32,
	maxPages int32,
) {
	for targetURL := range queue {
		fmt.Printf("[Worker %d] Fetching: %s\n", id, targetURL)
		<-rateLimit
		body, err := fetcher.Fetch(targetURL)
		if err != nil {
			fmt.Printf("[Worker %d] Error fetching %s: %v\n", id, targetURL, err)
			wg.Done()
			continue
		}

		currentCount := pagesParsed.Add(1)
		if currentCount >= maxPages {
			fmt.Println("[Limit reached] Stopping crawler...")
		}

		links := parser.ExtractLinks(body, targetURL)
		if err := body.Close(); err != nil {
			fmt.Printf("Body closing error: %s\n", err)
		}

		for _, link := range links {
			if ok := store.Add(targetURL, link); ok {
				wg.Add(1)
				go func(l string) {
					queue <- l
				}(link)
			}

		}

		wg.Done()
	}
}
