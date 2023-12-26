package utils

import (
	"fmt"
	"io"
	"os"
)

func Upload(filePath string, fileData []byte) error {
	// Create or open the file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// Write the []byte content to the file
	_, err = file.Write(fileData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	return nil
}

func Download(filePath string) ([]byte, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return nil, err
	}
	fileSize := fileInfo.Size()

	// Read the file content
	content := make([]byte, fileSize)
	_, err = file.Read(content)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading file:", err)
		return nil, err
	}
	return content, nil
}
