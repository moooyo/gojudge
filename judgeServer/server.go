package main


import (
    "net"
)

type Server interface {

    InitServer(net.Listener) error

    AcceptConn(conn net.Conn)

    ExitServer()

    HandleAcceptErorr() error
}


func RunServer(server Server, addr string) {

    listener, err := net.Listen("tcp", addr)
   

    if err != nil {
        return
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
