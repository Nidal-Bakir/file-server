package p2p

import (
	"io"
	"net"
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type Message struct {
	From    net.Addr
	Payload []byte
}
