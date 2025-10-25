# Search-Engine in Golang

## Overview

**Search-Engine** is a simple Go application to crawl websites from a CSV file, build an inverted index of words and search for URLs containing a particular word.

## Installation

Ensure you have **Go 1.23** installed on your system.

Clone the repository:

```sh
git clone https://github.com/b0tShaman/search-engine.git
cd search-engine
```

Install dependencies:

```sh
go mod tidy
```

## Dependencies

This project uses the following external dependencies:

* [logrus](https://github.com/sirupsen/logrus): For structured logging.

## Usage

### 1. Crawl Websites

To crawl websites listed in a CSV file and build an inverted index:

```sh
go run crawler/crawler.go urls.csv
```

* `urls.csv` should contain the list of URLs to crawl.

### 2. Search for a Word

To search for URLs containing a particular word:

```sh
go run search/search.go
```

* You will be prompted to **Enter search word**.
* The program will display the list of URLs containing that word from the inverted index.

## Running Unit Tests

To run unit tests:

```sh
go test ./...
```

This will execute all tests in the project.

## Assumptions and Design Decisions

* The inverted index is stored as a map[string]map[string]int and persisted using Go's gob encoder for fast serialization and deserialization.
* The program supports concurrent crawling with a cap on goroutines for efficiency.
* Graceful shutdown is implemented to allow the program to exit cleanly when interrupted.
