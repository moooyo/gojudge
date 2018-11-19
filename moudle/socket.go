package moudle


import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"unsafe"
)

var intTemp int

const intSize int64=int64(unsafe.Sizeof((intTemp)))

func socketRead(conn net.Conn)(buf []byte,err error) {
	sizeReader:=io.LimitReader(conn,intSize)
	temp:=make([]byte,intSize)
	sizeReader.Read(temp)
	if err!=nil{
		return nil ,err
	}
	dataReader:=io.LimitReader(conn,int64(binary.LittleEndian.Uint64(temp)))
	buf = make([]byte, binary.LittleEndian.Uint64(temp))
	_,err=dataReader.Read(buf)
	return buf,err
}

func socketWrite(conn net.Conn,data []byte)(size int,err error){
	size=len(data)
	var buf=bytes.NewBuffer(make([]byte,0))
	binary.Write(buf,binary.LittleEndian, int64(size))
	binary.Write(buf,binary.LittleEndian,data)
	size,err=conn.Write(buf.Bytes())
	return size,err
}
