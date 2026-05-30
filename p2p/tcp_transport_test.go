package p2p

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const _ListenerAddress = ":4000"

func TestTcpTransport(t *testing.T) {
	tt := testUtilsCreateTcpTransport(newMockedEncoding())
	assert.Equal(t, tt.ListenerAddress, _ListenerAddress)

	con := testUtilsDialTcp()
	assert.Equal(t, con.RemoteAddr().String(), "127.0.0.1"+_ListenerAddress)
	con.Close()
	time.Sleep(time.Millisecond * 300)
	tt.Close()
	time.Sleep(time.Millisecond * 300)
}

func TestTcpSeding(t *testing.T) {
	mockedEnc := newMockedEncoding()
	tt := testUtilsCreateTcpTransport(mockedEnc)

	conn := testUtilsDialTcp()
	conn.Write([]byte("test1"))
	conn.Write([]byte("test2"))
	conn.Close()
	time.Sleep(time.Millisecond * 300)

	tt.Close()
	time.Sleep(time.Millisecond * 300)

	mockedEnc.AssertExpectations(t)
}

func testUtilsDialTcp() net.Conn {
	conn, err := net.Dial("tcp", _ListenerAddress)
	if err != nil {
		panic(err)
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

func testUtilsCreateTcpTransport(decoder Decoder) *TcpTransport {
	ctx := context.Background()
	params := TcpTransportParams{
		ListenerAddress: _ListenerAddress,
		Decoder:         decoder,
	}
	tt := NewTcpTransport(params).(*TcpTransport)
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ctx done for consume in main")
				return

			case msg, ok := <-tt.Consume():
				if !ok || msg == nil {
					return
				}
				fmt.Printf("new message from chan %+v\n", msg)
			}
		}
	}()
	go tt.ListenAndAccept(ctx)
	time.Sleep(time.Second)
	return tt
}

type mockedEncoding struct {
	mock.Mock
}

func (m *mockedEncoding) Decode(r io.Reader, msg *Message) error {
	args := m.Called(r, msg)
	return args.Error(0)
}

func newMockedEncoding() *mockedEncoding {
	testObj := new(mockedEncoding)
	testObj.On("Decode", mock.Anything, mock.Anything).Return(nil).Times(2)
	testObj.On("Decode", mock.Anything, mock.Anything).Return(io.EOF).Once()
	return testObj
}
