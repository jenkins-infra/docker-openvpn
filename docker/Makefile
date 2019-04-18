.PHONY: build run shell

IMAGE = 'jenkinsciinfra/openvpn'
TAG = $(shell git rev-parse HEAD | cut -c1-6)

build:
	docker build -f ../Dockerfile --no-cache -t $(IMAGE):$(TAG) -t $(IMAGE):latest ..

publish:
	docker push $(IMAGE):$(TAG)
	docker push $(IMAGE):latest

up:
	docker-compose build
	docker-compose up -d
	docker exec -i -t docker_ldap_1 /entrypoint/restore
