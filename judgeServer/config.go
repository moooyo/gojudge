package main

import (
	"./dispatcher"
	"./executor"
	"./listenServer"
	"./processServer"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ProcessConfig    processServer.ProcessServerConfig `json:"processConfig"`
	ListenConfig     listenServer.ListenServerConfig   `json:"listenConfig"`
	DispatcherConfig dispatcher.DispatcherConfig       `json:"dispatcherConfig"`
	ExecutorConfig   executor.ExecutorConfig           `json:"executorConfig"`
}

func ParseConfig(configPath string) (cfg Config, err error) {

	configContent, err := ioutil.ReadFile(configPath)

	var config Config

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(configContent, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
