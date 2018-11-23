package main

import (
	"fmt"
	"net"
	"strconv"
)

type ListenServerConfig struct {
	Port int `json:"port"`
}

type ListenServer struct {
	listener          net.Listener
	addr              string
	dispatcherChannel chan<- SubmitTaskWrap
}

func NewListenServer(config ListenServerConfig, dispatcherChannel chan<- SubmitTaskWrap) *ListenServer {
	return &ListenServer{
		addr:              "127.0.0.1:" + strconv.Itoa(config.Port),
		dispatcherChannel: dispatcherChannel,
	}
}

func (listenServer *ListenServer) InitServer(listener net.Listener) error {
	listenServer.listener = listener
	return nil
}

func (listenServer *ListenServer) AcceptConn(conn net.Conn) {
	go func(conn net.Conn) {
		fmt.Println("listenServer incoming")
		conn.Close()
	}(conn)
}

func (listenServer *ListenServer) HandleAcceptErorr() error {
	return nil
}

func (listenServer *ListenServer) ExitServer() {

}

func (listenServer *ListenServer) Addr() string {
	return listenServer.addr
}
