package main

type DispatcherConfig struct {
	QueueSize           int `json:"queueSize"`
	DispatchChannelSize int `json:"channelSize"`
	Ndocker             int `json:"ndocker"`
}

type Dispatcher struct {
	Tasks           []SubmitTaskWrap
	Ndocker         int
	dispatchChannel chan SubmitTaskWrap
	processChannel  <-chan SubmitTaskWrap
}

func NewDispatcher(config DispatcherConfig, dispatchChannel chan SubmitTaskWrap, processChannel <-chan SubmitTaskWrap) *Dispatcher {
	return &Dispatcher{
		Tasks:           make([]SubmitTaskWrap, config.QueueSize),
		Ndocker:         config.Ndocker,
		dispatchChannel: dispatchChannel,
		processChannel:  processChannel,
	}
}

func (dispatcher *Dispatcher) Run() {
	for {
		select {
		case _ = <-dispatcher.dispatchChannel:
			{
			}
		case _ = <-dispatcher.processChannel:
			{
			}
		}
		// do dispatch
	}
}
