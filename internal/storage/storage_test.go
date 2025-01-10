package storage

import "testing"

func TestGetBlockSize(t *testing.T) {
	blockSize, err := GetBlockSize(".")
	if err != nil {
		t.Fatalf("GetBlockSize failed: %v", err)
	}
	if blockSize <= 0 {
		t.Fatalf("Block size should be greater than 0")
	}
}
