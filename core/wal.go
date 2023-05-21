package core

import (
	"io"

	"gitub.com/sriramr98/go_kvdb/utils"
)

type Log struct {
}

type WalLogger struct {
	serializer utils.Serializer[Log]
	writer     io.WriteCloser
}

func NewWalLogger(writer io.WriteCloser, serializer utils.Serializer[Log]) *WalLogger {
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

func (wc *WalLogger) Close() error {
	return wc.writer.Close()
}
