package p2p

import "net"

type TcpPeer struct {
	conn       net.Conn
	isOutBound bool
}

func NewTcpPeer(conn net.Conn, isOutBound bool) Peer {
	return &TcpPeer{
		conn:       conn,
		isOutBound: isOutBound,
	}
}

func (p *TcpPeer) IsOutBound() bool {
	return p.isOutBound
}

func (p *TcpPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

func (p *TcpPeer) Read(buff []byte) (n int, err error) {
	return p.conn.Read(buff)
}

func (p *TcpPeer) Write(buff []byte) (n int, err error) {
	return p.conn.Write(buff)
}

func (p *TcpPeer) Close() error {
	return p.conn.Close()
}
