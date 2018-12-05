package listenServer

import (
	"../../def"
	"../../moudle"
	"../submitwrap"
	"log"
	"net"
)

type ListenServerConfig struct {
	ListenAddr string `json:"listenAddr"`
}

type ListenServer struct {
	listener          net.Listener
	addr              string
	dispatcherChannel chan<- submitwrap.SubmitTaskWrap
}

func NewListenServer(config ListenServerConfig, dispatcherChannel chan<- submitwrap.SubmitTaskWrap) *ListenServer {
	return &ListenServer{
		addr:              config.ListenAddr,
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
			log.Println("AcceptConn read ", err)
			socket.Close()
			return
		}
		log.Println("New submit from web front: ", &submit)
		socket.Close()
		listenServer.dispatcherChannel <- submitwrap.WrapSubmit(&submit)
	}(socket)
}

func (listenServer *ListenServer) HandleAcceptErorr() error {
	log.Println("HandleAcceptErorr error")
	return nil
}

func (listenServer *ListenServer) ExitServer() {

}

func (listenServer *ListenServer) Addr() string {
	return listenServer.addr
}
