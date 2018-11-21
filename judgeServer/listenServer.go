package main

import (
	"fmt"
	"net"
)


type ListenServerConfig struct {
    Port    int     `json:"port"`
}

type ListenServer struct {
    listener    net.Listener
}



func (listenServer *ListenServer) InitServer(listener net.Listener) error {
    listenServer.listener = listener
    return nil
}

func (listenServer *ListenServer) AcceptConn(conn net.Conn) {
		go func(conn net.Conn) {
            for {
                buf := make([]byte, 100)
                n, err := conn.Read(buf)
                if err != nil {
                    conn.Close()
                    break
                }
                fmt.Print(string(buf[:n]))
            }
            fmt.Println("conn Close")
            conn.Close()
        }(conn)
}

func (listenServer *ListenServer) HandleAcceptErorr() error {
    return nil
}

func (listenServer *ListenServer) ExitServer() {

}
