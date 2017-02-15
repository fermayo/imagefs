package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/docker/client"
	"github.com/docker/go-plugins-helpers/volume"
)

const (
	// working directory for volume containers -- lets us mount an infinitely looping binary
	// into any container
	runtimeDir = "/tmp/runtime/"
	loopBinary = "loop"
	version    = "0.1"
)

func main() {
	fmt.Println("Starting ImageFS plugin")
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	err = os.Mkdir(runtimeDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	loop, err := ioutil.ReadFile("/loop")
	if err != nil {
		panic(err)
	}
	// TODO handle error
	ioutil.WriteFile("/tmp/runtime/loop", loop, os.ModePerm)

	d := ImagefsDriver{cli: cli}
	h := volume.NewHandler(d)
	fmt.Println(h.ServeUnix("imagefs", 1))
}
