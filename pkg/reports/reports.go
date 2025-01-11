package reports

import (
	"fmt"
	"os"

	"github.com/khaledibra1015/storage-tester/internal/metrics"
)

func GenerateReport(metrics *metrics.Metrics, format string) (string, error) {
	switch format {
	case "text":
		return metrics.String(), nil

	case "json":
		jsonOutput, err := metrics.Json()
		if err != nil {
			return "", fmt.Errorf("failed to generate JSON report: %v", err)
		}
		return jsonOutput, nil
	case "csv":
		// Create a temporary file to store CSV data
		tempFile, err := os.CreateTemp("", "metrics-*.csv")
		if err != nil {
			return "", fmt.Errorf("failed to create temporary file for CSV: %v", err)
		}
		defer os.Remove(tempFile.Name())
		// Write CSV data to the temporary file
		csvData, err := metrics.CSV(tempFile.Name())
		if err != nil {
			return "", fmt.Errorf("failed to generate CSV report: %v", err)
		}
		return csvData, nil

	default:
		return "", fmt.Errorf("invalid output format: %s", format)
	}

}
