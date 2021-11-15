.PHONY: build publish client-connect up

IMAGE = jenkinsciinfra/openvpn
TAG = $(shell git rev-parse HEAD | cut -c1-6)
GIT_TREE_STATE = $(shell test -z "`git status --porcelain`" && echo "clean" || echo "dirty")
GIT_SCM_URL = $(shell git config --get remote.origin.url)
SCM_URI = $(subst git@github.com:,https://github.com/,$(GIT_SCM_URL))
BUILD_DATE = $(shell date --utc '+%Y-%m-%dT%H:%M:%S' 2>/dev/null || gdate --utc '+%Y-%m-%dT%H:%M:%S')

build:
	docker build -f ../Dockerfile \
		--no-cache \
		--label "org.opencontainers.image.source=$(GIT_SCM_URL)" \
		--label "org.label-schema.vcs-url=$(GIT_SCM_URL)" \
		--label "org.opencontainers.image.url=$(SCM_URI)" \
		--label "org.label-schema.url=$(SCM_URI)" \
		--label "org.opencontainers.image.revision=$(TAG)" \
		--label "org.label-schema.vcs-ref=$(TAG)" \
		--label "org.opencontainers.image.created=$(BUILD_DATE)" \
		--label "org.label-schema.build-date=$(BUILD_DATE)" \
		--label "io.jenkins-tools.tree.state=$(GIT_TREE_STATE)" \
		-t $(IMAGE):$(TAG) \
		-t $(IMAGE):latest \
		..

publish:
	docker push $(IMAGE):$(TAG)
	docker push $(IMAGE):latest

up:
	docker-compose build
	docker-compose up -d
	docker exec -i -t docker_ldap_1 /entrypoint/restore

client-connect:
	docker-compose exec --workdir=/root/openvpn-client vpn openvpn --config config.ovpn --auth-user-pass auth
