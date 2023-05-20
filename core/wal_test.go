package core

import "testing"

type MockWriter struct {
	writeCount int
}

func (mw *MockWriter) Write(p []byte) (n int, err error) {
	mw.writeCount++
	return len(p), nil
}

func (mw *MockWriter) Clear() {
	mw.writeCount = 0
}

func TestWalLogger(t *testing.T) {

}
