package judge

import "net"
import "../../def"
import "../../moudle"

const (
	SIGKILL      = "signal: killed"
	RuntimeError = "signal: segmentation fault (core dumped)"
	eps = 0.01
)


func buildResponse(resp *def.Response, code int, msg string) bool {
	resp.ErrCode = code
	resp.Msg = []byte(msg)
	return resp.ErrCode == def.AcceptCode
}

func sendResponse(conn net.Conn, resp *def.Response) {
	socket := moudle.NewSocket(conn)
	socket.WriteStruct(resp)
}