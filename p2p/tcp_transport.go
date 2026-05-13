package p2p

import (
	"context"
	"fmt"
	"net"
	"sync"
)

type TcpTransportParams struct {
	ListenerAddress string
	Shakehands      HandshakeFn
	Decoder         Decoder
}

type TcpTransport struct {
	TcpTransportParams

	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTcpTransport(params TcpTransportParams) Transport {
	return &TcpTransport{
		TcpTransportParams: params,
	}
}

func (t *TcpTransport) ListenAndAccept(ctx context.Context) (err error) {
	t.listener, err = net.Listen("tcp", t.ListenerAddress)
	if err != nil {
		return err
	}
	return t.startAcceptLoop(ctx)
}

func (t *TcpTransport) startAcceptLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			t.listener.Close()
			return context.Cause(ctx)
		default:
			con, err := t.listener.Accept()
			if err != nil {
				fmt.Println("tcp: error can not accept new connection")
				continue
			}
			go t.handleTcpConnection(con)
		}
	}
}

func (t *TcpTransport) handleTcpConnection(conn net.Conn) {
	peer := NewTcpPeerFromAccept(conn)
	if err := t.Shakehands(peer); err != nil {
		peer.Conn().Write([]byte(fmt.Sprintln(ErrInvalidHandshake.Error())))
		peer.Conn().Close()
		return
	}
	fmt.Printf("new connection accepted and hundled. RemoteAddr:%s \n", peer.Conn().RemoteAddr().String())

	msg := new(Message)
	msg.From = peer.Conn().RemoteAddr()
	for {
		err := t.Decoder.Decode(peer.Conn(), msg)
		if err != nil {
			fmt.Println("error from tcp decoder")
			break
		}
		fmt.Printf("%v \n", msg)
	}
}
