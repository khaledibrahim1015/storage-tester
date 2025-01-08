package metrics

import (
	"fmt"
	"testing"
	"time"
)

func TestMetrics(t *testing.T) {
	metrics := &Metrics{}

	metrics.Start()

	// simulate i/o operations
	time.Sleep(100 * time.Microsecond)
	metrics.RecordRead(4096)
	metrics.RecordWrite(4096)

	metrics.Stop()

	if metrics.TotalIOPS != 2 {
		t.Fatalf("Expected 2 IOPS, got %d", metrics.TotalIOPS)
	}
	fmt.Printf("TotalIOPS : %v \n", metrics.TotalIOPS)
	if metrics.Throughput() <= 0 {
		t.Fatalf("Throughput should be greater than 0")
	}
	fmt.Printf("Throughput : %v \n", metrics.Throughput())

}
