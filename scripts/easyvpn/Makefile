.PHONY: build.project build.docker build.easyvpn build

GOLANG_VERSION=1.12

build.project:
	go get
	go build -v -o bin/easyvpn

build.docker:
	docker build \
		--build-arg VERSION=$(GOLANG_VERSION) \
		--build-arg UID=$(shell id -u) \
		--build-arg GID=$(shell id -g) \
		-t golang:sandbox-$(GOLANG_VERSION) .

build.easyvpn:
	docker run \
		--rm \
		-i -t\
		-v "$$PWD":/easyvpn \
		-w /easyvpn \
		-e GOBIN=/easyvpn/bin \
		-e GOPATH=/easyvpn \
		golang:sandbox-$(GOLANG_VERSION) make build.project


build: build.docker build.easyvpn
