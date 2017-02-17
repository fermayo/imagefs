FROM alpine:latest

ADD rootfs/build/docker-volume-imagefs /
ADD rootfs/build/loop /

CMD ["/docker-volume-imagefs"]
