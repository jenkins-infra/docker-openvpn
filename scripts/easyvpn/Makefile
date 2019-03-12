.PHONY: build

GOLANG_VERSION="1.12"

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
		golang:sandbox-$(GOLANG_VERSION)
		go get && go build -v


build: build.docker build.easyvpn
	ln bin/easyvpn ../../easyvpn
