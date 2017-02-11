# fermayo/imagefs

**POC ONLY**

Volume plugin to mount images into containers


## Usage

Run the plugin as a service on a swarm:

    $ docker service create \
        --name imagefs \
        --mode global \
        --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
        --mount type=bind,source=/var/run/docker/plugins,destination=/run/docker/plugins \
        --mount type=bind,source=/var/lib/docker,destination=/var/lib/docker \
        fermayo/imagefs:latest

Now you will see all your local images as volumes:

    $ docker volume ls
    DRIVER                   VOLUME NAME
    imagefs                  2941e8b0086fa6140c59523fb12865461e8c6adeab50732a05a6fe4a886ef836
    imagefs                  5745de9bbf7c7cbabbc7dee2da35ade343f59ff07ce02787a10dfaa6a935efd7
    imagefs                  71fd841bdcb61b4e91d6d4b6df030f9371ac9dbef27e51c1a5994c38c63c7c44
    imagefs                  7afbc2b03b9e6259c8b85457ca94490a1856d13a798ec0040c423543b66a9511
    imagefs                  88e169ea8f46ff0d0df784b1b254a15ecfaf045aee1856dca1ec242fdd231ddd
    imagefs                  ad974e767ec4f06945b1e7ffdfc57bd10e06baf66cdaf5a003e0e6a36924e30b
    imagefs                  d34c8c7c784a39e07cd7978a4f541bde3666ced2b5b7b6c7509b161fc4a6ae64
    imagefs                  f49eec89601e8484026a8ed97be00f14db75339925fad17b440976cffcbfb88a
    
You can mount any image on your container:

    $ docker run -it -v 88e169ea8f46ff0d0df784b1b254a15ecfaf045aee1856dca1ec242fdd231ddd:/tmp:ro ubuntu cat /tmp/etc/os-release
    NAME="Alpine Linux"
    ID=alpine
    VERSION_ID=3.5.0
    PRETTY_NAME="Alpine Linux v3.5"
    HOME_URL="http://alpinelinux.org"
    BUG_REPORT_URL="http://bugs.alpinelinux.org"


## Limitations

It does not use any union FS driver, so:
* Only mounts one layer at a time (does not mount parent layers)
* Allows write directly to the image, corrupting it (does not have a write layer on top)
