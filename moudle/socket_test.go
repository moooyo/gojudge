package moudle

import (
	"../def"
	"fmt"
	"net"
	"reflect"
	"testing"
)

const count int = 1000000

func StartTestServer(addr string, sync chan<- int) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		sync <- 1
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept err")
			panic(err)
		}
		go func(conn net.Conn) {
			socket := NewSocket(conn)
			for i := 0; i < count; i++ {
				data, err := socket.SocketRead()

				if err != nil {
					panic(err)
				}
				//data = data[8:len(data)]
				socket.SocketWrite(data)
			}
			conn.Close()
		}(conn)

	}

}

func TestReadWriteData(t *testing.T) {

	sync := make(chan int, 1)
	addr := "127.0.0.1:8081"
	go StartTestServer(addr, sync)

	data := make([]byte, 1024)
	<-sync
	conn, _ := net.Dial("tcp", addr)

	socket := NewSocketWithSize(conn, 1024*1024*4, 1024*1024*4)

	for i := 0; i < count; i++ {

		n, err := socket.SocketWrite(data)

		if err != nil {
			t.Error(n, err)
		}

		resp, err := socket.SocketRead()

		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(data, resp) {
			t.Error("send not equal recv")
		}
	}
	conn.Close()
}

func TestReadWriteStruct(t *testing.T) {
	sync := make(chan int, 1)
	addr := "127.0.0.1:8080"
	go StartTestServer(addr, sync)

	<-sync
	conn, _ := net.Dial("tcp", addr)

	socket := NewSocket(conn)

	for i := 0; i < count; i++ {

		var data = def.Submit{
			i,
			i,
			[]byte("Hello, world"),
			i,
		}

		err := socket.WriteStruct(&data)

		if err != nil {
			t.Error(err)
		}

		var resp def.Submit

		err = socket.ReadStruct(&resp)

		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(data, resp) {
			t.Error("send not equal recv")
		}
	}
	conn.Close()
}
