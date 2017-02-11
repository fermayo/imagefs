package main

import (
	"fmt"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/docker/docker/client"
)

func main() {
	fmt.Println("Starting ImageFS plugin")
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	d := ImagefsDriver{cli: cli}
	h := volume.NewHandler(d)
	fmt.Println(h.ServeUnix("imagefs", 1))
}
