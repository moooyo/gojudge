package main

import (
	"./defs"
	"./dispatchServer"
	"./listenServer"
	"./processServer"
	"fmt"
)



func main() {
	/*
		get config listenServer port
		get config dispatchServer port
	 */
	dispatchChannel := make(chan defs.JudgeTaskWrap, 1024)
	fmt.Println(dispatchChannel)
	go dispatchServer.RunDispatcher("8080", dispatchChannel)

	var i defs.JudgeTask = 22222

	dispatchChannel <- defs.JudgeTaskWrap{
		Status:    defs.OK,
		Task: &i,
	}

	listenServer.RunListenServer("8081", dispatchChannel)

	processServer.RunProcessServer("8081", dispatchChannel)
	dispatchChannel <- defs.JudgeTaskWrap{
		Status:    defs.OK,
		Task: &i,
	}

	dispatchChannel <- defs.JudgeTaskWrap{
		Status:    defs.OK,
		Task: &i,
	}
	//b, _ :=json.Marshal(value)
	//s := string(b)
	//fmt.Println(s)
	// go run listenServer
	// channel
	// go run dispatchServer
	// go run processServer a signal service
	//fmt.Scanln()
	//fmt.Println("done")

}
