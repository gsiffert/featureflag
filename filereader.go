package featureflag

import (
	"context"
	"fmt"
	"io"
	"os"
)

type fileReader struct {
	filePath string
}

func (f *fileReader) Read(_ context.Context) (bytes []byte, err error) {
	var file *os.File
	file, err = os.Open(f.filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if e := file.Close(); e != nil {
			err = fmt.Errorf("close file: %w", e)
		}
	}()

	bytes, err = io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return bytes, nil
}

// NewFileReader returns a SourceReader that reads the content of the file at filePath.
func NewFileReader(filePath string) SourceReader[[]byte] {
	return &fileReader{filePath: filePath}
}
