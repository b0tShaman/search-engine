package main

import (
	"context"
	"os"
	"os/signal"
	download "search_engine/crawler/downloader"
	log "search_engine/crawler/logger"
	"search_engine/crawler/metrics"
	persist "search_engine/crawler/persistence"
	read "search_engine/crawler/reader"
	"sync"
	"syscall"
	"time"
)

const (
	MAX_ROUTINES = 200
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("filePath not provided as argument")
		return
	}

	filePath := os.Args[1]
	urlChan := make(chan string, MAX_ROUTINES) // buffered channel to prevent blocks
	resultChan := make(chan persist.Content, MAX_ROUTINES)
	semaphore := make(chan struct{}, MAX_ROUTINES) // semaphore to limit max routines to 50
	var mainWg sync.WaitGroup

	// Context to be used to cancel/exit from goroutines on Ctrl+C interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := &read.CSVReader{}

	persistMetrics := &metrics.Metrics{}
	persister := &persist.InvertedIndex{Metrics: persistMetrics}

	downloadMetrics := &metrics.Metrics{}
	downloader := &download.HTTPDownloader{Metrics: downloadMetrics}

	// Stage 1 - Read File.
	mainWg.Add(1)
	go reader.StreamURLs(ctx, filePath, urlChan, &mainWg)

	// Stage 2- Download contents with at max 50 goroutines.
	// Spawning an additional separate goroutine for reading from urlChan
	var wg sync.WaitGroup
	mainWg.Add(1)
	go func(mainWg *sync.WaitGroup) {
		defer mainWg.Done()
		for url := range urlChan {
			semaphore <- struct{}{}
			wg.Add(1)
			go func(ctx context.Context, url string) {
				defer wg.Done()
				downloader.Download(ctx, url, resultChan)
				<-semaphore
			}(ctx, url)
		}
		wg.Wait() // wait till all downloads complete
		close(resultChan)
	}(&mainWg)

	mainWg.Add(1)
	// Stage 3 - Persist contents.
	go persister.Persist(ctx, resultChan, &mainWg)

	// Capture OS signals (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan // Wait for Ctrl+C
		log.Warn("Interrupt signal received. Gracefully shutting down after 5 seconds...")

		// Wait for 5 seconds
		<-time.After(5 * time.Second)
		cancel() // Stop every goroutine
	}()

	mainWg.Wait()
	log.Info("Download Metrics:\n", downloadMetrics.String())
	log.Info("Persistence Metrics:\n", persistMetrics.String())
}
