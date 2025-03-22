package utils

import (
	"io"
	"os"
	"strings"
)

// TestData is a wrapped test file.
type TestData struct {
	FilePath   string
	OriginData []byte
}

// NewTestData create a new TestData.
func NewTestData(
	filePath string,
) (TestData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return TestData{}, err
	}
	return TestData{
		FilePath:   filePath,
		OriginData: data,
	}, nil
}

// Data returns a content.
func (d TestData) Data() ([]byte, error) {
	return os.ReadFile(d.FilePath)
}

// Restore writes the original content back to the file.
func (d TestData) Restore() error {
	newlineChar := "\n"
	lines := strings.Split(string(d.OriginData), newlineChar)
	return WriteLinesToExistingFile(d.FilePath, lines, newlineChar)
}

func WriteLinesToExistingFile(
	fileName string,
	lines []string,
	newlineChar string,
) error {
	data := strings.Join(lines, newlineChar)
	return WriteExistingFile(
		fileName,
		[]byte(data),
	)
}

// WriteExistingFile writes the byte array to an existing file.
func WriteExistingFile(
	fileName string,
	data []byte,
) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
