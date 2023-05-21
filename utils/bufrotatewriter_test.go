package utils

import (
	"os"
	"path"
	"testing"
)

func helperTestFileWriter(writer *BuffRotateFileWriter, data [][]byte, t *testing.T) {
	t.Helper()
	for _, data := range data {
		n, err := writer.Write(data)
		if err != nil {
			t.Fatalf("Write error: %v\n", err)
		}

		//we're checking for len(data)+1 because of the newline character added when written to file
		if n != len(data)+1 {
			t.Errorf("Unexpected number of bytes written: got %d, want %d\n", n, len(data))

		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close BuffRotateFileWriter: %v", err)
	}
}

func getLogDir(t *testing.T) string {
	t.Helper()
	dir := os.TempDir()
	log_dir := path.Join(dir, "kvdb_test_log")

	_, err := os.Stat(log_dir)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(log_dir, 0777)
			return log_dir
		}

		t.Fatalf("Failed to stat log directory: %v\n", err)
	}

	return log_dir
}

func TestBuffRotateFileWriter(t *testing.T) {

	bytesThreshold := int64(10) // 10 bytes
	log_dir := getLogDir(t)
	defer os.RemoveAll(log_dir)

	writer, err := NewBuffRotateFileWriter(log_dir, bytesThreshold)
	if err != nil {
		t.Fatalf("Failed to create BuffRotateFileWriter: %v\n", err)
	}
	data_to_write := [][]byte{[]byte("this is a test"), []byte("this is a test")}

	helperTestFileWriter(writer, data_to_write, t)

	total_data_size := 0
	for _, data := range data_to_write {
		total_data_size += len(data)
	}
	noFilesExpected := total_data_size / int(bytesThreshold)

	files, err := os.ReadDir(log_dir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v\n", err)
	}

	if len(files) != noFilesExpected {
		t.Errorf("Unexpected number of files: got %d, want %d\n", len(files), noFilesExpected)
	}

}

func TestBufRotateWriter_NewLineOnEveryWrite(t *testing.T) {
	bytesThreshold := int64(10) // 10 bytes
	log_dir := getLogDir(t)
	defer os.RemoveAll(log_dir)

	writer, err := NewBuffRotateFileWriter(log_dir, bytesThreshold)
	if err != nil {
		t.Fatalf("Failed to create BuffRotateFileWriter: %v\n", err)
	}

	defer os.RemoveAll(path.Join(log_dir, "*"))
	data_to_write := [][]byte{[]byte("ab"), []byte("df"), []byte("gh")}

	helperTestFileWriter(writer, data_to_write, t)

	files, err := os.ReadDir(log_dir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v\n", err)
	}

	if len(files) != 1 {
		t.Errorf("Unexpected number of files: got %d, want %d\n", len(files), 1)
	}

	file, err := os.ReadFile(path.Join(log_dir, files[0].Name()))
	if err != nil {
		t.Fatalf("Failed to read file: %v\n", err)
	}

	expected := "ab\ndf\ngh\n"
	if string(file) != expected {
		t.Errorf("Unexpected file contents: got %s, want %s\n", string(file), expected)
	}

}
