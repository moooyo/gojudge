package executor

import (
	"testing"
)

var config DockerClientConfig = DockerClientConfig{
	HttpAddr: "http://172.17.0.1:8118",
	Version:  "v1.37",
}

var network_name string = "gojudge"

func TestDockerExecutor(t *testing.T) {
	client, err := NewDockerClient(config.HttpAddr, config.Version)
	if err != nil {
		t.Error(err)
	}
	var args []string
	args = append(args, "./showargs")
	args = append(args, "--submitId=1000")
	resource := Resource{
		Memory: 256 * 1024 * 1024,
	}
	executor := NewDockerExecutor("127.0.0.1", "8118", "showargs", "gojudge", "gojudge1", args, resource, client)
	err = executor.CreateAndStart()
	if err != nil {
		t.Error(err)
	}
	err = executor.Destroy()
	if err != nil {
		t.Error(err)
	}
}
