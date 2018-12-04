package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	_ "github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	_ "log"
)

type DockerClient struct {
	*client.Client
	context context.Context
	addr    string
	version string
}

type DockerClientConfig struct {
	HttpAddr string `json:"httpAddr"`
	Version  string `json:"version"`
}

//cli, err := client.NewClient("http://172.17.0.1:8118", "v1.37", nil, nil)

func NewDockerClient(httpAddr string, version string) (*DockerClient, error) {
	cli, err := client.NewClient(httpAddr, version, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DockerClient{
		cli,
		context.Background(),
		httpAddr,
		version,
	}, nil
}
func (dockerClient *DockerClient) String() string {
	return "addr " + dockerClient.addr + " version " + dockerClient.version
}

type Resource struct {
	Memory int64
	cpu    int64
}

const fm int64 = 4 * 1024 * 1024
const om int64 = 1024 * 1024
const ok int64 = 1024
const og int64 = 1024 * 1024 * 1024

func (resc *Resource) SetMemoryByMb(mem int64) {
	if mem < 4 {
		mem = 5
	}
	resc.Memory = mem * om
}

func (resc *Resource) SetMemoryByByte(mem int64) {
	if mem < fm {
		mem = 5 * om
	}
	resc.Memory = mem
}

func (resc *Resource) SetMemoryByGB(mem int64) {
	if mem < 1 {
		mem = 1
	}
	resc.Memory = mem * og
}

func (resc *Resource) SetMemoryByKB(mem int64) {
	if mem/om < 4 {
		mem = 1024 * 5
	}
	resc.Memory = mem * ok
}

const (
	waiting int = 0
	started int = 1
	stoped  int = 2
)

type DockerExecutor struct {
	localAddr string
	localPort string
	container string
	cmd       []string
	image     string
	network   string
	resource  Resource
	client    *DockerClient
	status    int
}

func NewDockerExecutor(localAddr, localPort, image, network, containerName string,
	cmd []string, resource Resource, client *DockerClient) *DockerExecutor {
	return &DockerExecutor{
		localAddr: localAddr,
		localPort: localPort,
		cmd:       cmd,
		image:     image,
		network:   network,
		resource:  resource,
		client:    client,
		container: containerName,
		status:    waiting,
	}
}

type DockerExecutorError struct {
	msg string
}

func (executorError *DockerExecutorError) Error() string {
	return executorError.msg
}

func (executor *DockerExecutor) CreateAndStart() error {
	if executor.status != waiting {
		return &DockerExecutorError{
			msg: "already started",
		}
	}
	nconfig := network.NetworkingConfig{}
	hostconfig := &container.HostConfig{
		Resources: container.Resources{
			Memory: executor.resource.Memory,
		},
	}
	config := container.Config{
		Image:        executor.image,
		Cmd:          executor.cmd,
		OpenStdin:    true,
		AttachStdout: true,
		AttachStdin:  true,
		AttachStderr: true,
	}
	_, err := executor.client.ContainerCreate(executor.client.context, &config, hostconfig, &nconfig, executor.container)
	if err != nil {
		return err
	}
	executor.client.NetworkConnect(executor.client.context, executor.network, executor.container, &network.EndpointSettings{})
	err = executor.client.ContainerStart(executor.client.context, executor.container, types.ContainerStartOptions{})
	if err != nil {
		executor.Destroy()
		return err
	}
	executor.status = started
	return nil
}

func (executor *DockerExecutor) Destroy() error {
	if executor.status != started {
		return &DockerExecutorError{
			msg: "not start",
		}
	}
	err := executor.client.ContainerRemove(executor.client.context, executor.container, types.ContainerRemoveOptions{Force: true})
	executor.status = stoped
	if err != nil {
		return err
	}
	return nil
}
