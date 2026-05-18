package p2p

import (
	"context"
	"net"
)

// Peer is an interface that represent the remote node
type Peer interface {
	Conn() net.Conn
	IsOutBound() bool
}

// Transport is anything that can handel the communication
// between peers in the network, this could be tcp, udp, ftp, websockets, ....
type Transport interface {
	ListenAndAccept(ctx context.Context) (err error)
	Consume() <-chan *Message
	Close()
}
