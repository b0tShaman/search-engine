package read

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"
)

func createTempCSVFile(t *testing.T, content string) string {
	tempFile, err := os.CreateTemp("", "test.csv")
	if err != nil {
		t.Fatal("Failed to create temp file")
	}
	tempFile.WriteString(content)
	tempFile.Close()
	return tempFile.Name()
}

func TestCSVReader_StreamURLs_Success(t *testing.T) {
	filePath := createTempCSVFile(t, "URL\nhttp://example.com\nhttp://test.com\n")
	defer os.Remove(filePath)

	r := &CSVReader{}
	ctx := context.Background()
	urlChan := make(chan string, 2)
	var wg sync.WaitGroup
	wg.Add(1)

	go r.StreamURLs(ctx, filePath, urlChan, &wg)

	expectedURLs := []string{"http://example.com", "http://test.com"}
	for _, expected := range expectedURLs {
		select {
		case url := <-urlChan:
			if url != expected {
				t.Fatalf("Expected %s, got %s", expected, url)
			}
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for URL")
		}
	}
	wg.Wait()
}

func TestCSVReader_StreamURLs_EmptyFile(t *testing.T) {
	filePath := createTempCSVFile(t, "")
	defer os.Remove(filePath)

	r := &CSVReader{}
	ctx := context.Background()
	urlChan := make(chan string, 2)
	var wg sync.WaitGroup
	wg.Add(1)

	go r.StreamURLs(ctx, filePath, urlChan, &wg)

	select {
	case _, ok := <-urlChan:
		if ok {
			t.Fatal("Expected channel to be closed due to file not found")
		}
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for URL")
	}
	wg.Wait()
}

func TestCSVReader_StreamURLs_FileNotFound(t *testing.T) {
	r := &CSVReader{}
	ctx := context.Background()
	urlChan := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go r.StreamURLs(ctx, "non_existent.csv", urlChan, &wg)

	select {
	case _, ok := <-urlChan:
		if ok {
			t.Fatal("Expected channel to be closed due to file not found")
		}
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for failure handling")
	}
	wg.Wait()
}

func TestCSVReader_StreamURLs_CancelContext(t *testing.T) {
	filePath := createTempCSVFile(t, "URL\nhttp://example.com\nhttp://test.com\n")
	defer os.Remove(filePath)

	r := &CSVReader{}
	ctx, cancel := context.WithCancel(context.Background())
	urlChan := make(chan string, 2)
	var wg sync.WaitGroup
	wg.Add(1)

	go func(){
		time.Sleep(100 * time.Millisecond)
		r.StreamURLs(ctx, filePath, urlChan, &wg)
	}()

	cancel()
	
	_, ok := <-urlChan
	if ok {
		t.Fatal("Expected channel to be closed due to context cancellation")
	}
	wg.Wait()
}
