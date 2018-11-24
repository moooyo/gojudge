package main

import (
	"./dispatcher"
	"./docker"
	"./listenServer"
	"./processServer"
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	configCase := Config{
		ProcessConfig: processServer.ProcessServerConfig{
			ListenAddr: "judgeServer:8081",
		},
		ListenConfig: listenServer.ListenServerConfig{
			ListenAddr: "judgeServer:8080",
		},
		DispatcherConfig: dispatcher.DispatcherConfig{
			QueueSize:           1024,
			DispatchChannelSize: 1024,
			Ndocker:             8,
			Network:             "maymomo",
			ClientConfig: docker.DockerClientConfig{
				HttpAddr: "http://172.17.0.1:8118",
				Version:  "v1.37",
			},
			ProcessAddrMap: "judgeServer",
			ProcessPort:    8081,
			ListenAddrMap:  "judgeServer",
			ListenPort:     8080,
		},
	}

	configPath := "./config.json"
	config, err := ParseConfig(configPath)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(configCase, config) {
		t.Error("test case: ", configCase, "\n not equal\n", "result :", config)
	}
}
