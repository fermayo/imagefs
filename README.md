# fermayo/imagefs

**POC ONLY**

Volume plugin to mount images into containers


## Usage

Install the plugin

```
docker plugin install fermayo/imagefs
```

Create a task with a new volume

```
docker service create --name ping --replicas 1 --restart-condition none --mount type=volume,volume-driver=fermayo/imagefs,volume-opt=source=alpine,target=/context,volume-opt=target=new alpine ping www.google.com
```

* `volume-opt=source` is the image to mount
* `target` is the directory in the service where the image will be mounted
* `volume-opt=target` is the tag to give to the new volume snapshot
  Note that if the tag contains a registry URL the image will be pushed
  when the volume is unmounted (Not yet implemented)

Modify the volume from within the service

```
tid=ping.1.$(docker service ps -q ping)
docker exec -it $tid touch /context/hello
```

Remove the service

```
docker service rm ping
```

and see the contents of your newly created image

```
docker run -it --rm --entrypoint ls new
```

## Limitations

* Does not pull layers that are not present locally in the node
* Only supports overlay2 graphdriver
