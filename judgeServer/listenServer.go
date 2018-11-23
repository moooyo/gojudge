package main

import (
	"../def"
	"../moudle"
	"log"
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
		addr:              "0.0.0.0:" + strconv.Itoa(config.Port),
		dispatcherChannel: dispatcherChannel,
	}
}

func (listenServer *ListenServer) InitServer(listener net.Listener) error {
	listenServer.listener = listener
	return nil
}

func (listenServer *ListenServer) AcceptConn(conn net.Conn) {
	socket := moudle.NewSocket(conn)
	go func(socket *moudle.Socket) {
		var submit def.Submit
		err := socket.ReadStruct(&submit)
		if err != nil {
			log.Println(err)
			socket.Close()
			return
		}
		log.Println("New submit from web front: ", submit)
		listenServer.dispatcherChannel <- WrapSubmit(&submit)
		socket.Close()
	}(socket)
}

func (listenServer *ListenServer) HandleAcceptErorr() error {
	return nil
}

func (listenServer *ListenServer) ExitServer() {

}

func (listenServer *ListenServer) Addr() string {
	return listenServer.addr
}
