package executor

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	_ "github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"log"
	"strconv"
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

type DockerContainerConfig struct {
	NetWork        string `json:"network"`
	ProcessAddrMap string `json:"processAddrMap"`
	ProcessPort    int    `json:"processPort"`
	ListenAddrMap  string `json:"listenAddrMap"`
	Volume         string `json:"volume"`
	Image          string `json:"image"`
	Cmd            string `json:"cmd"`
	MountPoint     string `json:"mountPoint"`
	LabelName      string `json:"labelName"`
	LabelValue     string `json:"labelValue"`
}

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

func (dockerClient *DockerClient) RemoveByLabel(labelName, labelValue string) {
	args := filters.NewArgs()
	args.Add("label", labelName+"="+labelValue)
	_, err := dockerClient.ContainersPrune(dockerClient.context, args)
	if err != nil {
		log.Println(err)
	}
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
	localAddr  string
	localPort  int
	container  string
	cmd        []string
	image      string
	network    string
	resource   Resource
	client     *DockerClient
	status     int
	volume     string
	mountPoint string
	labelName  string
	labelValue string
}

func NewDockerExecutor(submitID int, containerName string, resource Resource, client *DockerClient, config DockerContainerConfig) *DockerExecutor {

	cmd := make([]string, 0)
	cmd = append(cmd, config.Cmd)
	cmd = append(cmd, "--adress="+config.ProcessAddrMap)
	cmd = append(cmd, "--port="+strconv.Itoa(config.ProcessPort))
	cmd = append(cmd, "--submitID="+strconv.Itoa(submitID))

	return &DockerExecutor{
		localAddr:  config.ListenAddrMap,
		localPort:  config.ProcessPort,
		cmd:        cmd,
		image:      config.Image,
		network:    config.NetWork,
		resource:   resource,
		client:     client,
		container:  containerName,
		status:     waiting,
		volume:     config.Volume,
		mountPoint: config.MountPoint,
		labelName:  config.LabelName,
		labelValue: config.LabelValue,
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
	var mountVolume mount.Mount = mount.Mount{
		Type:        "volume",
		Source:      executor.volume,
		Target:      executor.mountPoint,
		ReadOnly:    true,
		Consistency: "default",
	}
	nconfig := network.NetworkingConfig{}

	hostconfig := &container.HostConfig{
		Resources: container.Resources{
			Memory: executor.resource.Memory,
		},
		Mounts: []mount.Mount{mountVolume},
	}

	labels := make(map[string]string)
	labels[executor.labelName] = executor.labelValue
	config := container.Config{
		Image:        executor.image,
		Cmd:          executor.cmd,
		OpenStdin:    true,
		AttachStdout: false,
		AttachStdin:  false,
		AttachStderr: false,
		Labels:       labels,
	}

	_, err := executor.client.ContainerCreate(executor.client.context, &config, hostconfig, &nconfig, executor.container)
	if err != nil {
		err = executor.Destroy()
		//		log.Println(err)
		if err != nil {
			return err
		}
		_, err := executor.client.ContainerCreate(executor.client.context, &config, hostconfig, &nconfig, executor.container)
		//		log.Println("create ", err)
		err = executor.client.ContainerStart(executor.client.context, executor.container, types.ContainerStartOptions{})
		//		log.Println("start", err)
		if err != nil {
			return err
		}
	}

	executor.client.NetworkConnect(executor.client.context, executor.network, executor.container, &network.EndpointSettings{})

	executor.client.NetworkDisconnect(executor.client.context, "bridge", executor.container, true)

	err = executor.client.ContainerStart(executor.client.context, executor.container, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	executor.status = started
	return err
}

func (executor *DockerExecutor) Destroy() error {
	err := executor.client.ContainerRemove(executor.client.context, executor.container, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		return err
	}
	executor.status = stoped
	return nil
}
