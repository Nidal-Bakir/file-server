package p2p

import (
	"context"
	"io"
	"net"
)

// Peer is an interface that represent the remote node
type Peer interface {
	io.ReadWriteCloser
	RemoteAddr() net.Addr

	// Outbound peer (we dialed and acquired the connection) => true
	// Inbound peer (we accepted the connection) => false
	IsOutBound() bool
}

// Transport is anything that can handel the communication
// between peers in the network, this could be tcp, udp, ftp, websockets, ....
type Transport interface {
	ListenAndAccept(ctx context.Context) (err error)
	Consume() <-chan *Message
	Close() error
	Dial(addr string) error
}
