package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"
	"net"
	"unsafe"
)

var port *string=flag.String("port","7777","JudgeServerPort")
var submitId *int=flag.Int("submitID",0,"submitID")
var adress *string=flag.String("adress","127.0.0.1","JugdeServerAdress")

func main(){
	//parse args
	flag.Parse()
	serverConn,err:=net.Dial("tcp",*adress+":"+*port)
	defer serverConn.Close()
	if err!=nil {
		log.Fatal(err)
	}

	buf:=bytes.NewBuffer(make([]byte,0))
	//网络字节序为大端序
	binary.Write(buf,binary.BigEndian,unsafe.Sizeof(*submitId))
	binary.Write(buf,binary.BigEndian,*submitId)
	_,err=serverConn.Write(buf.Bytes())
	if err!=nil{
		log.Fatal(err)
	}
	buf.Reset()
	

}

