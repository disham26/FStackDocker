# FStackDocker
Welcome to the FStack docker wrapper.

## Varibles
1. Docker d
2. Docker Client cli
3. initiated bool (to check if the client is activated)
4. NetworkType int
5. ProcessSpaceType int

## Structs
1. Docker
    - Name
    - ContainerId
    - ImageId
	- ListenPortMap
	- Proxy
	- Privileged
	- Network
	- Process
	- VolumeMap
	- VirtualEthDevice
	- CreatedTime
	- Cmdline
2. ImageData
    - Id
	- Name
	- Tag 
	- Size
	- BuildTime

## Interfaces
1. Containers
    - Methods
        1. IsInstalled - returns if the container type is installed on the machine
            - Input - None
            - Output - result bool
        2. GetContainerForProcess - returns the container ID on which the process is running
            - Input - Process ID (int)
            - Output - Container ID (string)
        3. GetContainerForListenPort - returns the container ID which is running on the given port
            - Input - Port (int)
            - Output - Container ID (string)
        4. GetContainerForInterface - returns container ID in which the interface is running
            - Input - VirtualEthDevice (string)
            - Output - Container ID (string)
        5. GetContainerData - returns the container struct for the given container ID
            - Input - Container ID (string)
            - Output - Docker struct
        6. GetHashForPath - returns Sha-256 of an internal path in container.
            - Input - Path (string)
            - Output - Hash (hash []byte)
        7. GetUsernameForUid - returns username for internal UID
            - Input - UID (int)
            - Output - Username (string)
        8. GetImageData - returns information about the image
            - Input - Image ID (string)
            - Output - ImageData struct
## Methods
1. InitPlugin
    - This method initialized instance of Docker and sets initialized bool to true
2. CheckInit
    - If InitPlugin is not called, CheckInit method is called before every function call to check if InitPlugin is called.
3. GetAllDockerContainers
    - This method returns all the docker instances present
4. IsDockerInstalled
    - This method returns a bool stating if Docker is installed on the machines

