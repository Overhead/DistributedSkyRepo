package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"code.google.com/p/go.net/websocket"
	"github.com/fsouza/go-dockerclient"
)

type ContainerArgs struct {
	ID string
}

type CreateContainerArgs struct {
	ContainerName string
	ImageName string
}

type RemoveContainerArgs struct {
    // The ID of the container.
    ID  string

    // A flag that indicates whether Docker should remove the volumes
    // associated to the container.
    RemoveVolumes bool

    // A flag that indicates whether Docker should remove the container
    // even if it is currently running.
    Force bool
}

//Replycodes
const (
    ErrorCode		= 0
    CreateContainerCode 	= 1
    StartContainerCode      = 2
    StopContainerCode     	= 3
    KillContainerCode	= 4
    RestartContainerCode    = 5
    RemoveContainerCode	= 6
    ListContainersCode      = 7
    PullImageCode		= 8
    RemoveImageCode		= 9
    ListImagesCode		= 10
    CommitContainerCode 	= 11
    PushImageCode		= 12
)

type DockerListArgs struct {
	ShowAll bool
}

type ImageArgs struct {
	ID string
	Registry string
	Repository string
}

type RemoveImageArgs struct {
	Name string
}

type ContainerCommitArgs struct {
	ContainerID string
	Repository  string
	Tag string
	Message string
	Author string
}
//Reply codes are used to define what method that it returns from
// 
// 
type RpcOutput struct {
	Content string
	ReplyCode int
}

type RpcOutputCreateCont struct {
	Content string
	ID string
	ReplyCode int
}

type RpcOutputCreateContainer struct {
	Content string
	ID string
	ReplyCode int
}

type Container struct {
	ID string
	Image string
	Created string
	Status string
}

type ContainerCollection struct {
	ReplyCode int
	Containers []Container
}

type InspectedContainer struct {
	ReplyCode int
	Container *docker.Container
}

type Image struct {
	ID string
	Created string
	Size string
	VirtualSize string
	RepoTags []string
}

type ImageCollection struct {
	ReplyCode int
	Images []Image
}

func Test(input* ContainerArgs) string {
	return fmt.Sprintf("test")
}

func CreateContainer(args* CreateContainerArgs) string {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	config := docker.Config{Image: args.ImageName}
	createArgs := docker.CreateContainerOptions{Name: args.ContainerName, Config: &config}	
	container, err := client.CreateContainer(createArgs)
	rpcOutput := RpcOutputCreateContainer{}
	rpcOutput.Content = ""
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container created successfully with new ID: %s", container.ID)
		rpcOutput.ID = container.ID
		rpcOutput.ReplyCode = CreateContainerCode
	}
	b, _ := json.Marshal(rpcOutput)
	return fmt.Sprintf(string(b))

}

func StartContainer(id string) string {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	err := client.StartContainer(id, nil)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container %s started", id)
		rpcOutput.ReplyCode = StartContainerCode
	}
	b, _ := json.Marshal(rpcOutput)
	return fmt.Sprintf(string(b))
	
}

func StopContainer(id string) string {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	err := client.StopContainer(id, 3)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	} else {
		rpcOutput.Content += fmt.Sprintf("Stopped container %s", id)
		rpcOutput.ReplyCode = StopContainerCode
	}
	b, _ := json.Marshal(rpcOutput)
	return fmt.Sprintf(string(b))
}

func KillContainer(id string) string {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	err := client.KillContainer(docker.KillContainerOptions{ID: id})
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container %s was killed", id)
		rpcOutput.ReplyCode = KillContainerCode
	}
	b, _ := json.Marshal(rpcOutput)
	return fmt.Sprintf(string(b))
}


func RestartContainer(args* ContainerArgs) string {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	err := client.RestartContainer(args.ID, 500)
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container %s is restarting", args.ID)
		rpcOutput.ReplyCode = RestartContainerCode
	}
	b, _ := json.Marshal(rpcOutput)
	return fmt.Sprintf(string(b))
}

func RemoveContainer(id string) string {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	err := client.RemoveContainer(docker.RemoveContainerOptions{ID: id, RemoveVolumes: true, Force: true})
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container %s was removed", id)
		rpcOutput.ReplyCode = RemoveContainerCode
	}
	b, _ := json.Marshal(rpcOutput)
	return fmt.Sprintf(string(b))
}


/*
	Repository  string
	Tag string
	Message string
	Author string
*/

