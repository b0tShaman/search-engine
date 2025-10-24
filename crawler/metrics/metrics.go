package metrics

import (
	"fmt"
	"sync/atomic"
	"time"
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

func (m *Metrics) String() string {
	totalProcessed := atomic.LoadInt64(&m.successCount) + atomic.LoadInt64(&m.failureCount)
	successCount := atomic.LoadInt64(&m.successCount)
	failureCount := atomic.LoadInt64(&m.failureCount)
	timeTaken := atomic.LoadInt64(&m.totalTime)
	avgDuration := time.Duration(0)
	if successCount != 0 {
		avgDuration = time.Duration(timeTaken / successCount)
	}

	return "Total processed: " + fmt.Sprintf("%d", totalProcessed) + "\n" +
		"Success count: " + fmt.Sprintf("%d", successCount) + "\n" +
		"Failure count: " + fmt.Sprintf("%d", failureCount) + "\n" +
		"Average duration: " + fmt.Sprintf("%.2f secs", avgDuration.Seconds())
}
