// Implement Concurrency Management
// Purpose
//     Enable parallel I/O operations to simulate real-world workloads.
//     Use Goroutines and channels for efficient concurrency.

package concurrency

import (
	"context"
	"fmt"
	"sync"
)

// Type
type Func func(ctx context.Context, workerID int) error

func RunWorkers(numWorkers int, job func()) {
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			job()
		}()
	}
	wg.Wait()
}

// RunWorkers runs the given job concurrently using the specified number of workers.
// It supports cancellation via context and returns an error if any worker fails.

func AdvancedRunWorkers(ctx context.Context, numWorkers int, job Func) error {

	var wg sync.WaitGroup

	errChan := make(chan error, numWorkers) // Buffered channel to collect errors

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			if err := job(ctx, workerID); err != nil {
				// add errors in channels
				errChan <- fmt.Errorf("worker %d failed: %w", workerID, err)
			}

		}(i)
	}
	// Wait for all workers to finish
	wg.Wait()

	// Check for errors
	close(errChan)
	for err := range errChan {
		return err // Return the first error encountered
	}
	return nil
}
