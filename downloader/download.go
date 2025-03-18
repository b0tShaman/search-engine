package download

import (
	"context"
	"io"
	"net/http"
	"time"
	"websiteCopier/logger"
	"websiteCopier/metrics"
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

func (h *HTTPDownloader) Download(ctx context.Context, url string, resultChan chan []byte) {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, HTTP_PREFIX+url, nil)
	if err != nil {
		h.Metrics.IncrementFailure()
		log.Error("http NewRequestWithContext failed for url=", url, " with error=", err)
		return
	}

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

	resultChan <- content
}
