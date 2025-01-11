package metrics

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// This package will collect and calculate performance metrics.

type Metrics struct {
	StartTime    time.Time
	EndTime      time.Time
	BytesWritten int64
	BytesRead    int64
	TotalIOPS    int64
	CPUUsage     float64 // CPU usage in percentage
	MemoryUsage  uint64  // Memory usage in bytes
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

// Memory Usage Tracking:
// RecordMemoryUsage method already tracks memory usage using Go's runtime package.
// It records the allocated memory in bytes (memStats.Alloc).

func (m *Metrics) RecordMemoryUsage() {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.MemoryUsage = memStats.Alloc // Allocated memory in bytes
}

// CPU Usage Tracking:
// RecordCPUUsage method to track CPU usage using the gopsutil/cpu package.
// The cpu.Percent(0, false) function returns the CPU usage as a percentage for all cores.
// Setting the second argument to false gives the average CPU usage across all cores.

func (m *Metrics) RecordCPUUsage() error {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}
	if len(percentages) > 0 {
		m.CPUUsage = percentages[0] // Average CPU usage across all cores
	}

	return nil

}

func (m *Metrics) String() string {

	return fmt.Sprintf(
		"Latency: %v\n"+
			"Throughput: %.2f bytes/sec\n"+
			"Bytes Written: %d\n"+
			"Bytes Read: %d\n"+
			"Total IOPS: %d\n"+
			"CPU Usage: %.2f%%\n"+
			"Memory Usage: %d bytes\n",
		m.Latency(),
		m.Throughput(),
		m.BytesWritten,
		m.BytesRead,
		m.TotalIOPS,
		m.CPUUsage,
		m.MemoryUsage,
	)

}

func (m *Metrics) Json() (string, error) {

	jsonbytes, err := json.MarshalIndent(m, "", "  ") // 2 spaces
	if err != nil {
		return "", fmt.Errorf("failed to serialize metrics to JSON: %v", err)
	}
	jsonstring := string(jsonbytes)
	return string(jsonstring), nil
}

func (m *Metrics) CSV(filename string) (string, error) {

	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	//  intaiate csv writer
	csvWriter := csv.NewWriter(file)
	// TODO : USING REFLECTION
	// Write the header
	header := []string{
		"Latency",
		"Throughput (bytes/sec)",
		"Bytes Written",
		"Bytes Read",
		"Total IOPS",
		"CPU Usage (%)",
		"Memory Usage (bytes)",
	}
	if err := csvWriter.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write the data

	// Write the data
	record := []string{
		m.Latency().String(),
		fmt.Sprintf("%.2f", m.Throughput()),
		fmt.Sprintf("%d", m.BytesWritten),
		fmt.Sprintf("%d", m.BytesRead),
		fmt.Sprintf("%d", m.TotalIOPS),
		fmt.Sprintf("%.2f", m.CPUUsage),
		fmt.Sprintf("%d", m.MemoryUsage),
	}
	if err := csvWriter.Write(record); err != nil {
		return "", fmt.Errorf("failed to write CSV record: %v", err)
	}

	// Flush the CSV writer
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		fmt.Errorf("failed to flush CSV writer: %v", err)
	}
	// Read data and return
	csvDataBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read CSV data: %v", err)
	}

	return string(csvDataBytes), nil

}
