package metrics

import (
	log "search_engine/crawler/logger"
	"sync/atomic"
	"time"
)

type Metrics struct {
	Name         string
	successCount int64
	failureCount int64
	totalTime    int64
}

func (m *Metrics) IncrementSuccess() {
	atomic.AddInt64(&m.successCount, 1)
}

func (m *Metrics) IncrementFailure() {
	atomic.AddInt64(&m.failureCount, 1)
}

func (m *Metrics) AddTotalTime(duration time.Duration) {
	atomic.AddInt64(&m.totalTime, int64(duration))
}

func (m *Metrics) LogMetrics() {
	totalProcessed := atomic.LoadInt64(&m.successCount) + atomic.LoadInt64(&m.failureCount)
	log.Infof("[%s] Total processed: %d", m.Name, totalProcessed)
	log.Infof("[%s] Success count: %d", m.Name, atomic.LoadInt64(&m.successCount))
	log.Infof("[%s] Failure count: %d", m.Name, atomic.LoadInt64(&m.failureCount))

	successCount := atomic.LoadInt64(&m.successCount)
	timeTaken := atomic.LoadInt64(&m.totalTime)
	avgDuration := time.Duration(0)
	if successCount != 0 {
		avgDuration = time.Duration(timeTaken / successCount)
	}

	log.Infof("[%s] Average duration: %v secs", m.Name, avgDuration.Seconds())
}
