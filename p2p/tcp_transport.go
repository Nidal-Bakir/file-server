package p2p

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
)

type TcpTransportParams struct {
	ListenerAddress string
	Shakehands      HandshakeFn
	Decoder         Decoder
	OnNewPeer       func(Peer) error
}

type TcpTransport struct {
	TcpTransportParams

	listener   net.Listener
	consumChan chan *Message
	closeChan  chan struct{}
}

func NewTcpTransport(params TcpTransportParams) Transport {
	return &TcpTransport{
		TcpTransportParams: params,
		closeChan:          make(chan struct{}),
		consumChan:         make(chan *Message),
	}
}

func (t *TcpTransport) Consume() <-chan *Message {
	return t.consumChan
}

func (t *TcpTransport) Close() {
	t.consumChan <- nil
	t.closeChan <- struct{}{}
}

func (t *TcpTransport) ListenAndAccept(ctx context.Context) (err error) {
	t.listener, err = net.Listen("tcp", t.ListenerAddress)
	if err != nil {
		return err
	}
	return t.startAcceptLoop(ctx)
}

func (t *TcpTransport) startAcceptLoop(ctx context.Context) error {
	go func() {
		<-t.closeChan
		t.listener.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			t.Close()
			return context.Cause(ctx)
		default:
			con, err := t.listener.Accept()
			if err != nil {
				fmt.Printf("tcp: error can not accept new connection, err: %v\n", err)
				return err
			}
			go t.handleTcpConnection(con)
		}
	}
}

func (t *TcpTransport) handleTcpConnection(conn net.Conn) {
	peer := NewTcpPeerFromAccept(conn)
	if t.OnNewPeer != nil {
		if err := t.OnNewPeer(peer); err != nil {
			peer.Conn().Write([]byte(err.Error()))
			peer.Conn().Close()
			return
		}
	}

	if t.Shakehands != nil {
		if err := t.Shakehands(peer); err != nil {
			peer.Conn().Write([]byte(fmt.Sprintln(ErrInvalidHandshake.Error())))
			peer.Conn().Close()
			return
		}
	}

	fmt.Printf("new connection accepted and hundled. RemoteAddr:%s \n", peer.Conn().RemoteAddr().String())

	msg := new(Message)
	msg.From = peer.Conn().RemoteAddr()
	for {
		err := t.Decoder.Decode(peer.Conn(), msg)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error from tcp decoder, err: %v\n", err)
		}
		t.consumChan <- msg
	}
}
