package FStackContainers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

type NetworkType int
type ProcessSpaceType int

//Constant for docker type
const DOCKER_TYPE = "Docker"

//Constants for network types
const NETWORK_TYPE_BRIDGE = 1
const NETWORK_TYPE_HOST = 2
const NETWORK_TYPE_CONTAINER = 3
const NETWORK_TYPE_NONE = 4
const NETWORK_TYPE_DEFAULT = 5
const NETWORK_TYPE_OTHER = 6

//Instance of Docker
var d Docker

//Instance of client
var cli *client.Client

//Flag to check if instance is initiated
var initiated = false

type Containers interface {

	// Is docker installed on host?
	IsInstalled() bool

	// Get container associated with various objects
	GetContainerForProcess(pid int) (containerId string)

	GetContainerForListenPort(port int) (containerId string)

	GetContainerForInterface(virtualEthDevice string) (containerId string)

	//Get data about a container.
	GetContainerData(containerId string) Containers

	//Get Sha-256 of an internal path in container.
	GetHashForPath(path string) (hash []byte)

	//Get username for internal UID
	GetUsernameForUid(uid int) string

	// Get information about the image
	GetImageData(id string) *ImageData
}

//Docker struct
type Docker struct {
	ContainerType    string
	Name             string
	ContainerId      string
	ImageId          string
	ListenPortMap    map[uint16]uint16
	Proxy            int // Pid of docker-proxy
	Privileged       bool
	Network          NetworkType //ask
	Process          ProcessSpaceType
	VolumeMap        map[string]string
	VirtualEthDevice string
	CreatedTime      time.Time
	Cmdline          strslice.StrSlice
	NetworkId        string
}

//ImageData struct
type ImageData struct {
	Id        string
	Name      string
	Tag       []string
	Size      int64
	BuildTime time.Time
}

//Returns the container ID for the given port
func (d Docker) GetContainerForListenPort(port int) string {
	if CheckInit() {
		containers, _ := GetAllDockerContainers()
		for _, container := range containers {
			result, err := cli.ContainerInspect(context.Background(), container.ContainerId)
			if err != nil {
				fmt.Println(err)
			}
			if result.State.Pid == port {
				return container.ContainerId
			}
		}
	}

	return ""
}

//ToDo
func (d Docker) GetContainerForInterface(virtualEthDevice string) string {
	CheckInit()
	return ""
}

//Returns docker for a given container ID
func (d Docker) GetContainerData(containerId string) Docker {
	var docker Docker
	if CheckInit() {
		containersList, _ := GetAllDockerContainers()

		for _, container := range containersList {
			if container.ContainerId == containerId {
				docker = container
				return docker
			}
		}
	}
	return docker

}

//ToDo
func (d Docker) GetHashForPath(path string) []byte {
	CheckInit()
	return nil
}

//ToDO
func (d Docker) GetUsernameForUid(uid int) string {
	CheckInit()
	return ""
}

//Returns Image details for a given image ID
func (d Docker) GetImageData(id string) *ImageData {
	if CheckInit() {
		out, _, err := cli.ImageInspectWithRaw(context.Background(), id)
		if err != nil {
			panic(err)
		}
		var imageData ImageData
		imageData.Id = out.ID
		imageData.Name = out.GraphDriver.Name
		imageData.Tag = out.RepoTags
		imageData.Size = out.Size
		imageData.BuildTime = out.Metadata.LastTagTime
		return &imageData
	}
	return nil
}

//Enter a container process ID and the function will return the container ID
func (d Docker) GetContainerForProcess(pid int) string {
	if CheckInit() {
		containers, _ := GetAllDockerContainers()
		for _, container := range containers {
			result, err := cli.ContainerInspect(context.Background(), container.ContainerId)
			if err != nil {
				fmt.Println(err)
			}
			if result.State.Pid == pid {
				return container.ContainerId
			}
		}
	}

	return ""
}

//Get the list of all the docker containers
func GetAllDockerContainers() ([]Docker, bool) {
	DockerExists := true
	var dockers []Docker
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		DockerExists = false
		return dockers, DockerExists
	}
	for _, container := range containers {
		containerInspectResult, _ := cli.ContainerInspect(context.Background(), container.ID)
		NetworkMode := containerInspectResult.HostConfig.NetworkMode.NetworkName()
		var networkType NetworkType
		NetworkImage, _ := cli.NetworkInspect(context.Background(), NetworkMode, types.NetworkInspectOptions{Verbose: false, Scope: "global"})
		switch NetworkMode {
		case "bridge":
			networkType = NETWORK_TYPE_BRIDGE
		case "host":
			networkType = NETWORK_TYPE_HOST
		case "container":
			networkType = NETWORK_TYPE_CONTAINER
		case "none":
			networkType = NETWORK_TYPE_NONE
		case "default":
			networkType = NETWORK_TYPE_DEFAULT
		default:
			networkType = NETWORK_TYPE_DEFAULT
		}
		i, _ := strconv.ParseInt(containerInspectResult.Created, 10, 64)
		networkId := container.NetworkSettings.Networks["bridge"].NetworkID
		ports := container.Ports
		var listenPorts map[uint16]uint16
		for _, port := range ports {
			listenPorts[port.PublicPort] = port.PrivatePort
		}
		docker := Docker{
			ContainerType: DOCKER_TYPE,
			Name:          containerInspectResult.Name,
			ContainerId:   container.ID,
			ImageId:       container.ImageID,
			ListenPortMap: listenPorts,
			//Proxy:      "",
			Privileged: containerInspectResult.HostConfig.Privileged,
			Network:    networkType,
			//Process: ,
			//VolumeMap:        "",
			VirtualEthDevice: NetworkImage.Name,
			CreatedTime:      time.Unix(i, 0),
			Cmdline:          containerInspectResult.Config.Cmd,
			NetworkId:        networkId,
		}
		dockers = append(dockers, docker)
	}
	return dockers, DockerExists
}

//Returns boolean for docker installed in machine
func IsDockerInstalled() bool {
	return CheckInit()
}

//Checks if instance is initiated
func CheckInit() bool {
	if !initiated {
		return InitPlugin()
	}
	return initiated
}

//Initiates the plugin
func InitPlugin() bool {
	var err error
	cli, err = client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		return false
	}
	_, exists := GetAllDockerContainers()
	initiated = exists
	return initiated
}
