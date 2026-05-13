package main

import (
	"context"
	"log"

	"github.com/Nidal-Bakir/file-server/p2p"
)

func main() {
	ctx := context.Background()
	params := p2p.TcpTransportParams{
		Shakehands:      p2p.NoOpHandshake,
		ListenerAddress: ":4000",
		Decoder:         p2p.NewDefaultEncoding(),
	}
	transport := p2p.NewTcpTransport(params)
	if err := transport.ListenAndAccept(ctx); err != nil {
		log.Fatal(err)
	}
}
