package file_server

import (
	"context"

	"github.com/Nidal-Bakir/file-server/p2p"
)

type Server interface {
	Start(ctx context.Context) error
	ListenAddr() string

	RootName() string

	Nodes(ctx context.Context) ([]p2p.Peer, error)
	AddNode(ctx context.Context, addr string) error

	Close(ctx context.Context) error
	Clear(ctx context.Context) error
}
