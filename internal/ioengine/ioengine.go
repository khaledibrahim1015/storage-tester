package ioengine

import (
	"io"
	"os"
)

// WriteFile writes data to a file at the specified path.
func WriteFile(filename string, data []byte) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// WriteFile writes data to a file in chunks.
func WriteFileWithChunk(filename string, chunkSize int, data []byte) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for offset := 0; offset < len(data); offset += chunkSize {
		endOffset := offset + chunkSize

		if endOffset > len(data) {
			// to get last offset in data
			endOffset = len(data)
		}
		_, err = file.Write(data[offset:endOffset])

		if err != nil {
			return err
		}

	}

	return nil

}

// ReadFile reads data from a file at the specified path.
func ReadFile(filename string, data []byte) (int, error) {

	// data is allocated buffer to read data in
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Read(data)

}

//	another util function to get data from
//
// ReadFileData
func ReadFileData(filename string) ([]byte, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// FileInfo
	// Get the file size to allocate the appropriate buffer
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	filesize := fileinfo.Size()

	// allocate buffer to store file data
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	return buffer, err
}

// ReadFile reads data from a file in chunks.
// func ReadFileWithChunks(filename string, chunkSize int, data []byte) (int, error) {

// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer file.Close()

// 	//  Get FileStat FileInfo
// 	fileinfo, err := file.Stat()
// 	if err != nil {
// 		return 0, err
// 	}
// 	var n int
// 	for offset := 0; offset < int(fileinfo.Size()); offset += chunkSize {

// 		endOffset := offset + chunkSize
// 		n, err = file.ReadAt(data[offset:endOffset], int64(offset))
// 		if err != nil {
// 			return n, err
// 		}
// 	}

// 	return n, nil
// }

func ReadFileWithChunks(filename string, chunkSize int, data []byte) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var totalRead int
	chunkBuffer := make([]byte, chunkSize) // Buffer for reading chunks

	for {
		// Read a chunk into the temporary buffer
		n, err := file.Read(chunkBuffer)
		if err != nil && err != io.EOF {
			return totalRead, err // Return any error other than EOF
		}

		// Copy the chunk into the data buffer
		if totalRead+n > len(data) {
			// Avoid buffer overflow by only copying what fits
			copy(data[totalRead:], chunkBuffer[:len(data)-totalRead])
			totalRead = len(data)
			break
		}
		copy(data[totalRead:], chunkBuffer[:n])
		totalRead += n

		// Stop if we've reached the end of the file
		if err == io.EOF {
			break
		}
	}

	return totalRead, nil
}
