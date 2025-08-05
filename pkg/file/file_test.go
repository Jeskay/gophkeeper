package file

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	testDataSize   = 1024 // Size of random text to generate
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789	 "
	testFolderName = "test_files"
)

func generateRandomText(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

func setupTestFolder(t *testing.T) string {
	t.Helper()
	// Create test folder in temporary directory
	tmpDir := os.TempDir()
	testPath := filepath.Join(tmpDir, testFolderName)
	err := os.MkdirAll(testPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	return testPath
}

func cleanupTestFolder(t *testing.T, path string) {
	t.Helper()
	err := os.RemoveAll(path)
	if err != nil {
		t.Errorf("Failed to cleanup test directory: %v", err)
	}
}

func TestFileReadWrite(t *testing.T) {
	testPath := setupTestFolder(t)
	defer cleanupTestFolder(t, testPath)

	// Initialize reader and writer
	reader := NewFileReader()
	writer := NewFileWriter()

	testCases := []struct {
		name     string
		dataSize int
	}{
		{"Small file", 100},
		{"Medium file", 1024},
		{"Large file", 10240},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate random text
			expectedText := generateRandomText(tc.dataSize)
			fileName := filepath.Join(testPath, fmt.Sprintf("test_%d.txt", tc.dataSize))

			// Write test
			file, err := writer.CreateFile(fileName)
			if err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}

			n, err := writer.FileWriteString(file, expectedText)
			if err != nil {
				t.Fatalf("Failed to write to file: %v", err)
			}
			if n != len(expectedText) {
				t.Errorf("Written bytes %d doesn't match expected length %d", n, len(expectedText))
			}

			err = writer.FileClose(file)
			if err != nil {
				t.Fatalf("Failed to close file after writing: %v", err)
			}

			// Read test
			file, err = reader.OpenFile(fileName)
			if err != nil {
				t.Fatalf("Failed to open file for reading: %v", err)
			}

			data := make([]byte, tc.dataSize)
			n, err = reader.FileRead(file, data)
			if err != nil {
				t.Fatalf("Failed to read from file: %v", err)
			}
			if n != len(expectedText) {
				t.Errorf("Read bytes %d doesn't match expected length %d", n, len(expectedText))
			}

			actualText := string(data[:n])
			if actualText != expectedText {
				t.Errorf("Read content doesn't match written content.\nExpected: %s\nGot: %s",
					expectedText[:50], actualText[:50]) // Show first 50 chars for debugging
			}

			err = reader.FileClose(file)
			if err != nil {
				t.Fatalf("Failed to close file after reading: %v", err)
			}
		})
	}
}

func TestBufferedReadWrite(t *testing.T) {
	testPath := setupTestFolder(t)
	defer cleanupTestFolder(t, testPath)

	reader := NewFileReader()
	writer := NewFileWriter()

	testCases := []struct {
		name     string
		dataSize int
	}{
		{"Small buffered file", 100},
		{"Medium buffered file", 1024},
		{"Large buffered file", 10240},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedText := generateRandomText(tc.dataSize)
			fileName := filepath.Join(testPath, fmt.Sprintf("buffered_%d.txt", tc.dataSize))

			// Buffered write test
			file, err := writer.CreateFile(fileName)
			if err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}

			buffWriter := writer.NewBufferedWriter(file)
			n, err := writer.BufferedWriteString(buffWriter, expectedText)
			if err != nil {
				t.Fatalf("Failed to write to buffered writer: %v", err)
			}
			if n != len(expectedText) {
				t.Errorf("Written bytes %d doesn't match expected length %d", n, len(expectedText))
			}

			err = writer.BufferedFlush(buffWriter)
			if err != nil {
				t.Fatalf("Failed to flush buffered writer: %v", err)
			}

			err = writer.FileClose(file)
			if err != nil {
				t.Fatalf("Failed to close file after writing: %v", err)
			}

			// Buffered read test
			file, err = reader.OpenFile(fileName)
			if err != nil {
				t.Fatalf("Failed to open file for reading: %v", err)
			}

			buffReader := reader.NewBufferedReader(file)
			var actualText string
			for {
				chunk, err := reader.BufferedRead(buffReader)
				if err != nil {
					break
				}
				actualText += string(chunk)
			}

			if actualText != expectedText {
				t.Errorf("Read content doesn't match written content.\nExpected: %s\nGot: %s",
					expectedText[:50], actualText[:50]) // Show first 50 chars for debugging
			}

			err = reader.FileClose(file)
			if err != nil {
				t.Fatalf("Failed to close file after reading: %v", err)
			}
		})
	}
}
