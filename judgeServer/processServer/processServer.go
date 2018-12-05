package processServer

import (
	"../../def"
	"../../moudle"
	"../submitwrap"
	"encoding/binary"
	"log"
	"net"
	"sync"
)

type ProcessServerConfig struct {
	ListenAddr string `json:"listenAddr"`
}

type ProcessServer struct {
	taskMap        map[int]submitwrap.SubmitTaskWrap
	processChannel chan submitwrap.SubmitTaskWrap
	listener       net.Listener
	mutex          sync.Mutex
	addr           string
}

func NewProcessServer(config ProcessServerConfig, processChannel chan submitwrap.SubmitTaskWrap) *ProcessServer {
	return &ProcessServer{
		taskMap:        make(map[int]submitwrap.SubmitTaskWrap),
		processChannel: processChannel,
		addr:           config.ListenAddr,
	}
}

func (processServer *ProcessServer) AddSubmit(submitTaskWrap submitwrap.SubmitTaskWrap) {
	processServer.mutex.Lock()
	processServer.taskMap[submitTaskWrap.Task.SubmitID] = submitTaskWrap
	processServer.mutex.Unlock()
}

func (processServer *ProcessServer) RemoveSubmit(submitTaskWrap submitwrap.SubmitTaskWrap) {
	processServer.mutex.Lock()
	delete(processServer.taskMap, submitTaskWrap.Task.SubmitID)
	processServer.mutex.Unlock()
}

func (processServer *ProcessServer) CheckoutSubmit(submitID int) (submitwrap.SubmitTaskWrap, bool) {
	processServer.mutex.Lock()
	submit, ok := processServer.taskMap[submitID]
	if ok {
		processServer.mutex.Unlock()
		return submit, true
	}
	processServer.mutex.Unlock()
	return submit, false
}

func (processServer *ProcessServer) Addr() string {
	return processServer.addr
}

func (processServer *ProcessServer) InitServer(listener net.Listener) error {
	processServer.listener = listener
	return nil
}

func (processServer *ProcessServer) AcceptConn(conn net.Conn) {
	socket := moudle.NewSocket(conn)
	go func(socket *moudle.Socket) {
		data := make([]byte, 8)
		_, err := socket.Read(data)
		if err != nil {
			log.Println("processServer: ", err)
			socket.Close()
			return
		}
		temp := uint64(binary.LittleEndian.Uint64(data))
		tempdata := make([]byte, temp)
		socket.Read(tempdata)
		submitid := int(binary.LittleEndian.Uint32(tempdata))
		submitTaskWrap, ok := processServer.CheckoutSubmit(submitid)
		if !ok {
			log.Println("processServer access a bad submit")
			socket.Close()
			return
		}
		socket.WriteStruct(submitTaskWrap.Task)
		for {
			var resp def.Response
			err := socket.ReadStruct(&resp)
			if err != nil {
				submitTaskWrap.Status = submitwrap.ERROR
				log.Println("judgeCore error")
				break
			} else {
				log.Println(&resp)
			}
		}
		socket.Close()
		processServer.RemoveSubmit(submitTaskWrap)
		processServer.processChannel <- submitTaskWrap
	}(socket)
}

func (processServer *ProcessServer) HandleAcceptErorr() error {
	return nil
}

func (processServer *ProcessServer) ExitServer() {
}
