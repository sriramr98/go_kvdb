package core

import (
	"io"
	"sync"

	"gitub.com/sriramr98/go_kvdb/utils"
)

type Log struct {
}

type WalLogger struct {
	mu         sync.Mutex
	serializer utils.Serializer[Log]
	writer     io.Writer
}

func NewWalLogger(writer io.Writer, serializer utils.Serializer[Log]) *WalLogger {
	return &WalLogger{
		writer:     writer,
		serializer: serializer,
	}
}

func (wl *WalLogger) Write(log Log) error {

	b, err := wl.serializer.Serialize(log)
	if err != nil {
		return err
	}

	_, err = wl.writer.Write(b)
	return err

}

func (wl *WalLogger) Close() {
	wl.mu.Lock()
	defer wl.mu.Unlock()
}
