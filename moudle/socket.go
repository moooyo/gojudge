package moudle

import (
	"bufio"
	"io"
	"net"
	"time"
)

type Socket struct {
	wbuf *bufio.Writer
	rbuf *bufio.Reader
	conn net.Conn
}

func Dial(addr string) (*Socket, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Socket{
		conn: conn,
		wbuf: bufio.NewWriterSize(conn, 4096),
		rbuf: bufio.NewReaderSize(conn, 1024*16),
	}, nil
}

func (socket *Socket) Close() error {
	return socket.conn.Close()
}

func SocketFromConn(conn net.Conn) *Socket {
	return &Socket{
		conn: conn,
		wbuf: bufio.NewWriterSize(conn, 4096),
		rbuf: bufio.NewReaderSize(conn, 1024*16),
	}
}

func (socket *Socket) SetDeadline(t time.Time) error {
	return socket.conn.SetDeadline(t)
}

func (socket *Socket) SetReadDeadline(t time.Time) error {
	return socket.conn.SetReadDeadline(t)
}

func (socket *Socket) SetWriteDeadline(t time.Time) error {
	return socket.conn.SetWriteDeadline(t)
}

func (socket *Socket) Write(data []byte) (int, error) {
	all := len(data)
	remain := all
	nwrite := 0

	for remain > 0 {
		n, err := socket.wbuf.Write(data[nwrite:all])
		nwrite += n
		if err != nil {
			return nwrite, err
		}
		remain -= n
	}
	return all, socket.Flush()
}

func (socket *Socket) Flush() error {
	for {
		err := socket.wbuf.Flush()
		if err == io.ErrShortWrite {
			continue
		}
		if err == nil {
			return nil
		}
		return err
	}
}

func (socket *Socket) Read(data []byte) (int, error) {
	all := len(data)
	remain := all
	nread := 0

	for remain > 0 {
		n, err := socket.rbuf.Read(data[nread:all])
		nread += n
		if err != nil {
			return nread, err
		}
		remain -= n
	}
	return all, nil
}

func (socket *Socket) RemoteAddr() net.Addr {
	return socket.conn.RemoteAddr()
}

func (socket *Socket) LocalAddr() net.Addr {
	return socket.conn.LocalAddr()
}
