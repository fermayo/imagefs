package main

import (
	"context"
	"github.com/docker/go-plugins-helpers/volume"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"strings"
)

type ImagefsDriver struct {
	cli *client.Client
}

func (d ImagefsDriver) Create(r volume.Request) volume.Response {
	fmt.Printf("-> Create %+v\n", r)
	response := volume.Response{Err: "Not implemented"}
	fmt.Printf("<- %+v\n", response)
	return response
}

func (d ImagefsDriver) List(r volume.Request) volume.Response {
	fmt.Printf("-> List %+v\n", r)
	images, err := d.cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		panic(err)
	}

	volumes := []*volume.Volume{}
	for _, image := range images {
		fmt.Printf("%s\n", image.ID)
		volumes = append(volumes, &volume.Volume{Name: image.ID[7:]})
	}

	response := volume.Response{Volumes:volumes}
	fmt.Printf("<- %+v\n", response)
	return response
}

func (d ImagefsDriver) Get(r volume.Request) volume.Response {
	fmt.Printf("-> Get %+v\n", r)
	inspect, _, err := d.cli.ImageInspectWithRaw(context.Background(), r.Name)
	if err != nil {
		return volume.Response{Err: fmt.Sprint(err)}
	}

	response := volume.Response{}
	if inspect.GraphDriver.Name == "overlay2" {
		response = volume.Response{Volume: &volume.Volume{Name: r.Name,
			Mountpoint: strings.Split(inspect.GraphDriver.Data["UpperDir"], ":")[0]}}
	} else {
		response = volume.Response{Err: fmt.Sprintf("GraphDriver %s not supported", inspect.GraphDriver.Name)}
	}
	fmt.Printf("<- %+v\n", response)
	return response
}

func (d ImagefsDriver) Remove(r volume.Request) volume.Response {
	fmt.Printf("-> Remove %+v\n", r)
	response := volume.Response{Err: "Not implemented"}
	fmt.Printf("<- %+v\n", response)
	return response
}

func (d ImagefsDriver) Path(r volume.Request) volume.Response {
	fmt.Printf("-> Path %+v\n", r)
	response := volume.Response{Mountpoint: d.Get(r).Volume.Mountpoint}
	fmt.Printf("<- %+v\n", response)
	return response
}

func (d ImagefsDriver) Mount(r volume.MountRequest) volume.Response {
	fmt.Printf("-> Mount %+v\n", r)
	response := volume.Response{Mountpoint: d.Get(volume.Request{Name: r.Name}).Volume.Mountpoint}
	fmt.Printf("<- %+v\n", response)
	return response
}

func (d ImagefsDriver) Unmount(r volume.UnmountRequest) volume.Response {
	fmt.Printf("-> Unmount %+v\n", r)
	response := volume.Response{}
	fmt.Printf("<- %+v\n", response)
	return response
}

func (d ImagefsDriver) Capabilities(r volume.Request) volume.Response {
	fmt.Printf("-> Capabilities %+v\n", r)
	response := volume.Response{Capabilities: volume.Capability{Scope: "local"}}
	fmt.Printf("<- %+v\n", response)
	return response
}
