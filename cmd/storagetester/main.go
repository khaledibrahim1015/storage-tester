package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/khaledibra1015/storage-tester/internal/concurrency"
	"github.com/khaledibra1015/storage-tester/internal/ioengine"
	"github.com/khaledibra1015/storage-tester/internal/metrics"
	"github.com/khaledibra1015/storage-tester/pkg/reports"
)

const (
	chunk_Size    = "chunk-size"
	file_name     = "file-name"
	file_size     = "file-size"
	operation     = "operation"
	ConcurrentOps = "concurrency"
	output_format = "output-format"
)

// CLICMD defines the structure for command-line flags
type CLICMD struct {
	ChunkSize     int64
	FileName      string
	FileSize      int64
	Operation     string
	ConcurrentOps int
	OutputFormat  string
}

// Execute_CMD parses and validates command-line flags
func Execute_CMD() CLICMD {
	// Define default values
	cliflags := CLICMD{
		ChunkSize:     4096,
		FileName:      "test.txt",
		FileSize:      1024 * 1024 * 1024, // 1 GB
		Operation:     "write",
		ConcurrentOps: 1,
		OutputFormat:  "text",
	}

	// Define flags
	chunkSize := flag.Int64(chunk_Size, cliflags.ChunkSize, "Chunk size for I/O operations (in bytes)")
	fileName := flag.String(file_name, cliflags.FileName, "Name of the test file")
	fileSize := flag.Int64(file_size, cliflags.FileSize, "Size of the test file (in bytes)")
	operation := flag.String(operation, cliflags.Operation, "Type of I/O operation (read, write, mixed)")
	concurrency := flag.Int(ConcurrentOps, cliflags.ConcurrentOps, "Number of workers for parallel I/O operations")
	outputFormat := flag.String(output_format, cliflags.OutputFormat, "Format for performance reports (text, json, csv)")

	// Parse flags
	flag.Parse()

	// Validate flags
	if *chunkSize <= 0 || *fileSize <= 0 || *concurrency <= 0 {
		log.Fatal("Invalid flag values: chunk-size, file-size, and concurrency must be greater than 0")
	}

	// Return parsed flags as a CLICMD struct
	return CLICMD{
		ChunkSize:     *chunkSize,
		FileName:      *fileName,
		FileSize:      *fileSize,
		Operation:     *operation,
		ConcurrentOps: *concurrency,
		OutputFormat:  *outputFormat,
	}
}

func main() {
	// Execute command-line parsing and validation
	cmd := Execute_CMD()

	// Print parsed values
	fmt.Println("Configuration : ")
	fmt.Printf("Chunk Size: %d bytes\n", cmd.ChunkSize)
	fmt.Printf("File Name: %s\n", cmd.FileName)
	fmt.Printf("File Size: %d bytes\n", cmd.FileSize)
	fmt.Printf("Operation: %s\n", cmd.Operation)
	fmt.Printf("Concurrency: %d\n", cmd.ConcurrentOps)
	fmt.Printf("Output Format: %s\n", cmd.OutputFormat)

	//  Run The Test

}

func RunTest(filename string, filesize int64, chunksize int64, opertation string, concurency int, outputformat string) error {
	// Generate the test file if it doesn't exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := ioengine.GenerateTestFile(filename, filesize); err != nil {
			return fmt.Errorf("failed to generate test file: %v", err)
		}

	}

	// Initialize metrics
	metrics := &metrics.Metrics{}
	metrics.Start()

	// operation
	switch opertation {
	case "write":
		// strat run writetest ops
		if err := RunWriteTest(filename, int(chunksize), filesize, concurency, metrics); err != nil {
			return fmt.Errorf("write test failed: %v", err)
		}

	case "read":
		// RunReadTest(filename , filesize , chunksize , concurency, metrics)
		if err := RunReadTest(filename, chunksize, concurency, metrics); err != nil {
			return err
		}

	default:
		return fmt.Errorf("invalid operation: %s", operation)
	}

	// Record CPU and Memory Usage

	if err := metrics.RecordCPUUsage(); err != nil {
		log.Fatalf("Error recording CPU usage:%v", err)
	}
	metrics.RecordMemoryUsage()
	// Stop Tracking Metrics
	metrics.Stop()

	// Generate the report
	report, err := reports.GenerateReport(metrics, outputformat)
	if err != nil {
		return fmt.Errorf("failed to generate report: %v", err)
	}
	// Print the report
	fmt.Println("\nPerformance Report:")
	fmt.Println(report)
	return nil

}

func RunReadTest(filename string, chunksize int64, concurency int, metrics *metrics.Metrics) error {

	job := func() {
		data, err := ioengine.ReadFileWithChunks(filename, int(chunksize))
		if err != nil {
			log.Printf("Read failed: %v", err)
		}
		metrics.RecordRead(int64(len(data)))

	}
	concurrency.RunWorkers(concurency, job)
	return nil
}

func RunWriteTest(filename string, chunksize int, fileSize int64, concurency int, metrics *metrics.Metrics) error {

	// Generate random data for the test file
	data := make([]byte, fileSize)
	// Fill the slice with random data
	// Tmetricshe slice should now contain random bytes instead of only zeroes.
	//  here we can replace it with GenerateTestFile
	_, err := rand.Read(data)
	if err != nil {
		log.Fatalf("Failed to generate random data: %v", err)
	}

	job := func() {
		if err := ioengine.WriteFileWithChunk(filename, chunksize, data); err != nil {
			log.Printf("Write failed: %v", err)
		}
		metrics.RecordWrite(int64(len(data)))
	}
	concurrency.RunWorkers(concurency, job)

	// ctx := context.Background()
	// job2 := func(ctx context.Context, workerID int) error {
	// 	if err := ioengine.WriteFileWithChunk(filename, chunksize, data); err != nil {
	// 		log.Printf("Write failed: %v", err)
	// 	}
	// 	metrics.RecordWrite(int64(len(data)))
	// 	return nil
	// }
	// if err :=  concurrency.AdvancedRunWorkers(ctx, concurency, job2) ; err != nil {
	// 	return err
	// }
	return nil

}
