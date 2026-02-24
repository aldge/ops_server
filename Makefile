PROJECT = ops_server
VERSION = $(shell date +%m%d%H%M)
GITLAB = registry.gitlab.com/cinemae
URL = $(GITLAB)/$(PROJECT)
REGISTRY_RELEASE = $(URL):$(VERSION)
REGISTRY_LATEST = $(URL):latest


net:
	@docker network inspect shared-network >/dev/null 2>&1 || docker network create shared-network

build:
	docker build -t $(PROJECT):latest .

run:
	go run main.go

up:
	docker compose up -d ops_center ops_server --scale mysql=0 --scale redis=0

up-all:
	docker compose up -d

down:
	docker compose  down

swag:
	swag init

push:
	docker tag $(PROJECT) $(REGISTRY_RELEASE)
	docker tag $(PROJECT) $(REGISTRY_LATEST)
	docker push $(REGISTRY_RELEASE)
	docker push $(REGISTRY_LATEST)
	docker rmi $(REGISTRY_RELEASE)
	docker rmi $(REGISTRY_LATEST)