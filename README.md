# Installing

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone git clone https://github.com/nakiner/faceit-subscriber.git

#move to project
$ cd faceit-subscriber

# Build the docker image first (or skip, to use docker.io prebuilt one)
$ make docker

# Run subscriber
$ cd ../ && cd faceit && docker-compose up subscriber -d 

# check if the containers are running
$ docker ps

# Run tests (optional)
$ make test
```