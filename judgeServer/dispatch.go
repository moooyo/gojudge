package main

import (
	"fmt"
)

type Dispatcher struct {
	Tasks [] JudgeTaskWrap
}


func RunDispatcher(port string, dispatchChannel chan JudgeTaskWrap) {
	dispatcher := Dispatcher{make([] JudgeTaskWrap, 10)}
	for {
		task := <- dispatchChannel
		dispatcher.Tasks = append(dispatcher.Tasks, task)
		fmt.Println("dispatch", task)
		fmt.Println(task)
	}
}
