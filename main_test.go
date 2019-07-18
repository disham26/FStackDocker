package FStackContainers

import (
	"testing"
)

func TestInit(t *testing.T) {
	InitPlugin()
	if !initiated {
		t.Error("Initiation failed")
	}
}

func TestInitNegative(t *testing.T) {
	InitPlugin()
	if initiated {
		t.Error("Initiation succes, passed the negative test")
	}
}

func TestCheckInit(t *testing.T) {
	initiated = false
	CheckInit()
	if !initiated {
		t.Error("CheckInit did not work")
	}
}

func TestCheckInitNegative(t *testing.T) {
	initiated = false
	CheckInit()
	if initiated {
		t.Error("CheckInit worked,, passed the negative test")
	}
}

func TestDockerInstall(t *testing.T) {
	if !IsDockerInstalled() {
		t.Error("Docker is installed, still gave false")
	}
}

func TestDockerInstallNegative(t *testing.T) {
	if IsDockerInstalled() {
		t.Error("Docker is installed, passed the negative test")
	}
}

func TestSizeOfDockerContainers(t *testing.T) {
	containers := GetAllDockerContainers()
	if len(containers) != 1 {
		t.Error("I have only 1 container running, got ", len(containers))
	}
}

func TestSizeOfDockerContainersNegative(t *testing.T) {
	containers := GetAllDockerContainers()
	if len(containers) == 1 {
		t.Error("I have only 1 container running, got ", len(containers), ". Passed negative test")
	}
}

func TestGetContainerFromPort(t *testing.T) {
	containerId := d.GetContainerForProcess(2783)
	if containerId != "9c0d34b30cc61587fd310e9bfe8cf341b1bc74111ec6f60ea5e99eda3a543803" {
		t.Error("Expected 9c0d34b30cc61587fd310e9bfe8cf341b1bc74111ec6f60ea5e99eda3a543803, got ", containerId)
	}

}

func TestGetContainerFromPortNegative(t *testing.T) {
	containerId := d.GetContainerForProcess(2783)
	if containerId == "9c0d34b30cc61587fd310e9bfe8cf341b1bc74111ec6f60ea5e99eda3a543803" {
		t.Error("Expected 9c0d34b30cc61587fd310e9bfe8cf341b1bc74111ec6f60ea5e99eda3a543803, got ", containerId, ". Negative test passed")
	}

}
