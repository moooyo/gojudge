package main

func RunSystem(conf Config) {

	dispatcherChannel := make(chan SubmitTaskWrap, conf.DispatcherConfig.DispatchChannelSize)

	processChannel := make(chan SubmitTaskWrap, conf.DispatcherConfig.Ndocker)

	dispatcher := NewDispatcher(conf.DispatcherConfig, dispatcherChannel, processChannel)

	go dispatcher.Run()

	processServer := NewProcessServer(conf.ProcessConfig, processChannel)

	go RunServer(processServer)

	listenServer := NewListenServer(conf.ListenConfig, dispatcherChannel)

	RunServer(listenServer)
}
