// Block Size Detection
// Purpose
//     Detect the filesystem’s block size programmatically.
//     Validate chunk size alignment.
//     Compare the chunk size with the block size and provide warnings.

package storage

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

// GetBlockSize Function
//	The function retrieves the block size of the filesystem at the given path using syscall.Statfs.
//	Issues:
//	 The syscall package is platform-specific. Using syscall.Statfs limits your code to Unix-like systems (Linux, macOS, etc.). It won't work on Windows.
//   golang.org/x/sys/unix package instead of syscall for better cross-platform compatibility and more modern APIs

// [deprecated]
func GetBlockSizeVold(path string) (int64, error) {

	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {

		return 0, err
	}

	return stat.Bsize, nil
}

// GetBlockSize returns the block size of the filesystem at the given path.
func GetBlockSize(path string) (int64, error) {

	var fsStat unix.Statfs_t

	if err := unix.Statfs(path, &fsStat); err != nil {
		return 0, fmt.Errorf("failed to get filesystem stats for path %s: %w", path, err)
	}

	return int64(fsStat.Bsize), nil
}

// CheckChunkSize checks if the provided chunk size is aligned with the filesystem's block size.
func CheckChunkSize(path string, chunkSize int64) error {
	if path == "" {
		return fmt.Errorf("path is required")
	}
	if chunkSize <= 0 {
		return fmt.Errorf("chunk size must be a positive value, got %d", chunkSize)
	}

	blockSize, err := GetBlockSize(path)
	if err != nil {
		return fmt.Errorf("failed to check chunk size: %w", err)
	}

	if chunkSize%blockSize != 0 {
		return fmt.Errorf("chunk size (%d) is not aligned with the filesystem block size (%d)", chunkSize, blockSize)
	}
	return nil
}
