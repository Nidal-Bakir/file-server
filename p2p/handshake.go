package p2p

import "errors"

var (
	ErrInvalidHandshake = errors.New("invalid handshake")
)

type HandshakeFn func(Peer) error

func NoOpHandshake(_ Peer) error { return nil }
