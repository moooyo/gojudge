package  moudle

import (
	"../def"
	"encoding/json"
	"net"
	"reflect"
	"testing"
)

const count int  =  1000
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
			for i := 0; i < count; i++ {
				data, _ := socketRead(conn)
				var resp def.Submit
				json.Unmarshal(data, &resp)
				data, _ = resp.StructToBytes()
				socketWrite(conn, data)
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
	for i := 0; i < count; i++ {
		var test = def.Submit{
			SubmitID:   i,
			ProblemID:  i,
			CodeSource: []byte("Hello, world"),
			Language:   i,
		}

		data, _ := test.StructToBytes()
		socketWrite(conn, data)
		data, _ = socketRead(conn)

		var resp def.Submit

		json.Unmarshal(data ,&resp)

		if !reflect.DeepEqual(resp, test) {
			t.Error("send not equal recv")
		}
	}
	conn.Close()
}
