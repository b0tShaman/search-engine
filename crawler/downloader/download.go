package download

import (
	"context"
	"io"
	"net/http"
	"time"
	log "search_engine/crawler/logger"
	"search_engine/crawler/metrics"
	persist "search_engine/crawler/persistence"
)

const (
	HTTP_PREFIX = "http://"
)

type Downloader interface {
	Download(context.Context, string, chan []byte)
}

type HTTPDownloader struct {
	Metrics *metrics.Metrics
}

func (h *HTTPDownloader) Download(ctx context.Context, url string, resultChan chan persist.Content) {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, HTTP_PREFIX+url, nil)
	if err != nil {
		h.Metrics.IncrementFailure()
		log.Error("http NewRequestWithContext failed for url=", url, " with error=", err)
		return
	}
	req.Header.Set("User-Agent", "search_engine/crawler/1.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.Metrics.IncrementFailure()
		log.Error("http Send request failed for url=", url, " with error=", err, " RTT=", time.Since(start))
		return
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		h.Metrics.IncrementFailure()
		log.Error("http Reading response failed for url=", url, " with error=", err, " RTT=", time.Since(start))
		return
	}

	h.Metrics.IncrementSuccess()
	h.Metrics.AddTotalTime(time.Since(start))

	resultChan <- persist.Content{
		URL:     url,
		Payload: content,
	}
}
