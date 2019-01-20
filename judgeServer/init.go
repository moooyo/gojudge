package main

import (
	"fmt"
	"github.com/ferriciron/gojudge/judgeServer/dispatcher"
	"github.com/ferriciron/gojudge/judgeServer/executor"
	"github.com/ferriciron/gojudge/judgeServer/listenServer"
	"github.com/ferriciron/gojudge/judgeServer/processServer"
	"github.com/ferriciron/gojudge/judgeServer/submitwrap"
)

func RunSystem(conf Config) {

	fmt.Println(conf)

	dispatcherChannel := make(chan submitwrap.SubmitTaskWrap, conf.DispatcherConfig.DispatchChannelSize)

	processServer := processServer.NewProcessServer(conf.ProcessConfig, dispatcherChannel)

	executorChannel := make(chan submitwrap.SubmitTaskWrap, conf.ExecutorConfig.ChannelSize)

	dispatcher := dispatcher.NewDispatcher(conf.DispatcherConfig, processServer, dispatcherChannel, executorChannel)

	executor := executor.NewExecutor(executorChannel, dispatcherChannel, conf.ExecutorConfig)

	listenServer := listenServer.NewListenServer(conf.ListenConfig, dispatcherChannel)

	executor.Run()

	go dispatcher.Run()

	go RunServer(processServer)

	listenServer.Run()
}
