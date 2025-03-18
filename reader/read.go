package read

import (
	"bufio"
	"context"
	"os"
	"sync"
	log "websiteCopier/logger"
)

type URLReader interface {
	StreamURLs(context.Context, string, chan string, *sync.WaitGroup)
}

type CSVReader struct{}

func (r *CSVReader) StreamURLs(ctx context.Context, filePath string, urlChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		log.Error("Error opening file=", filePath, " with error=", err)
		close(urlChan)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // Read header

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			log.Warn("Stopping CSV reading due to ctx cancelled")
			close(urlChan)
			return
		case urlChan <- scanner.Text():
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error("Error reading CSV file=", filePath, " with error=", err)
	}
	close(urlChan)
	log.Info("Finished reading CSV file=", filePath)
}
