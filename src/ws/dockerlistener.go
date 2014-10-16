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

func Test(input* ContainerArgs, output* string) error {
	*output += fmt.Sprintf("test")
	return nil
}

func CreateContainer(args* CreateContainerArgs, output* string) error {
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
	*output += fmt.Sprintf(string(b))
	return nil

}

func StartContainer(args* ContainerArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	err := client.StartContainer(args.ID, nil)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container %s started", args.ID)
		rpcOutput.ReplyCode = StartContainerCode
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
	
}

func StopContainer(args* ContainerArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	err := client.StopContainer(args.ID, 3)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	} else {
		rpcOutput.Content += fmt.Sprintf("Stopped container %s", args.ID)
		rpcOutput.ReplyCode = StopContainerCode
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}

func KillContainer(args* ContainerArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	err := client.KillContainer(docker.KillContainerOptions{ID: args.ID})
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container %s was killed", args.ID)
		rpcOutput.ReplyCode = KillContainerCode
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}


func RestartContainer(args* ContainerArgs, output* string) error {
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
	*output += fmt.Sprintf(string(b))
	return nil
}

func RemoveContainer(args* RemoveContainerArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	err := client.RemoveContainer(docker.RemoveContainerOptions{ID: args.ID, RemoveVolumes: args.RemoveVolumes, Force: args.Force})
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container %s was removed", args.ID)
		rpcOutput.ReplyCode = RemoveContainerCode
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}


/*
	Repository  string
	Tag string
	Message string
	Author string
*/
func CommitContainer(args* ContainerCommitArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	image, err := client.CommitContainer(docker.CommitContainerOptions{Container: args.ContainerID, 
						Repository: args.Repository,
						Tag: args.Tag,
						Message: args.Message,
						Author: args.Author})
	rpcOutput := RpcOutput{}
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
		b, _ := json.Marshal(rpcOutput)
		*output += fmt.Sprintf(string(b))
	} else {
		rpcOutput.Content += fmt.Sprintf("Commited container, new Image ID is: %s", image.ID)
		rpcOutput.ReplyCode = CommitContainerCode
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}

func ListContainers(args* DockerListArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	imgs, err := client.ListContainers(docker.ListContainersOptions{All: args.ShowAll})
	rpcOutput := RpcOutput{}
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
		b, _ := json.Marshal(rpcOutput)
		*output += fmt.Sprintf(string(b))
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
		*output += fmt.Sprintf(string(b))
	}
	
	return nil
}

func PullImage(args* ImageArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	err := client.PullImage(docker.PullImageOptions{Repository: args.Repository, Registry: args.Registry}, docker.AuthConfiguration{})
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	} else {
		rpcOutput.Content += fmt.Sprintf("Pulled Image: " + args.Registry+"/"+args.Repository)
		rpcOutput.ReplyCode = PullImageCode
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}

func PushImage(args* ImageArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	fmt.Println("Args repo: %s", args.Repository)
	err := client.PushImage(docker.PushImageOptions{Name: args.Repository, Registry: args.Registry}, docker.AuthConfiguration{})
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	} else {
		rpcOutput.Content += fmt.Sprintf("Pushed Image: " + args.Registry+"/"+args.Repository)
		rpcOutput.ReplyCode = PushImageCode
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}

func RemoveImage(args* RemoveImageArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	err := client.RemoveImage(args.Name)
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	} else {
		rpcOutput.Content += fmt.Sprintf("Removed %s Image", args.Name)
		rpcOutput.ReplyCode = RemoveImageCode
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}



func ListImages(args* DockerListArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
        client, _ := docker.NewClient(endpoint)
        imgs, err := client.ListImages(args.ShowAll)
	rpcOutput := RpcOutput{}
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
		b, _ := json.Marshal(rpcOutput)
		*output += fmt.Sprintf(string(b))
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
		*output += fmt.Sprintf(string(b))
	}
	return nil

}

type Msg struct {
		Action int
		Date int64
}

func handleDockerAction(action int) {

	switch(action) {
		case 1: //List Images
			break;
		case 2: //List containers
			break;
		case 3: //Create container
			break;
		case 4: //Start container
			break;
		case 5: //Stop container
			break;
		case 6: //Remove container
			break;
	}

}

func echoHandler(ws *websocket.Conn) {
	msg := make([]byte, 512)
	n, err := ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg[:n])
	var res Msg
	json.Unmarshal([]byte(msg[:n]), &res)
  
	fmt.Println(res)
	fmt.Println(res.Action)
	fmt.Println(res.Date)

	m, err := ws.Write(msg[:n])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", msg[:m])
}

func main() {
	http.Handle("/docker", websocket.Handler(echoHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}