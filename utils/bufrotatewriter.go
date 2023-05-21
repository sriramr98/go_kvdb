package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// This helps us to write a rotated file writer which is also buffered
// All writes will first go into an in-memory buffer of configured size
// Once the buffer is full, the buffer will be flushed to the file
// If the file size exceeds the configured size, the file will be rotated
type BuffRotateFileWriter struct {
	dirpath           string        // The path of the directory where all the files are stored
	mu                sync.Mutex    // Mutex to ensure thread safety
	writer            *bufio.Writer // buffered writer for the current file
	file              *os.File      // The current file in the buffered writer, required to close file when rotating
	bytesWritten      int64         // The number of bytes written to the current file
	fileBytesTreshold int64         // The number of bytes after which we rotate the file
	lastFileIndex     int           // The index of the last file
}

func NewBuffRotateFileWriter(dirpath string, bytesTreshold int64) (*BuffRotateFileWriter, error) {
	dir, err := os.Stat(dirpath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err = os.MkdirAll(dirpath, 0755); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else if !dir.IsDir() {
		return nil, errors.New("path is not a directory")
	}

	lastFileIndex := 0
	file, err := os.Create(filepath.Join(dirpath, strconv.Itoa(lastFileIndex)))
	if err != nil {
		return nil, err
	}

	bufWriter := bufio.NewWriterSize(file, 1024*1024) // 1MB default in-mem buffer size

	return &BuffRotateFileWriter{
		dirpath:           dirpath,
		writer:            bufWriter,
		file:              file,
		lastFileIndex:     lastFileIndex,
		fileBytesTreshold: bytesTreshold,
	}, nil
}

// returns the amount of bytes written and error if any
func (brfw *BuffRotateFileWriter) Write(data []byte) (int, error) {
	brfw.mu.Lock()
	defer brfw.mu.Unlock()

	data_to_write := []byte(fmt.Sprintf("%s\n", data))
	totalBytesWritten := brfw.bytesWritten + int64(len(data_to_write))
	if totalBytesWritten > brfw.fileBytesTreshold && brfw.bytesWritten > 0 {
		if err := brfw.rotateFile(); err != nil {
			return 0, err
		}
	}

	n, err := brfw.file.Write([]byte(data_to_write))
	brfw.bytesWritten += int64(n)
	return n, err
}

func (brfw *BuffRotateFileWriter) rotateFile() error {
	if brfw.file != nil {
		err := brfw.writer.Flush()
		if err != nil {
			return err
		}

		err = brfw.file.Close()
		if err != nil {
			return err
		}
	}

	fmt.Printf("creating file with index %d", brfw.lastFileIndex+1)
	file, err := os.Create(filepath.Join(brfw.dirpath, strconv.Itoa(brfw.lastFileIndex+1)))
	if err != nil {
		return err
	}

	brfw.file = file
	brfw.writer.Reset(brfw.file)
	brfw.lastFileIndex++
	brfw.bytesWritten = 0

	return nil
}

func (brfw *BuffRotateFileWriter) Close() error {
	brfw.mu.Lock()
	defer brfw.mu.Unlock()

	err := brfw.writer.Flush()
	if err != nil {
		return err
	}

	err = brfw.file.Close()
	if err != nil {
		return err
	}

	// this is done so that the next write will create a new file
	brfw.file = nil

	return nil
}
