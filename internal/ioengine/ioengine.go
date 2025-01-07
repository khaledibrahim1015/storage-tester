package ioengine

import "os"

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
