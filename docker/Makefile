.PHONY: build run shell

IMAGE = 'olblak/openvpn'
TAG = $(shell git rev-parse HEAD | cut -c1-6)

build:
	docker build -f ../Dockerfile --no-cache -t $(IMAGE):$(TAG) ..

publish:
	docker push $(IMAGE):$(TAG)

up: 
	docker-compose build
	docker-compose up -d
	docker exec -i -t openvpn_ldap_1 /entrypoint/restore
