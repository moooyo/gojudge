package dispatchServer

import (
    "../defs"
	"fmt"
)

type Dispatcher struct {
	Tasks [] defs.JudgeTaskWrap
}


func RunDispatcher(port string, dispatchChannel chan defs.JudgeTaskWrap) {
	dispatcher := Dispatcher{make([] defs.JudgeTaskWrap, 10)}
	for {
		task := <- dispatchChannel
		dispatcher.Tasks = append(dispatcher.Tasks, task)
		fmt.Println("dispatch", task)
		fmt.Println(*task.Task)
	}
}
