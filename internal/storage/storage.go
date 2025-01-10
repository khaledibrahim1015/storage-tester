// Block Size Detection
// Purpose
//     Detect the filesystem’s block size programmatically.
//     Validate chunk size alignment.
//     Compare the chunk size with the block size and provide warnings.

package storage

import (
	"fmt"
	"syscall"
)

func GetBlockSize(path string) (int64, error) {

	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {

		return 0, err
	}

	return stat.Bsize, nil
}

func CheckChunkSize(path string, chunkSize int64) error {

	if path == "" {
		return fmt.Errorf("path not provided ")
	}
	blockSize, err := GetBlockSize(path)
	if err != nil {
		return err
	}

	if chunkSize%blockSize != 0 {
		return fmt.Errorf("chunk size (%d) is not aligned with block size (%d)", chunkSize, blockSize)
	}
	return nil
}
