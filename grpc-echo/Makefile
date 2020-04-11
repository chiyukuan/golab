# makefile for echo service

#export GO111MODULE=on
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: pb
pb:
	protoc -I pb echo.proto --go_out=plugins=grpc:pb

build:
	@go build -o bin/grpc-echo main/*.go

docker-build:
	@docker build -f package/docker/Dockerfile.grpc-echo -t quay.io/raykuan/grpc-echo .

docker-cp:
	@docker create -ti --name dummy quay.io/raykuan/grpc-echo /bin/sh
	@docker cp dummy:/data/go.mod .
	@docker cp dummy:/data/go.sum .
	@docker rm -fv dummy

docker-server-1:
	@docker run --rm -it -p 127.0.0.1:5001:35001/tcp quay.io/raykuan/grpc-echo /app/grpc-echo -mode s -ports 35001

docker-server-2:
	@docker run --rm -it -p 127.0.0.1:5002:35001/tcp quay.io/raykuan/grpc-echo /app/grpc-echo -mode s -ports 35001

server: build
	bin/grpc-echo -mode s

client: build
	bin/grpc-echo -mode c

server-1: build
	bin/grpc-echo -mode s -ports 5001

server-2: build
	bin/grpc-echo -mode s -ports 5002

client-12: build
	bin/grpc-echo -mode c -servers localhost:5001,localhost:5002

# with arguments
runargs := $(word 2, $(MAKECMDGOALS) )

clients: build
	bin/grpc-echo -mode c -servers $(runargs)

servers: build
	bin/grpc-echo -mode s -ports $(runargs)
