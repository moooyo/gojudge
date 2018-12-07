package moudle

import (
	"../def"
	"log"
	"net"
	"reflect"
	"testing"
)

func server(listen net.Listener) {
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		socket := SocketFromConn(conn)

		coder := NewDECoderWithSize(1024 * 1024 * 2)

		for i := 0; i < 1000; i++ {
			var submit def.Submit
			err = coder.ReadStruct(socket, &submit)
			if err != nil {
				log.Fatal(err)
			}

			coder.AppendStruct(&submit)

			coder.Send(socket)
		}

		socket.Close()
	}
}

func TestSocket(t *testing.T) {

	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Fatal(err)
	}

	go server(listen)

	socket, err := Dial("127.0.0.1:8080")
	if err != nil {
		t.Fatal(err)
	}

	coder := NewDECoderWithSize(2 * 1024 * 1024)

	defer socket.Close()

	for i := 0; i < 1000; i++ {
		submit := def.Submit{
			SubmitID:   i,
			ProblemID:  i,
			CodeSource: make([]byte, 1024*5),
			Language:   i,
		}
		coder.AppendStruct(&submit)

		coder.Send(socket)

		var result def.Submit

		err = coder.ReadStruct(socket, &result)

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(result, submit) {
			t.Fail()
		}
	}
}