func ListContainers(args* DockerListArgs) string {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	imgs, err := client.ListContainers(docker.ListContainersOptions{All: args.ShowAll})
	rpcOutput := RpcOutput{}
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
		b, _ := json.Marshal(rpcOutput)
		return fmt.Sprintf(string(b))
	} else {
		
		list := []Container{}
		for _, img := range imgs {
			cont := Container{}
			cont.ID += fmt.Sprintf(img.ID)
			cont.Image += fmt.Sprintf(img.Image)
			cont.Created += fmt.Sprintf("%v",img.Created)
			cont.Status += fmt.Sprintf(img.Status)
			list = append(list,cont)
		}
		containerColl := ContainerCollection{ReplyCode: ListContainersCode, Containers: list}
		b, _ := json.Marshal(containerColl)
		return fmt.Sprintf(string(b))
	}
	
	return ""
}

func ListImages(args* DockerListArgs) string {
	endpoint := "unix:///var/run/docker.sock"
        client, _ := docker.NewClient(endpoint)
        imgs, err := client.ListImages(args.ShowAll)
	rpcOutput := RpcOutput{}
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
		b, _ := json.Marshal(rpcOutput)
		return fmt.Sprintf(string(b))
	} else {
		list := []Image{}
		for _, img := range imgs {
			image := Image{}
		        image.ID += fmt.Sprintf(img.ID)
		        image.RepoTags = img.RepoTags
		        image.Created += fmt.Sprintf("%v", img.Created)
		        image.Size += fmt.Sprintf("%d",img.Size)
		        image.VirtualSize += fmt.Sprintf("%d", img.VirtualSize)
			list = append(list, image)
        	}
		imageColl := ImageCollection{ReplyCode: ListImagesCode, Images: list}
		b, _ := json.Marshal(imageColl)
		return fmt.Sprintf(string(b))
	}
	return ""

}

func InspectContainer(id string) (*docker.Container, error) {
	endpoint := "unix:///var/run/docker.sock"
  client, _ := docker.NewClient(endpoint)
  container, err := client.InspectContainer(id)
	if err != nil {
			return nil, err
	} else {
		return container, nil
	}
}

type Msg struct {
		Action int
		Container_ID string
		ContainerName string
		ImageID string
		Date int64
}


/*  
    CreateContainerCode 	= 1
    StartContainerCode      = 2
    StopContainerCode     	= 3
    KillContainerCode	= 4
    RestartContainerCode    = 5
    RemoveContainerCode	= 6
    ListContainersCode      = 7
    PullImageCode		= 8
    RemoveImageCode		= 9
    ListImagesCode		= 10
    CommitContainerCode 	= 11
    PushImageCode		= 12
*/
func handleDockerAction(msg* Msg) string {
	var response string
	switch(msg.Action) {
		case 1: //Create container
			var args = CreateContainerArgs{ContainerName: msg.ContainerName, ImageName: msg.ImageID}
			response = CreateContainer(&args)
			break;
		case 2: //Start container
			response = StartContainer(msg.Container_ID)
			break;
		case 3: //Stop container
			response = StopContainer(msg.Container_ID)
			break;
		case 4: //Kill container
			response = KillContainer(msg.Container_ID)
			break;
		case 5: //Restart container
			break;
		case 6: //Remove container
			response = RemoveContainer(msg.Container_ID)
			break;
		case 7: //List containers
			fmt.Println("Get container list\n")
			var args = DockerListArgs{ShowAll: true}
			response = ListContainers(&args)
			break;
		case 10: //List images
			fmt.Println("Get image list\n")
			var args = DockerListArgs{ShowAll: true}
			response = ListImages(&args)
			break;
		case 11: //Tell node to join 
			container, err := InspectContainer(msg.Container_ID)
			if err != nil {
				fmt.Println(err)			
			} else {
				fmt.Println(container)
				data := InspectedContainer{ReplyCode: 11, Container: container}	
				b, _ := json.Marshal(data)
				return fmt.Sprintf(string(b))	
			}
			break;
	}
	return response
}

func echoHandler(ws *websocket.Conn) {
	var err error
	var msg Msg	
	for {
		dec := json.NewDecoder(ws)
		err = dec.Decode(&msg)
		if err != nil {
			log.Fatal(err)
			break
		}
		fmt.Printf("Receive: %s\n", &msg)
		
		fmt.Println(&msg)
	
		response := handleDockerAction(&msg)

		_, err := ws.Write([]byte(response))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Send: %s\n", response)
	}
}

func main() {
	http.Handle("/docker", websocket.Handler(echoHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
