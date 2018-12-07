package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
)
import "../def"
import "../moudle"
import "./complie"
import "./judge"

var port *string = flag.String("port", "7777", "JudgeServerPort")
var submitId *int = flag.Int("submitID", 0, "submitID")
var adress *string = flag.String("adress", "127.0.0.1", "JugdeServerAdress")

const BasePath = "./problem"
const CompliePath = "./submit"

func main() {
	//parse args
	flag.Parse()
	serverConn, err := net.Dial("tcp", *adress+":"+*port)
	socket:=moudle.SocketFromConn(serverConn)
	//socket := moudle.NewSocket(serverConn)
	encoder:=moudle.NewEnCoder()
	defer socket.Close()
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	encoder.SendInt(socket,*submitId)
	decodere:=moudle.NewDecoder()
	_, err = socket.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	buf.Reset()
	var submit def.Submit
	//err = socket.ReadStruct(&submit)
	err=decodere.ReadStruct(socket,&submit)
	if err != nil {
		log.Fatal(err)
	}
	//complie SourceCode
	problem, err := Complie(&submit)
	var resp def.Response
	if err != nil {
		resp.ErrCode = def.ComplierError
		resp.Msg = []byte(err.Error())
		encoder.SendStruct(socket,&resp)
		return
	}
	//Run
	err = RunJudge(&submit, &problem, serverConn)
	return
}
func Complie(submit *def.Submit) (problem def.Problem, err error) {
	filename := BasePath + "/" + strconv.Itoa(submit.ProblemID) + "/problem.json"
	log.Println(submit.ProblemID)
	log.Println(filename)
	err = ParseProblemFile(filename, &problem)
	if err != nil {
		log.Fatal(fmt.Errorf("Parse problem %s FAILD", filename))
		return
	}
	//complie
	switch submit.Language {
	default:
		return problem, fmt.Errorf("gojudge not support this language")
	case def.CLanguage:
		err = complie.GccComplie(submit)
	}
	return
}

func RunJudge(submit *def.Submit, problem *def.Problem, conn net.Conn) (err error) {
	problemPath:="./problem/"+strconv.Itoa(submit.ProblemID)
	switch submit.Language {
	default:
		return fmt.Errorf("gojudge not support this language")
	case def.CLanguage:
		err = judge.ElfJudge(problemPath,CompliePath, problem, conn)
	case def.Cpp11Language:
		err = judge.ElfJudge(problemPath,CompliePath, problem, conn)
	case def.Cpp17Language:
		err = judge.ElfJudge(problemPath,CompliePath, problem, conn)
	case def.Cpp99Language:
		err = judge.ElfJudge(problemPath,CompliePath, problem, conn)
	}
	return
}
