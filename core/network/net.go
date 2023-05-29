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

type Listener interface {
	Accept() (Conn, error)
	Close() error
}

type ListenFunc func(network, laddr string) (Listener, error)

type NetDialer struct {
}

func (n NetDialer) Dial(network, address string) (Conn, error) {
	return net.Dial(network, address)
}

type NetListener struct {
	ln net.Listener
}

func (n *NetListener) Accept() (Conn, error) {
	return n.ln.Accept()
}

func (n *NetListener) Close() error {
	return n.ln.Close()
}

func Listen(network, addr string) (Listener, error) {
	ln, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}

	return &NetListener{ln: ln}, nil
}
