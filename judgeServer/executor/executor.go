package executor

import (
	"github.com/ferriciron/gojudge/judgeServer/submitwrap"
	"log"
	"strconv"
	"time"
)

type ExecutorConfig struct {
	Workers          int                   `json:"workers"`
	RecliamFrequence int                   `json:"recliamFrequence"`
	ChannelSize      int                   `json:"channelSize"`
	ContainerConfig  DockerContainerConfig `json:"containerConfig"`
	ClientConfig     DockerClientConfig    `json:"clientConfig"`
}

type Executor struct {
	workers []worker
}

var clientConfig DockerClientConfig
var containerConfig DockerContainerConfig
var executorConfig ExecutorConfig

type worker interface {
	run()
}

type createWorker struct {
	client          *DockerClient
	executorChannel <-chan submitwrap.SubmitTaskWrap
	dispatchChannel chan<- submitwrap.SubmitTaskWrap
}

type recliamWorker struct {
	client *DockerClient
	ticker *time.Ticker
}

func NewExecutor(executorChannel chan submitwrap.SubmitTaskWrap,
	dispatchChannel chan<- submitwrap.SubmitTaskWrap, config ExecutorConfig) *Executor {
	containerConfig = config.ContainerConfig
	clientConfig = config.ClientConfig
	executorConfig = config

	workers := make([]worker, 0)

	for i := 0; i < config.Workers; i++ {
		client, err := NewDockerClient(clientConfig.HttpAddr, clientConfig.Version)
		if err != nil {
			log.Fatal(err)
		}
		w := &createWorker{
			client:          client,
			executorChannel: executorChannel,
			dispatchChannel: dispatchChannel,
		}
		workers = append(workers, w)
	}

	client, err := NewDockerClient(clientConfig.HttpAddr, clientConfig.Version)
	if err != nil {
		log.Fatal(err)
	}
	rworker := &recliamWorker{
		client: client,
		ticker: time.NewTicker(time.Duration(config.RecliamFrequence) * time.Second),
	}
	workers = append(workers, rworker)
	return &Executor{
		workers: workers,
	}
}

func (executor *Executor) Run() {
	for i := 0; i < len(executor.workers); i++ {
		executor.workers[i].run()
	}
}

func (w *createWorker) run() {
	go func() {
		for {
			task := <-w.executorChannel
			if task.Status != submitwrap.EXECUTING {
				log.Println("BUG ON createWorker task's status error")
			}
			err := w.executeTask(task)
			if err != nil {
				task.Status = submitwrap.EXECUTERROR
				w.dispatchChannel <- task
			}
		}
	}()
}

func (w *createWorker) executeTask(task submitwrap.SubmitTaskWrap) error {
	var resource Resource
	resource.cpu = 1
	resource.SetMemoryByMb(200)
	dockerExecutor := NewDockerExecutor(task.Task.SubmitID, task2Name(task), resource, w.client, containerConfig)
	err := dockerExecutor.CreateAndStart()
	return err
}

func (w *recliamWorker) run() {
	go func() {
		for {
			<-w.ticker.C
			log.Println("remove exited containers")
			w.client.RemoveByLabel(containerConfig.LabelName, containerConfig.LabelValue)
		}
	}()
}

func task2Name(task submitwrap.SubmitTaskWrap) string {
	return "judge_" + strconv.Itoa(task.Task.SubmitID) + "_" + strconv.Itoa(task.ExecCount)

}
