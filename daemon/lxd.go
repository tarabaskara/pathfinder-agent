package daemon

import (
	"github.com/giosakti/pathfinder-agent/model"
	client "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

type apiContainers []api.Container

type LXD struct {
	SocketPath string
}

func (a apiContainers) ToContainerList() *model.ContainerList {
	containerList := make(model.ContainerList, len(a))
	for i, c := range a {
		containerList[i] = model.Container{
			Name: c.Name,
		}
	}
	return &containerList
}

func (l LXD) ListContainers() (*model.ContainerList, error) {
	conn, err := client.ConnectLXDUnix(l.SocketPath, nil)
	if err != nil {
		return nil, err
	}

	res, err := conn.GetContainers()
	if err != nil {
		return nil, err
	}

	containerList := apiContainers(res).ToContainerList()

	return containerList, nil
}

func (l LXD) CreateContainer(name string, image string) (bool, error) {
	conn, err := client.ConnectLXDUnix(l.SocketPath, nil)
	if err != nil {
		return false, err
	}

	// Container creation request
	req := api.ContainersPost{
		Name: name,
		Source: api.ContainerSource{
			Type:     "image",
			Server:   "https://cloud-images.ubuntu.com/releases",
			Protocol: "simplestreams",
			Alias:    image,
		},
	}

	// Get LXD to create the container (background operation)
	op, err := conn.CreateContainer(req)
	if err != nil {
		return false, err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return false, err
	}

	// Get LXD to start the container (background operation)
	reqState := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err = conn.UpdateContainerState(name, reqState, "")
	if err != nil {
		return false, err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return false, err
	}

	return true, nil
}