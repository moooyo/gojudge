package moudle


import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
    "bufio"
	"net"
	"unsafe"
)
import "../def"
var intTemp int
const intSize int64=int64(unsafe.Sizeof((intTemp)))

type Socket struct {
    reader *bufio.Reader
    writer *bufio.Writer
    conn   net.Conn
}

func NewSocket(conn net.Conn) *Socket {
    return &Socket {
        reader: bufio.NewReader(conn),
        writer: bufio.NewWriter(conn),
        conn: conn,
    }
}

func (socket *Socket) SocketRead() (data []byte, err error) {
	temp:=make([]byte,intSize)
    socket.reader.Read(temp)
	if err!=nil{
		return nil ,err
	}
	dataReader:=io.LimitReader(socket.reader,int64(binary.LittleEndian.Uint64(temp)))
    buf := make([]byte, binary.LittleEndian.Uint64(temp))
	_,err=dataReader.Read(buf)
	return buf,err
}

func (socket *Socket) SocketWrite(data []byte) (size int, err error) {
    size = len(data)
	var buf=bytes.NewBuffer(make([]byte,0))
	binary.Write(buf,binary.LittleEndian, int64(size))
	binary.Write(buf,binary.LittleEndian,data)
	size,err= socket.writer.Write(buf.Bytes())
    socket.writer.Flush()
	return size,err
}

func (socket *Socket) StructWrite(resp def.SocketInterface)(err error){
	temp,err:=resp.StructToBytes()
	if err!=nil{
		return err
	}
	_,err= socket.SocketWrite(temp)
	return err
}

func (socket *Socket) StructRead(resp interface{})(err error){
	data,err:= socket.SocketRead()
	if err!=nil{
		return err
	}
	err =json.Unmarshal(data,&resp)
	return nil
}
