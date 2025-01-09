package concurrency

import "testing"

func TestRunWorkers(t *testing.T) {

	counter := 0

	job := func() {
		counter++
	}

	RunWorkers(10, job)

	if counter != 10 {
		t.Fatalf("Expected 10 jobs to run, got %d", counter)
	}

}
