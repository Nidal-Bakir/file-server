package file_server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Nidal-Bakir/file-server/p2p"
	"github.com/Nidal-Bakir/file-server/storage"
)

type ServerOpt struct {
	ListenAddr      string
	PathTransformer storage.PathTransformFunc
	Root            *os.Root
	Nodes           []string
}

type defaultServer struct {
	opt       ServerOpt
	transport p2p.Transport
	storage   storage.Storage

	peersMux sync.RWMutex
	peers    map[string]p2p.Peer
}

func New(opt ServerOpt) Server {
	s := &defaultServer{
		opt:      opt,
		peers:    map[string]p2p.Peer{},
		peersMux: sync.RWMutex{},
		storage: storage.NewSimpleStorage(
			storage.SimpleStoreParam{
				RootDir:           opt.Root,
				PathTransformFunc: storage.CASPathTransform,
			},
		),
	}
	s.transport = p2p.NewTcpTransport(
		p2p.TcpTransportParams{
			ListenerAddress: opt.ListenAddr,
			Decoder:         p2p.NewDefaultEncoding(),
			OnPeer:          s.onPeer,
		})
	return s
}

func (s *defaultServer) onPeer(p p2p.Peer) error {
	s.peersMux.Lock()
	defer s.peersMux.Unlock()
	s.peers[p.RemoteAddr().String()] = p
	if p.IsOutBound() {
		fmt.Printf("connected to remote node: %s, using: %s\n", p.RemoteAddr().String(), p.RemoteAddr().Network())
	} else {
		fmt.Printf("new client connected to this node: %s, using: %s\n", p.RemoteAddr().String(), p.RemoteAddr().Network())
	}
	return nil
}

func (s *defaultServer) ListenAddr() string {
	return s.opt.ListenAddr
}

func (s *defaultServer) RootName() string {
	return s.opt.Root.Name()
}

func (s *defaultServer) Clear(ctx context.Context) error {
	return s.storage.Clear()
}

func (s *defaultServer) Close(ctx context.Context) error {
	return errors.Join(
		s.transport.Close(),
		s.storage.Close(),
	)
}

func (s *defaultServer) Start(ctx context.Context) error {
	go s.consumeNetworkMessages(ctx)
	s.connectToNodes()
	return s.transport.ListenAndAccept(ctx)
}

func (s *defaultServer) Nodes(ctx context.Context) ([]p2p.Peer, error) {
	s.peersMux.RLock()
	defer s.peersMux.RUnlock()

	peers := make([]p2p.Peer, 0)
	if s.peers == nil {
		return peers, nil
	}
	for _, p := range s.peers {
		if p.IsOutBound() {
			peers = append(peers, p)
		}
	}
	return peers, nil
}

func (s *defaultServer) AddNode(ctx context.Context, addr string) error {
	err := s.transport.Dial(addr)
	if err != nil {
		fmt.Println("error, can not connect to node: ", addr)
		return err
	}
	return nil
}

func (s *defaultServer) consumeNetworkMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, ok := <-s.transport.Consume()
			if !ok || msg == nil {
				return
			}
			err := s.storage.StreamStore("TODO:", bytes.NewBuffer(msg.Payload))
			if err != nil {
				fmt.Println("streamStore: error can not store the payload")
			}
		}
	}
}

func (s *defaultServer) connectToNodes() {
	if len(s.opt.Nodes) == 0 {
		return
	}
	var err error
	for _, n := range s.opt.Nodes {
		err = s.transport.Dial(n)
		if err != nil {
			fmt.Println("error, can not connect to node: ", n)
		}
	}
}
