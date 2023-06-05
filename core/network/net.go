package network

import (
	"io"
	"net"
)

type Conn interface {
	io.ReadWriteCloser
}

type Dialer interface {
	Dial(network, address string) (Conn, error)
}

type Processor interface {
	Accept() (Conn, error)
	Close() error
}

type Listener interface {
	Listen(network, addr string) (Processor, error)
}

type NetDialer struct {
}

func (n NetDialer) Dial(network, address string) (Conn, error) {
	return net.Dial(network, address)
}

type NetProcessor struct {
	ln net.Listener
}

func (n NetProcessor) Accept() (Conn, error) {
	return n.ln.Accept()
}

func (n NetProcessor) Close() error {
	return n.ln.Close()
}

type NetworkListener struct {
}

func (l NetworkListener) Listen(network, addr string) (Processor, error) {
	ln, err := net.Listen(network, addr)
	if err != nil {
		return NetProcessor{}, err
	}

	return NetProcessor{ln: ln}, nil
}
