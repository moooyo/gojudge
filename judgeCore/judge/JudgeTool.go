package judge

import "net"
import "github.com/ferriciron/gojudge/def"
import "github.com/ferriciron/gojudge/moudle"

const (
	SIGKILL      = "signal: killed"
	RuntimeError = "signal: segmentation fault (core dumped)"
	eps          = 0.01
)

func buildResponse(resp *def.Response, code int, msg string) bool {
	resp.ErrCode = code
	resp.Msg = []byte(msg)
	return resp.ErrCode == def.AcceptCode
}

func sendResponse(conn net.Conn, resp *def.Response) {
	encoder := moudle.NewEnCoder()
	socket := moudle.SocketFromConn(conn)
	err := encoder.SendStruct(socket, resp)
	if err != nil {
		panic("Socket Couldn't Send Response to Server!")
	}
}
