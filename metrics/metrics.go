package metrics

import (
	"sync/atomic"
	"time"
	"websiteCopier/logger"
)

type Metrics struct {
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
	log.Infof("Total URLs processed: %d", totalProcessed)
	log.Infof("Success count: %d", atomic.LoadInt64(&m.successCount))
	log.Infof("Failure count: %d", atomic.LoadInt64(&m.failureCount))

	time_taken := m.totalTime
	if m.successCount != 0{
		time_taken = m.totalTime/m.successCount
	}
	
	log.Infof("Average Download duration: %v secs",time.Duration(time_taken).Seconds() )
}
