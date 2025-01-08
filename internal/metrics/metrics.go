package metrics

import "time"

// This package will collect and calculate performance metrics.

type Metrics struct {
	StartTime    time.Time
	EndTime      time.Time
	BytesWritten int64
	BytesRead    int64
	TotalIOPS    int64
}

func (m *Metrics) Start() {
	m.StartTime = time.Now()
}

func (m *Metrics) Stop() {
	m.EndTime = time.Now()
}

func (m *Metrics) RecordWrite(bytes int64) {
	m.BytesWritten += bytes
	m.TotalIOPS++
}

func (m *Metrics) RecordRead(bytes int64) {
	m.BytesRead += bytes
	m.TotalIOPS++
}

func (m *Metrics) Latency() time.Duration {
	return m.EndTime.Sub(m.StartTime)
}

func (m *Metrics) Throughput() float64 {
	durationInSeconds := m.Latency().Seconds()
	if durationInSeconds == 0 {
		return 0
	}
	return float64(m.BytesRead+m.BytesWritten) / durationInSeconds
}
