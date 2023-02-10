# GoChat

## v0.3 Release

* ui
* ...

### Run

* server

![gochat-server](docs/images/gochat-server.gif)

* client (huoyijie)

![gochat-huoyijie](docs/images/gochat-huoyijie.gif)

* client (jack)

![gochat-jack](docs/images/gochat-jack.gif)

### Features

![gochat-features-uml](docs/images/gochat-features-uml.svg)

### Diagrams

* lib

![gochat-lib-uml](docs/images/gochat-lib-uml.svg)

* server

![gochat-server-uml](docs/images/gochat-server-uml.svg)

* client

![gochat-client-uml](docs/images/gochat-client-uml.svg)

* sequence

![gochat-sequence-uml](docs/images/gochat-sequence-uml.svg)

## Docker

```bash
# work dir
cd server

# build executable
go build -o target/gochat-server

# build docker image
docker build -t gochat-server:latest .

# run docker c
docker run -it -v "$(pwd)"/target:/root/.gochat gochat-server:latest

# open container's shell
docker exec -it af2e58909af8 /bin/bash
```

## v0.4 todo

* tls
* emoji
* send file
* group chat