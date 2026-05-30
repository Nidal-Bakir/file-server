package file_server

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Nidal-Bakir/file-server/storage"
	"github.com/Nidal-Bakir/file-server/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileServerAndClose(t *testing.T) {
	ctx := context.Background()
	server := utilTestNewFileServer(t, "server_3000", ":3000")
	go func() {
		time.Sleep(time.Second)
		teardown(t, server, ctx)
	}()
	err := server.Start(ctx)
	assert.ErrorIs(t, err, net.ErrClosed)
}

func TestStartNewFileServerAndAcceptConn(t *testing.T) {
	ctx := context.Background()
	addr := ":3000"
	dirName := "server_3000"
	server := utilTestNewFileServer(t, dirName, addr)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		require.ErrorIs(t, server.Start(ctx), net.ErrClosed)
	})
	time.Sleep(time.Millisecond * 200)

	conn := utilTestDialTcp(t, addr)
	defer conn.Close()
	conn.Write([]byte("hi from: " + dirName))
	time.Sleep(time.Millisecond * 200)

	teardown(t, server, ctx)
	wg.Wait()
}

func TestStartNewFileServersAndAcceptConn(t *testing.T) {
	ctx := context.Background()
	server1 := utilTestNewFileServer(t, "server_3001", ":3001")
	server2 := utilTestNewFileServer(t, "server_3002", ":3002")

	wg := sync.WaitGroup{}
	wg.Go(func() {
		require.ErrorIs(t, server1.Start(ctx), net.ErrClosed)
	})
	wg.Go(func() {
		require.ErrorIs(t, server2.Start(ctx), net.ErrClosed)
	})
	time.Sleep(time.Millisecond * 200)

	conn1 := utilTestDialTcp(t, server1.ListenAddr())
	defer conn1.Close()
	conn1.Write([]byte("hi from: " + server1.RootName()))

	conn2 := utilTestDialTcp(t, server2.ListenAddr())
	defer conn2.Close()
	conn2.Write([]byte("hi from: " + server2.RootName()))

	time.Sleep(time.Millisecond * 200)
	teardown(t, server2, ctx)
	teardown(t, server1, ctx)
	wg.Wait()
}

func TestConnectionToOtherNodes(t *testing.T) {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	server1 := utilTestNewFileServer(t, "server_3001", "127.0.0.1:3001")
	wg.Go(func() {
		require.ErrorIs(t, server1.Start(ctx), net.ErrClosed)
	})
	time.Sleep(time.Millisecond * 200)

	server2 := utilTestNewFileServer(t, "server_3002", "127.0.0.1:3002", server1.ListenAddr())
	wg.Go(func() {
		require.ErrorIs(t, server2.Start(ctx), net.ErrClosed)
	})
	time.Sleep(time.Millisecond * 200)

	peers2, err := server2.Nodes(ctx)
	assert.NoError(t, err)
	isConnectedToServer1 := false
	for _, p := range peers2 {
		if p.RemoteAddr().String() == server1.ListenAddr() {
			isConnectedToServer1 = true
			break
		}
	}
	assert.True(t, isConnectedToServer1)

	err = server1.AddNode(ctx, server2.ListenAddr())
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 200)
	peers1, err := server1.Nodes(ctx)
	assert.NoError(t, err)
	isConnectedToServer2 := false
	for _, p := range peers1 {
		if p.RemoteAddr().String() == server2.ListenAddr() {
			isConnectedToServer2 = true
			break
		}
	}
	assert.True(t, isConnectedToServer2)

	teardown(t, server1, ctx)
	teardown(t, server2, ctx)
	wg.Wait()
}

func utilTestNewFileServer(t *testing.T, dirName string, addr string, nodes ...string) *defaultServer {
	opt := ServerOpt{
		ListenAddr:      addr,
		PathTransformer: storage.CASPathTransform,
		Root:            utilCreateRoot(t, dirName),
		Nodes:           nodes,
	}
	server := New(opt).(*defaultServer)
	return server
}

func utilCreateRoot(t *testing.T, dirName string) *os.Root {
	root, err := utils.OpenRoot(dirName)
	require.NoError(t, err)
	return root
}

func utilTestDialTcp(t *testing.T, addr string) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		require.NoError(t, err)
	}
	go func() {
		m := make([]byte, 1024)
		for {
			n, err := conn.Read(m)
			if err != nil {
				fmt.Printf("client recived an error: %v \n", err)
				conn.Close()
				break
			}
			fmt.Printf("client recived a message: %s \n", string(m[:n]))
		}
	}()
	return conn
}

func teardown(t *testing.T, s *defaultServer, ctx context.Context) {
	assert.NoError(t, s.Clear(ctx))
	assert.NoError(t, s.Close(ctx))
	assert.NoError(t, os.RemoveAll(s.RootName()))
}
