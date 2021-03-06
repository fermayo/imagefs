BUILD_DIR=rootfs/build
BINARY=$(BUILD_DIR)/docker-volume-imagefs
REPONAME=fermayo/imagefs

test:
	docker run --rm -v $(CURDIR):/go/src/app -w /go/src/app golang:1.7 sh -c "go get -v && go test -v"

binary: $(BINARY)

$(BINARY):
	docker run --rm -v $(CURDIR):/go/src/app -w /go/src/app golang:1.7 sh -c "go get -v && CGO_ENABLED=0 GOOS=linux go build -ldflags '-s' -a -installsuffix cgo -v -o $(BINARY)"

clean:
	rm -fr $(BUILD_DIR)

image: binary
	docker build -f Dockerfile -t $(REPONAME) .

plugin: binary
	docker plugin create $(REPONAME) .

plugin-push:
	docker plugin push $(REPONAME)

image-push:
	docker push $(REPONAME)

deploy:
	docker service create \
		--name imagefs \
		--mode global \
		--mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
		--mount type=bind,source=/var/run/docker/plugins,destination=/run/docker/plugins \
		--mount type=bind,source=/var/lib/docker,destination=/var/lib/docker \
		fermayo/imagefs
