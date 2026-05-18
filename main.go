package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Nidal-Bakir/file-server/p2p"
)

func main() {
	ctx := context.Background()
	params := p2p.TcpTransportParams{
		ListenerAddress: ":4000",
		Decoder:         p2p.NewDefaultEncoding(),
	}
	transport := p2p.NewTcpTransport(params)

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ctx done for consume in main")
				return

			case msg := <-transport.Consume():
				if msg == nil {
					return
				}
				fmt.Printf("new message from chan %+v\n", msg)
			}
		}
	}()

	if err := transport.ListenAndAccept(ctx); err != nil {
		log.Fatal(err)
	}
}
