package main

import (
	"fmt"
	"sync"
	"time"

	"WebCrawler/internal/fetcher"
	"WebCrawler/internal/parser"
	"WebCrawler/internal/storage"
)

func main() {
	seedURL := "https://golang.org/"
	numWorkers := 10
	rateLimit := time.Tick(200 * time.Millisecond)

	store := storage.NewURLStorage()
	queue := make(chan string, 1000)
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		go worker(i, queue, store, &wg, rateLimit)
	}

	store.Add(seedURL)
	wg.Add(1)
	queue <- seedURL

	wg.Wait()
	close(queue)

	fmt.Println("Crawling completed")
}

func worker(id int, queue chan string, store *storage.URLStorage, wg *sync.WaitGroup, rateLimit <-chan time.Time) {
	for targetURL := range queue {
		fmt.Printf("[Worker %d] Fetching: %s\n", id, targetURL)
		<-rateLimit
		body, err := fetcher.Fetch(targetURL)
		if err != nil {
			fmt.Printf("[Worker %d] Error fetching %s: %v\n", id, targetURL, err)
			wg.Done()
			continue
		}

		links := parser.ExtractLinks(body, targetURL)
		body.Close()

		for _, link := range links {
			if ok := store.Add(link); ok {
				wg.Add(1)
				go func(l string) {
					queue <- l
				}(link)
			}

		}

		wg.Done()
	}
}
