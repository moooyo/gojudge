package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"unsafe"
)
import "../def"
import "../moudle"
import "./complie"
var port *string=flag.String("port","7777","JudgeServerPort")
var submitId *int=flag.Int("submitID",0,"submitID")
var adress *string=flag.String("adress","127.0.0.1","JugdeServerAdress")
const BasePath="./problem"
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

	var submit def.Submit
	err=moudle.StructRead(serverConn,&submit)
	if err!=nil{
		log.Fatal(err)
	}

	//complie SourceCode
	err=Complie(&submit)
	var resp def.Response
	if err!=nil{
		//todo
		resp.ErrCode=def.ComplierError
		resp.Msg=[]byte(err.Error())
		moudle.StructWrite(serverConn,&resp)
		return
	}

	//Run

}
func Complie(submit *def.Submit)(err error){
	filename:=BasePath+string(submit.ProblemID)+"/problem.json"
	var problem Problem
	err=ParseProblemFile(filename,&problem)
	if err!=nil{
		log.Fatal(fmt.Errorf("Parse problem %s FAILD",filename))
		return
	}
	//complie
	switch submit.Language{
	default:
		return fmt.Errorf("gojudge not support this language")
	case def.CLanguage:
		err=complie.GccComplie(submit)
	}
	return
}

