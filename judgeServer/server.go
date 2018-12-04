package main

import (
	"log"
	"net"
)

type Server interface {
	InitServer(net.Listener) error

	AcceptConn(conn net.Conn)

	ExitServer()

	HandleAcceptErorr() error

	Addr() string
}

func RunServer(server Server) {

	listener, err := net.Listen("tcp", server.Addr())

	if err != nil {
		log.Fatalln("Listen: ", server.Addr(), ": ", err)
	}

	server.InitServer(listener)

	for {

		conn, err := listener.Accept()

		// err EMFILE ENFILE  no more fd
		// ECONNABORTED connection has been aborted

		if err != nil {
			err = server.HandleAcceptErorr()
			if err != nil {
				server.ExitServer()
				return
			}
		} else {
			server.AcceptConn(conn)
		}
	}
}
