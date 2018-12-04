package moudle

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
	"time"
	"unsafe"
)
import "../def"

var intTemp int

const intSize int64 = int64(unsafe.Sizeof((intTemp)))

type Socket struct {
	reader *bufio.Reader
	writer *bufio.Writer
	conn   net.Conn
}

func NewSocket(conn net.Conn) *Socket {
	return &Socket{
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		conn:   conn,
	}
}

func NewSocketWithSize(conn net.Conn, rs, ws int) *Socket {
	return &Socket{
		reader: bufio.NewReaderSize(conn, rs),
		writer: bufio.NewWriterSize(conn, ws),
		conn:   conn,
	}
}

func (socket *Socket) LocalAddr() net.Addr {
	return socket.conn.LocalAddr()
}

func (socket *Socket) RemoteAddr() net.Addr {
	return socket.conn.RemoteAddr()
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

func (socket *Socket) Close() {
	socket.writer.Flush()
	socket.conn.Close()
}

func (socket *Socket) SocketRead() (data []byte, err error) {

	temp := make([]byte, intSize)

	_, err = socket.Read(temp)

	if err != nil {
		return nil, err
	}

	buf := make([]byte, binary.LittleEndian.Uint64(temp))

	_, err = socket.Read(buf)

	if err != nil {
		return nil, err
	}

	return buf, err
}

func (socket *Socket) SocketWrite(data []byte) (size int, err error) {
	size = len(data)
	var buf = bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, int64(size))
	binary.Write(buf, binary.LittleEndian, data)
	size, err = socket.Write(buf.Bytes())
	return size, err
}

func (socket *Socket) Read(p []byte) (n int, err error) {
	psize := len(p)
	premain := psize

	for premain > 0 {
		nreaded, err := socket.reader.Read(p[psize-premain : psize])
		if err != nil {
			return psize - premain, err
		}
		premain = premain - nreaded
	}

	return psize, nil
}

func (socket *Socket) Write(p []byte) (n int, err error) {

	n, err = socket.WriteWithOutFlush(p)

	socket.writer.Flush()

	return n, err
}

func (socket *Socket) WriteWithOutFlush(p []byte) (n int, err error) {

	psize := len(p)
	premain := psize

	for premain > 0 {
		nwrited, err := socket.writer.Write(p[psize-premain : psize])
		if err != nil {
			return psize - premain, err
		}
		premain = premain - nwrited
	}
	return psize, nil
}

func (socket *Socket) WriteStruct(resp def.SocketInterface) (err error) {

	temp, err := resp.StructToBytes()

	if err != nil {
		return err
	}

	_, err = socket.SocketWrite(temp)

	return err
}

func (socket *Socket) ReadStruct(resp interface{}) (err error) {

	data, err := socket.SocketRead()

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &resp)
	return nil
}
