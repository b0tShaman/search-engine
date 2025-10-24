package persist

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"
)

func TestTextFileSaver_SaveToFile_Success(t *testing.T) {
	saver := &TextFileSaver{}
	ctx := context.Background()
	resultChan := make(chan Content, 2)
	var wg sync.WaitGroup
	wg.Add(1)

	resultChan <- Content{URL: "http://test1",Payload: []byte("test data 1")}
	resultChan <- Content{URL: "http://test2",Payload: []byte("test data 2")}
	close(resultChan)

	go saver.SaveToFile(ctx, resultChan, &wg)
	wg.Wait()

	files, err := os.ReadDir(OUTPUT_PATH)
	if err != nil {
		t.Fatal("Failed to read directory")
	}

	var found bool
	for _, file := range files {
		if len(file.Name()) == RAND_BYTES_LENGTH*2+4 && file.Name()[RAND_BYTES_LENGTH*2:] == ".txt" {
			found = true
			os.Remove(OUTPUT_PATH + file.Name()) // Clean up
		}
	}

	if !found {
		t.Fatal("Expected file not found")
	}
}

func TestTextFileSaver_SaveToFile_ContextCancelled(t *testing.T) {
	saver := &TextFileSaver{}
	ctx, cancel := context.WithCancel(context.Background())
	resultChan := make(chan Content, 2)
	var wg sync.WaitGroup
	wg.Add(1)

	go saver.SaveToFile(ctx, resultChan, &wg)

	time.Sleep(100 * time.Millisecond)
	cancel()
	close(resultChan)
	wg.Wait()
}
