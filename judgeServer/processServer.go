package main

import (
	"fmt"
	"net"
	"strconv"
)

type ProcessServerConfig struct {
	Port int
}

type ProcessServer struct {
	taskMap        map[int]SubmitTaskWrap
	processChannel chan SubmitTaskWrap
	listener       net.Listener
	addr           string
}

func NewProcessServer(config ProcessServerConfig, processChannel chan SubmitTaskWrap) *ProcessServer {
	return &ProcessServer{
		taskMap:        make(map[int]SubmitTaskWrap),
		processChannel: processChannel,
		addr:           "0.0.0.0:" + strconv.Itoa(config.Port),
	}
}

func (processServer *ProcessServer) Addr() string {
	return processServer.addr
}

func (processServer *ProcessServer) InitServer(listener net.Listener) error {
	processServer.listener = listener
	return nil
}

func (processServer *ProcessServer) AcceptConn(conn net.Conn) {
	fmt.Println("processServer incoming")
	conn.Close()
}

func (processServer *ProcessServer) HandleAcceptErorr() error {
	return nil
}

func (processServer *ProcessServer) ExitServer() {
}
