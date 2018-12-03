package main

import (
	"./dispatcher"
	"./listenServer"
	"./processServer"
	"./submitwrap"
)

func RunSystem(conf Config) {

	dispatcherChannel := make(chan submitwrap.SubmitTaskWrap, conf.DispatcherConfig.DispatchChannelSize)

	processChannel := make(chan submitwrap.SubmitTaskWrap, conf.DispatcherConfig.Ndocker)

	processServer := processServer.NewProcessServer(conf.ProcessConfig, processChannel)
	dispatcher := dispatcher.NewDispatcher(conf.DispatcherConfig, processServer, dispatcherChannel, processChannel)

	go dispatcher.Run()
	go RunServer(processServer)

	listenServer := listenServer.NewListenServer(conf.ListenConfig, dispatcherChannel)

	RunServer(listenServer)
}
