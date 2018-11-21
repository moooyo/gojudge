package main

import (
	"fmt"
)

type DispatcherConfig struct {
    QueueSize   int     `json:"queueSize"`
    DispatchChannelSize int `json:"channelSize"`
}

type Dispatcher struct {
	Tasks [] SubmitTaskWrap
}


func RunDispatcher(port string, dispatchChannel chan SubmitTaskWrap) {
	dispatcher := Dispatcher{make([] SubmitTaskWrap, 10)}
	for {
		task := <- dispatchChannel
		dispatcher.Tasks = append(dispatcher.Tasks, task)
		fmt.Println("dispatch", task)
		fmt.Println(task)
	}
}
