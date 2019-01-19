package processServer

import (
	"github.com/ferriciron/gojudge/def"
	"github.com/ferriciron/gojudge/judgeServer/submitwrap"
	"github.com/ferriciron/gojudge/moudle"
	"log"
	"net"
	"sync"
)

type ProcessServerConfig struct {
	ListenAddr  string `json:"listenAddr"`
	ChannelSize int    `json:"channelSize"`
}

type ProcessServer struct {
	taskMap           map[int]submitwrap.SubmitTaskWrap
	dispathcerChannel chan<- submitwrap.SubmitTaskWrap
	listener          net.Listener
	mutex             sync.Mutex
	addr              string
}

func NewProcessServer(config ProcessServerConfig, dispathcerChannel chan<- submitwrap.SubmitTaskWrap) *ProcessServer {
	return &ProcessServer{
		taskMap:           make(map[int]submitwrap.SubmitTaskWrap),
		dispathcerChannel: dispathcerChannel,
		addr:              config.ListenAddr,
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
	socket := moudle.SocketFromConn(conn)
	go func(socket *moudle.Socket) {
		coder := moudle.NewDECoderWithSize(1024 * 4)
		submitid, err := coder.ReadInt(socket)
		if err != nil {
			log.Println(err)
			socket.Close()
			return
		}
		submitTaskWrap, ok := processServer.CheckoutSubmit(submitid)
		if !ok {
			log.Println("processServer access a bad submit")
			socket.Close()
			return
		}
		coder.SendStruct(socket, submitTaskWrap.Task)
		for {
			var resp def.Response
			err := coder.ReadStruct(socket, &resp)
			if err != nil {
				log.Println("judgeCore error")
				break
			} else {
				log.Println(&resp)
				if resp.ErrCode != def.AcceptCode {

				} else if resp.AllNode != resp.JudgeNode {
					continue
				} else {
				}
				submitTaskWrap.Status = submitwrap.OK
				break
			}
		}
		socket.Close()
		processServer.RemoveSubmit(submitTaskWrap)
		processServer.dispathcerChannel <- submitTaskWrap
	}(socket)
}

func (processServer *ProcessServer) HandleAcceptErorr() error {
	return nil
}

func (processServer *ProcessServer) ExitServer() {
}
