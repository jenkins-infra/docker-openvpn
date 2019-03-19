.PHONY: build publish run

build_osx:
	go get -d
	CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -v .

build_windows:
	go get -d
	CGO_ENABLED=0 GOOS=windows go build -v .

build_linux:
	go get -d
	CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -v .

init_osx:
	docker run --rm -v "$$PWD":/easyvpn -w /easyvpn -e GOOS=windows golang:1.12 make build_osx
	if [ -f ../../easyvpn ]; then rm ../../easyvpn ;fi

init_linux:
	docker run --rm -v "$$PWD":/easyvpn -w /easyvpn -e GOOS=windows golang:1.12 make build_linux
	if [ -f ../../easyvpn ]; then rm ../../easyvpn fi

init_windows:
	docker run --rm -v "$$PWD":/easyvpn -w /easyvpn -e GOOS=windows golang:1.12 make build_windows

