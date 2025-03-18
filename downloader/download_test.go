package download

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"websiteCopier/metrics"
)

func TestHTTPDownloader_Download_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("this is a test response"))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	metrics := &metrics.Metrics{}
	downloader := &HTTPDownloader{Metrics: metrics}
	ctx := context.Background()
	resultChan := make(chan []byte, 1)
	url := ts.URL[len("http://"):]

	downloader.Download(ctx, url, resultChan)

	select {
	case result := <-resultChan:
		if string(result) != "this is a test response" {
			t.Fatalf("Expected 'this is a test response', got '%s'", string(result))
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for download result")
	}
}

func TestHTTPDownloader_Download_500_Response(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("this is a error response"))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	metrics := &metrics.Metrics{}
	downloader := &HTTPDownloader{Metrics: metrics}
	ctx := context.Background()
	resultChan := make(chan []byte, 1)
	url := ts.URL[len("http://"):]

	downloader.Download(ctx, url, resultChan)

	select {
	case result := <-resultChan:
		if string(result) != "this is a error response" {
			t.Fatalf("Expected 'this is a error response', got '%s'", string(result))
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for download result")
	}
}

func TestHTTPDownloader_Download_Failure(t *testing.T) {
	metrics := &metrics.Metrics{}
	downloader := &HTTPDownloader{Metrics: metrics}
	ctx := context.Background()
	resultChan := make(chan []byte, 1)

	downloader.Download(ctx, "invalid-url", resultChan)
	select {
	case <-resultChan:
		t.Fatal("Expected failure, but got a response")
	case <-time.After(time.Second):
		// Expected behavior: no response should be sent
	}
}

func TestHTTPDownloader_Download_Timeout(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // add delay
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("delayed response"))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	metrics := &metrics.Metrics{}
	downloader := &HTTPDownloader{Metrics: metrics}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	resultChan := make(chan []byte, 1)

	downloader.Download(ctx, ts.URL[len("http://"):], resultChan)
	select {
	case <-resultChan:
		t.Fatal("Expected timeout, but got a response")
	case <-time.After(time.Second):
		// Expected behavior: timeout should prevent response
	}
}
