package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
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
	Network          NetworkType
	Process          ProcessSpaceType
	VolumeMap        map[string]string
	VirtualEthDevice string
	CreatedTime      time.Time
	Cmdline          string
}

func (d Docker) IsInstalled() bool {
	cmd := exec.Command("docker", "version")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
func (d Docker) GetContainerForProcess(pid int) {
	fmt.Println("to do")
}
func (d Docker) GetContainerForListenPort(port int) string {
	return ""
}
func (d Docker) GetContainerForInterface(virtualEthDevice string) string {
	return ""
}
func (d Docker) GetContainerData(containerId string) {

}
func (d Docker) GetHashForPath(path string) []byte {
	return nil
}
func (d Docker) GetUsernameForUid(uid int) string {
	return ""
}
func (d Docker) GetImageData(id string) *ImageData {
	return nil
}

type ImageData struct {
	Id        string
	Name      string
	Tag       string
	Mtime     time.Time
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

func DockerInstalled() string {
	result := "Nothing installed"
	if d.IsInstalled() {
		result = "Docker is installed"
	}
	//will include rest of the containers logic ToDo
	return result
}

var d Docker

func main() {

	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		fmt.Println(err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Size of the container is :", len(containers))
	for _, container := range containers {
		fmt.Println(container)
	}
}
