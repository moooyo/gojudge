package listenServer

import (
	"../defs"
	"fmt"
	"net"
)

func RunListenServer(port string, dispatcheChannel chan <- defs.JudgeTaskWrap) {
	fmt.Println("port ", port)

	addr := net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080}
	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		fmt.Println(err)
	}

	for {
		conn, err := listener.AcceptTCP()
		go func(conn *net.TCPConn) {
			fmt.Println(conn, err)
			buf := make([]byte, 100)
			n, _ := conn.Read(buf)
			fmt.Println("read ", n)
			fmt.Println(string(buf[:n]))
			conn.Close()
		}(conn)
	}
}