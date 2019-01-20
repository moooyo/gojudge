package dispatcher

import (
	"container/list"
	"github.com/ferriciron/gojudge/judgeServer/processServer"
	"github.com/ferriciron/gojudge/judgeServer/submitwrap"
	"log"
	_ "strconv"
)

type DispatcherConfig struct {
	QueueSize           int `json:"queueSize"`
	DispatchChannelSize int `json:"channelSize"`
	Ndocker             int `json:"ndocker"`
}

type Dispatcher struct {
	tasks           *list.List
	dispatchChannel chan submitwrap.SubmitTaskWrap
	processServer   *processServer.ProcessServer
	executorChannel chan<- submitwrap.SubmitTaskWrap
	ndockeRemain    int
}

func NewDispatcher(config DispatcherConfig, processServer *processServer.ProcessServer, dispatchChannel chan submitwrap.SubmitTaskWrap, executorChannel chan<- submitwrap.SubmitTaskWrap) *Dispatcher {
	return &Dispatcher{
		tasks:           list.New(),
		dispatchChannel: dispatchChannel,
		processServer:   processServer,
		executorChannel: executorChannel,
		ndockeRemain:    config.Ndocker,
	}
}

var count int = 0

func (dispatcher *Dispatcher) dealTask() {
	task := <-dispatcher.dispatchChannel
	switch task.Status {
	case submitwrap.OK:
		{
			dispatcher.reclaim(task)
		}
	case submitwrap.TIMEOUTERROR:
		{
		}
	case submitwrap.EXECUTERROR:
		{
			dispatcher.ndockeRemain--
			dispatcher.processServer.RemoveSubmit(task)
			log.Println("task executor error")
		}
	case submitwrap.WAITING:
		{
			dispatcher.tasks.PushBack(task)
		}
	case submitwrap.UNKONW:
		{
		}
	default:
		{
		}
	}
}

func (dispatcher *Dispatcher) Run() {
	for {
		dispatcher.dealTask()
		dispatcher.dispatch()
	}
}

func (dispatcher *Dispatcher) dispatch() {
	for dispatcher.ndockeRemain > 0 {
		if dispatcher.tasks.Len() <= 0 {
			return
		}
		item := dispatcher.tasks.Front()
		dispatcher.tasks.Remove(item)
		if newTask, ok := item.Value.(submitwrap.SubmitTaskWrap); ok {
			// do some time cal
			newTask.Status = submitwrap.EXECUTING
			dispatcher.processServer.AddSubmit(newTask)

			dispatcher.executorChannel <- newTask
		}
		dispatcher.ndockeRemain--
	}
}

func (dispatcher *Dispatcher) reclaim(submitTaskWrap submitwrap.SubmitTaskWrap) {
	count++
	log.Println("count: ", count)
	dispatcher.ndockeRemain++
}
