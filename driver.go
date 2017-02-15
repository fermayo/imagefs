package main

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-plugins-helpers/volume"
)

// ImagefsDriver implements the docker volume driver interface to
// mount images as volumes and commit volume diffs as image layers
type ImagefsDriver struct {
	cli *client.Client
}

// Create creates a volume
func (d ImagefsDriver) Create(r volume.Request) volume.Response {
	fmt.Printf("-> Create %+v\n", r)
	source, ok := r.Options["source"]
	if !ok {
		return volume.Response{Err: "no source volume specified"}
	}

	// pull the image
	readCloser, err := d.cli.ImagePull(context.Background(), source, types.ImagePullOptions{
		// HACK assume the registry ignores the auth header
		RegistryAuth: "null",
	})
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("unexpected error: %s", err)}
	}
	scanner := bufio.NewScanner(readCloser)
	for scanner.Scan() {
	}

	containerConfig := &container.Config{
		Image:      source,
		Entrypoint: []string{"/runtime/loop"},
		Labels: map[string]string{
			"com.docker.imagefs.version": version,
			"com.docker.imagefs.target":  r.Options["target"],
		},
	}
	// TODO handle error
	hostConfig := &container.HostConfig{
		Binds: []string{"/tmp/runtime:/runtime"},
	}
	networkConfig := &network.NetworkingConfig{}
	cont, err := d.cli.ContainerCreate(
		context.Background(),
		containerConfig,
		hostConfig,
		networkConfig,
		// TODO(rabrams) namespace
		r.Name,
	)
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("unexpected error: %s", err)}
	}
	d.cli.ContainerStart(
		context.Background(),
		cont.ID,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("unexpected error: %s", err)}
	}
	return volume.Response{Err: ""}
}

// List lists available volumes
func (d ImagefsDriver) List(r volume.Request) volume.Response {
	containers, err := d.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("unexpected error: %s", err)}
	}
	response := volume.Response{}
	for i := range containers {
		_, ok := containers[i].Labels["com.docker.imagefs.version"]
		if !ok {
			continue
		}
		response.Volumes = append(response.Volumes, &volume.Volume{
			// TODO(rabrams) fall back to id if no names
			Name: containers[i].Names[0],
		})
	}
	return response
}

// Get gets a volume
func (d ImagefsDriver) Get(r volume.Request) volume.Response {
	fmt.Printf("-> Mount %+v\n", r)
	container, err := d.cli.ContainerInspect(context.Background(), r.Name)
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("unexpected error: %s", err)}
	}
	if container.GraphDriver.Name != "overlay" {
		return volume.Response{Err: fmt.Sprintf("unexpected graph driver: %s", container.GraphDriver.Name)}
	}
	mergedDir, ok := container.GraphDriver.Data["MergedDir"]
	if !ok {
		return volume.Response{Err: fmt.Sprintf("missing MergedDir")}
	}
	// HACK directory is relative to host but docker will prepend rootfs of the
	// plugin container
	mergedDir = fmt.Sprintf("../../../../../../../../../../../%s", mergedDir)
	return volume.Response{
		Volume: &volume.Volume{
			Name:       r.Name,
			Mountpoint: mergedDir,
		},
		Mountpoint: mergedDir,
	}
}

// Remove removes a volume
func (d ImagefsDriver) Remove(r volume.Request) volume.Response {
	timeout := 60 * time.Second
	err := d.cli.ContainerStop(context.Background(), r.Name, &timeout)
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("unexpected error: %s", err)}
	}
	container, err := d.cli.ContainerInspect(context.Background(), r.Name)
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("unexpected error: %s", err)}
	}
	target, ok := container.Config.Labels["com.docker.imagefs.target"]
	if ok {
		_, err := d.cli.ContainerCommit(context.Background(), r.Name, types.ContainerCommitOptions{
			Reference: target,
		})
		if err != nil {
			return volume.Response{Err: fmt.Sprintf("unexpected error: %s", err)}
		}
		parts := strings.Split(target, "/")
		if len(parts) == 3 {
			// push the image
			readCloser, err := d.cli.ImagePush(context.Background(), target, types.ImagePushOptions{
				// HACK assume the registry ignores the auth header
				RegistryAuth: "null",
			})
			if err != nil {
				return volume.Response{Err: fmt.Sprintf("unexpected error: %s", err)}
			}
			scanner := bufio.NewScanner(readCloser)
			for scanner.Scan() {
			}
		}
	}
	err = d.cli.ContainerRemove(context.Background(), r.Name, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("unexpected error: %s", err)}
	}

	// HACK remove duplicate container until we can figure out why it is being created
	containers, err := d.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err == nil {
		for i := range containers {
			otherTarget, ok := containers[i].Labels["com.docker.imagefs.target"]
			if ok {
				if target == otherTarget {
					d.cli.ContainerRemove(context.Background(), containers[i].ID, types.ContainerRemoveOptions{
						Force: true,
					})
				}
			}
		}
	}

	return volume.Response{}
}

// Path gets the mounted path of a volume
func (d ImagefsDriver) Path(r volume.Request) volume.Response {
	return d.Get(r)
}

// Mount mounts a volume
func (d ImagefsDriver) Mount(r volume.MountRequest) volume.Response {
	return d.Path(volume.Request{Name: r.Name})
}

// Unmount unmounts a volume
func (d ImagefsDriver) Unmount(r volume.UnmountRequest) volume.Response {
	fmt.Printf("-> Unmount %+v\n", r)
	response := volume.Response{}
	fmt.Printf("<- %+v\n", response)
	return response
}

// Capabilities returns the capabilities of the volume driver
func (d ImagefsDriver) Capabilities(r volume.Request) volume.Response {
	fmt.Printf("-> Capabilities %+v\n", r)
	response := volume.Response{Capabilities: volume.Capability{Scope: "local"}}
	fmt.Printf("<- %+v\n", response)
	return response
}
