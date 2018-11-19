package moudle

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"unsafe"
)
var intTemp int
const intSize int64=int64(unsafe.Sizeof((intTemp)))
func socketRead(conn net.Conn,data []byte)(size int,err error) {
	sizeReader:=io.LimitReader(conn,intSize)
	temp:=make([]byte,intSize)
	_,err=sizeReader.Read(temp)
	if err!=nil{
		return 0,err
	}
	dataReader:=io.LimitReader(conn,int64(binary.BigEndian.Uint64(temp)))
	var buf []byte
	size,err=dataReader.Read(buf)
	return size,err
}

func socketWrite(conn net.Conn,data []byte)(size int,err error){
		size=len(data)
		var buf=bytes.NewBuffer(make([]byte,0))
		binary.Write(buf,binary.BigEndian,intSize)
		binary.Write(buf,binary.BigEndian,data)
		size,err=conn.Write(buf.Bytes())
		return size,err
}
