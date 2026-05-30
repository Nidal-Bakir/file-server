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
	Decoder         Decoder

	Handshake HandshakeFn
	OnPeer  func(Peer) error
}

type TcpTransport struct {
	TcpTransportParams

	listener   net.Listener
	consumeChan chan *Message
	closeChan  chan struct{}
}

func NewTcpTransport(params TcpTransportParams) Transport {
	return &TcpTransport{
		TcpTransportParams: params,
		closeChan:          make(chan struct{}),
		consumeChan:         make(chan *Message),
	}
}

func (t *TcpTransport) Consume() <-chan *Message {
	return t.consumeChan
}

func (t *TcpTransport) Close() error {
	close(t.consumeChan)
	close(t.closeChan)
	return nil
}

func (t *TcpTransport) ListenAndAccept(ctx context.Context) (err error) {
	t.listener, err = net.Listen("tcp", t.ListenerAddress)
	if err != nil {
		return err
	}
	return t.startAcceptLoop(ctx)
}

func (t *TcpTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleTcpConnection(conn, true)
	return nil
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
			go t.handleTcpConnection(con, false)
		}
	}
}

func (t *TcpTransport) handleTcpConnection(conn net.Conn, isOutBound bool) {
	peer := NewTcpPeer(conn, isOutBound)
	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			peer.Write([]byte(err.Error()))
			peer.Close()
			return
		}
	}

	if t.Handshake != nil {
		if err := t.Handshake(peer); err != nil {
			peer.Write([]byte(fmt.Sprintln(ErrInvalidHandshake.Error())))
			peer.Close()
			return
		}
	}

	fmt.Printf("new connection accepted and hundled. RemoteAddr:%s \n", peer.RemoteAddr().String())

	msg := new(Message)
	msg.From = peer.RemoteAddr()
	for {
		err := t.Decoder.Decode(peer, msg)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error from tcp decoder, err: %v\n", err)
		}
		t.consumeChan <- msg
	}
}
