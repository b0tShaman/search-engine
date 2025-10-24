# Search-Engine in Golang

## Overview

**Search-Engine** is a simple Go application to crawl websites from a CSV file, build an inverted index of words and search for URLs containing a particular word.

The crawler generates inverted index files where each file is named after a word, and the content of the file contains URLs of websites that include that word. The results are saved in the `output` directory.

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
* Each word found in the crawled pages will have a corresponding file in the `output` directory.
* The `output` directory will be created automatically if it does not exist.

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

* The inverted index is built with a file per word, containing all URLs where the word appears.
* The program supports concurrent crawling with a cap on goroutines for efficiency.
* Graceful shutdown is implemented to allow the program to exit cleanly when interrupted.
