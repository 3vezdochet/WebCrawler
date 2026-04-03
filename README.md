# WebCrawler

A highly concurrent, scalable web crawler written in Go. 

This project is designed to efficiently traverse web pages, extract links, and manage network requests using Go's powerful concurrency model. The architecture is built with a focus on memory safety, avoiding race conditions, and graceful degradation.

## 🏗 Architecture

The crawler uses a **Worker Pool** pattern to manage concurrent HTTP requests without exhausting system resources. 

Key components:
*   **Dispatcher:** Manages the seed URLs and distributes tasks.
*   **Worker Pool:** A pool of goroutines executing HTTP fetches safely.
*   **Fetcher:** Handles network requests with strict timeouts to prevent goroutine leaks.
*   **Parser:** Extracts `<a>` tags and raw URLs from HTML streams.
*   **Storage (Seen URLs):** A thread-safe, mutex-guarded storage to track visited links and prevent infinite crawling loops.

## 📁 Project Structure

The repository follows the Standard Go Project Layout:

```text
webcrawler/
├── cmd/
│   └── crawler/
│       └── main.go       # Application entry point
├── internal/
│   ├── fetcher/          # HTTP client logic and timeout management
│   ├── parser/           # HTML parsing and link extraction
│   └── storage/          # Thread-safe in-memory cache for visited URLs
├── go.mod                # Go module dependencies
└── README.md
```

## 🚀 Getting Started

### Prerequisites
* Go 1.22 or higher

### Running the application
*(Instructions will be added as the CLI is implemented)*
```bash
go run cmd/crawler/main.go
```

## 🗺 Roadmap

- [ ] Implement a robust HTTP client with timeouts (`internal/fetcher`).
- [ ] Build a thread-safe URL storage using `sync.RWMutex` (`internal/storage`).
- [ ] Implement HTML parsing to extract links (`internal/parser`).
- [ ] Tie components together using a Worker Pool and `sync.WaitGroup` in `main`.
- [ ] Add graceful shutdown (context cancellation).
- [ ] Implement Rate Limiting to prevent server overload.

## 🛠 Tech Stack
*   **Language:** Go
*   **Libraries:** Standard library (`net/http`, `sync`, `context`)