.PHONY: build_osx build_windows build_linux init_osx init_linux init_windows

build_osx:
	go get -d
	GO111MODULE=on CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -v .

build_windows:
	go get -d
	GO111MODULE=on CGO_ENABLED=0 GOOS=windows go build -v .

build_linux:
	go get -d
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -v .

init_osx:
	docker run --rm -v "$(CURDIR)":/easyvpn -w /easyvpn -e GOOS=windows golang:1.13 make build_osx

init_linux:
	docker run --rm -v "$(CURDIR)":/easyvpn -w /easyvpn -e GOOS=windows golang:1.13 make build_linux

init_windows:
	docker run --rm -v "$(CURDIR)":/easyvpn -w /easyvpn -e GOOS=windows golang:1.13 make build_windows

