package dispatcher

import (
	"../docker"
	"../processServer"
	"../submitwrap"
	"container/list"
	"log"
	"strconv"
)

type DispatcherConfig struct {
	QueueSize           int                       `json:"queueSize"`
	DispatchChannelSize int                       `json:"channelSize"`
	Network             string                    `json:"network"`
	Ndocker             int                       `json:"ndocker"`
	ClientConfig        docker.DockerClientConfig `json:"clientConfig"`
	ProcessAddrMap      string                    `json:"processAddrMap"`
	ProcessPort         int                       `json:"processPort"`
	ListenAddrMap       string                    `json:"listenAddrMap"`
	ListenPort          int                       `json:"listenPort"`
}

type Dispatcher struct {
	tasks           *list.List
	dispatchChannel chan submitwrap.SubmitTaskWrap

	processChannel <-chan submitwrap.SubmitTaskWrap
	processServer  *processServer.ProcessServer
	processAddrMap string
	processPort    int

	Ndocker         int
	ndockeRemain    int
	dockerClient    *docker.DockerClient
	dockerExecutors map[int]*docker.DockerExecutor
}

func NewDispatcher(config DispatcherConfig, processServer *processServer.ProcessServer, dispatchChannel chan submitwrap.SubmitTaskWrap,
	processChannel <-chan submitwrap.SubmitTaskWrap) *Dispatcher {
	dockerClient, err := docker.NewDockerClient(config.ClientConfig.HttpAddr, config.ClientConfig.Version)
	if err != nil {
		log.Fatal(err)
	}
	return &Dispatcher{
		tasks:           list.New(),
		Ndocker:         config.Ndocker,
		ndockeRemain:    config.Ndocker,
		dispatchChannel: dispatchChannel,
		processChannel:  processChannel,
		dockerClient:    dockerClient,

		dockerExecutors: make(map[int]*docker.DockerExecutor),

		processServer:  processServer,
		processAddrMap: config.ProcessAddrMap,
		processPort:    config.ProcessPort,
	}
}

func (dispatcher *Dispatcher) Run() {
	for {
		var newTask submitwrap.SubmitTaskWrap
		var overTask submitwrap.SubmitTaskWrap
		select {
		case overTask = <-dispatcher.processChannel:
			{
				if overTask.Status == submitwrap.OK || overTask.Status == submitwrap.ERROR {
					log.Println("submit is over", overTask.Task)
					dispatcher.reclaim(overTask)
				} else {
					log.Println("submit Status error", overTask.Task)
				}
			}
		case newTask = <-dispatcher.dispatchChannel:
			{
				dispatcher.tasks.PushBack(newTask)
				log.Println("len of tasks: ", dispatcher.tasks.Len())
			}
		}
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
			dispatcher.processServer.AddSubmit(newTask)

			//now just write gojudgeCore and judgeServer
			cmd := make([]string, 0)
			cmd = append(cmd, "./judgeCore")
			cmd = append(cmd, "--adress="+dispatcher.processAddrMap)
			cmd = append(cmd, "--port="+strconv.Itoa(dispatcher.processPort))
			cmd = append(cmd, "--submitID="+strconv.Itoa(newTask.Task.SubmitID))
			log.Println(cmd)
			resource := docker.Resource{}
			resource.SetMemoryByMb(100)

			dockerExecutor := docker.NewDockerExecutor(dispatcher.processAddrMap,
				strconv.Itoa(dispatcher.processPort), "gojudgecore", "gojudge", submitID2String(newTask.Task.SubmitID),
				cmd, resource, dispatcher.dockerClient)
			err := dockerExecutor.CreateAndStart()
			if err != nil {
				log.Println("docker: ", dispatcher.dockerClient, " ", err)
				dispatcher.processServer.RemoveSubmit(newTask)
				continue
			}
			dispatcher.dockerExecutors[newTask.Task.SubmitID] = dockerExecutor
			newTask.Status = submitwrap.JUDGING
			log.Println(newTask)
		}
		dispatcher.ndockeRemain -= 1
	}
}

func (dispatcher *Dispatcher) reclaim(submitTaskWrap submitwrap.SubmitTaskWrap) {
	dockerExecutor, ok := dispatcher.dockerExecutors[submitTaskWrap.Task.SubmitID]
	if !ok {
		log.Println("dispatcher reclaim error: not found dockerExecutor by key ", submitTaskWrap.Task.SubmitID)
		return
	}
	dockerExecutor.Destroy()
	if dispatcher.ndockeRemain+1 > dispatcher.Ndocker {
		log.Println("reclaim ndockeRemain >= Ndocker")
	} else {
		dispatcher.ndockeRemain = dispatcher.ndockeRemain + 1
	}
}

func submitID2String(submitID int) string {
	return "juge_" + strconv.Itoa(submitID)
}
