package ioengine

import (
	"fmt"
	"os"
	"testing"
)

func TestWriteAndReadFile(tst *testing.T) {

	filename := "test.txt"
	data := []byte("kosom elsisi")

	// Test  WriteFile
	if err := WriteFile(filename, data); err != nil {
		tst.Fatalf("writefile failed %v", err)
	}

	// Test ReadFile
	readData := make([]byte, len(data))
	n, err := ReadFile(filename, readData)
	if err != nil {
		tst.Fatalf("ReadFile failed: %v", err)
	}

	if n != len(data) || string(readData) != string(data) {
		tst.Fatalf("Read data does not match written data")
	}

	// Test ReadFileData
	filedata, err := ReadFileData(filename)
	if err != nil {
		tst.Fatalf("ReadFileData  failed: %v", err)
	}
	fmt.Println(string(filedata))

	// Clean up
	os.Remove(filename)

}

// func TestWriteAndReadFileWithChuncks(t *testing.T) {
// 	filename := "testfile.txt"
// 	data := make([]byte, 1024*1024) // 1 MB
// 	chunkSize := 4096

// 	if err := WriteFileWithChunk(filename, chunkSize, data); err != nil {
// 		t.Fatalf("WriteFile failed: %v", err)
// 	}

// 	readData := make([]byte, chunkSize)
// 	n, err := ReadFileWithChunks(filename, chunkSize, readData)
// 	if err != nil {
// 		t.Fatalf("ReadFile failed: %v", err)
// 	}
// 	if n != len(readData) {
// 		t.Fatalf("Read data does not match written data")
// 	}

// 	// clean up
// 	os.Remove(filename)
// }
