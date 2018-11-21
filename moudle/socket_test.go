package  moudle

import (
	"../def"
	"net"
	"reflect"
	"testing"
)

const count int  =  10000
const addr  string = "127.0.0.1:8080"

func StartTestServer(addr string, sync chan <- int) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	sync <- 1
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go func(conn net.Conn) {
            socket := NewSocket(conn)
			for i := 0; i < count; i++ {
				var resp def.Submit
                socket.StructRead(&resp)
                socket.StructWrite(&resp)
			}
			conn.Close()
		}(conn)

	}

}
func TestReadWriteStruct (t *testing.T) {

	sync := make(chan  int, 1)
	go StartTestServer(addr, sync)

	<- sync
	conn, _ := net.Dial("tcp", addr)
    socket := NewSocket(conn)
	for i := 0; i < count; i++ {
		var test = def.Submit{
			SubmitID:   i,
			ProblemID:  i,
			CodeSource: []byte("Hello, world"),
			Language:   i,
		}

        socket.StructWrite(&test)
        


		var resp def.Submit

        socket.StructRead(&resp)
    
		if !reflect.DeepEqual(resp, test) {
			t.Error("send not equal recv")
		}
	}
	conn.Close()
}
