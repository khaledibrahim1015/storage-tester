package metrics

import (
	"fmt"
	"testing"
	"time"
)

func TestMetrics(t *testing.T) {

	// Create a new Metrics instance
	metric := &Metrics{}

	//  start tracking
	metric.Start()

	// simulate i/o operations
	time.Sleep(2 * time.Second)

	metric.RecordRead(1024)  // record 1 KB Read
	metric.RecordWrite(2048) // record 2 KB writen

	// Record CPU and Memory Usage

	if err := metric.RecordCPUUsage(); err != nil {
		t.Fatalf("Error recording CPU usage:%v", err)
	}
	metric.RecordMemoryUsage()

	// Stop Tracking
	metric.Stop()

	// Print metrics
	fmt.Println(metric.String())

}
