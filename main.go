package FStackContainers

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

type NetworkType int
type ProcessSpaceType int

type Docker struct {
	Name             string
	ContainerId      string
	ImageId          string
	ListenPortMap    map[int]int
	Proxy            int // Pid of docker-proxy
	Privileged       bool
	Network          NetworkType //ask
	Process          ProcessSpaceType
	VolumeMap        map[string]string
	VirtualEthDevice string
	CreatedTime      time.Time
	Cmdline          strslice.StrSlice
}

//Get the list of all the docker containers
func GetAllDockerContainers() []Docker {
	CheckInit()
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		fmt.Println(err)
	}
	var dockers []Docker
	for _, container := range containers {
		containerInspectResult, _ := cli.ContainerInspect(context.Background(), container.ID)
		i, _ := strconv.ParseInt(containerInspectResult.Created, 10, 64)
		docker := Docker{
			Name:        containerInspectResult.Name,
			ContainerId: container.ID,
			ImageId:     container.ImageID,
			Privileged:  containerInspectResult.HostConfig.Privileged,
			CreatedTime: time.Unix(i, 0),
			Cmdline:     containerInspectResult.Config.Cmd,
		}
		dockers = append(dockers, docker)
	}
	return dockers
}

//Check is docker is installed on this machine
func (d Docker) IsInstalled() bool {
	CheckInit()
	cmd := exec.Command("docker", "version")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (d Docker) GetContainerForListenPort(port int) string {
	CheckInit()
	containers := GetAllDockerContainers()
	for _, container := range containers {
		result, err := cli.ContainerInspect(context.Background(), container.ContainerId)
		if err != nil {
			fmt.Println(err)
		}

		if result.State.Pid == port {
			return container.ContainerId
		}
	}
	return ""
}
func (d Docker) GetContainerForInterface(virtualEthDevice string) string {
	CheckInit()
	return ""
}
func (d Docker) GetContainerData(containerId string) Docker {
	CheckInit()
	dockerJSON, _ := cli.ContainerInspect(context.Background(), containerId)
	i, _ := strconv.ParseInt(dockerJSON.Created, 10, 64)
	docker := Docker{
		Name:        dockerJSON.Name,
		ContainerId: dockerJSON.ID,
		ImageId:     dockerJSON.Image,
		Privileged:  dockerJSON.HostConfig.Privileged,
		CreatedTime: time.Unix(i, 0),
	}
	return docker

}
func (d Docker) GetHashForPath(path string) []byte {
	CheckInit()
	return nil
}
func (d Docker) GetUsernameForUid(uid int) string {
	CheckInit()
	return ""
}
func (d Docker) GetImageData(id string) *ImageData {
	CheckInit()
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

//Enter a container process ID and the function will return the container ID
func (d Docker) GetContainerForProcess(pid int) string {
	CheckInit()
	containers := GetAllDockerContainers()
	for _, container := range containers {
		result, err := cli.ContainerInspect(context.Background(), container.ContainerId)
		if err != nil {
			fmt.Println(err)
		}

		if result.State.Pid == pid {
			return container.ContainerId
		}
	}
	return ""
}

type ImageData struct {
	Id        string
	Name      string
	Tag       []string
	Size      int64
	BuildTime time.Time
}

type Containers interface {

	// Is docker installed on host?
	IsInstalled() bool

	// Get container associated with various objects
	GetContainerForProcess(pid int) (containerId string)

	GetContainerForListenPort(port int) (containerId string)

	GetContainerForInterface(virtualEthDevice string) (containerId string)

	//Get data about a container.
	GetContainerData(containerId string)

	//Get Sha-256 of an internal path in container.
	GetHashForPath(path string) (hash []byte)

	//Get username for internal UID
	GetUsernameForUid(uid int) string

	// Get information about the image
	GetImageData(id string) *ImageData
}

func IsDockerInstalled() bool {
	CheckInit()
	if d.IsInstalled() {
		return true
	}
	//will include rest of the containers logic ToDo
	return false
}

var d Docker
var cli *client.Client
var initiated = false

func CheckInit() {
	if !initiated {
		InitPlugin()
	}
}
func InitPlugin() {
	var err error
	cli, err = client.NewClientWithOpts(client.WithVersion("1.39"))

	if err != nil {
		fmt.Println(err)
	}

	initiated = true
}
