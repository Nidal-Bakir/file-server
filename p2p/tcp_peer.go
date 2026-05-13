package p2p

import "net"

type TcpPeer struct {
	conn net.Conn

	// Outbound peer (we dialed and acquired the connection) => true
	// Inbound peer (we accepted the connection) => false
	isOutBound bool
}

func NewTcpPeer(conn net.Conn, isOutBound bool) Peer {
	return &TcpPeer{
		conn:       conn,
		isOutBound: isOutBound,
	}
}

func NewTcpPeerFromAccept(conn net.Conn) Peer {
	return NewTcpPeer(conn, false)
}

func NewTcpPeerFromDial(conn net.Conn) Peer {
	return NewTcpPeer(conn, true)
}

func (p *TcpPeer) Conn() net.Conn {
	return p.conn
}

func (p *TcpPeer) IsOutBound() bool {
	return p.isOutBound
}
