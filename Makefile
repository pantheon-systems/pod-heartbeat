REGISTRY := quay.io/getpantheon
APP := pod-heartbeat

# determinse the docker tag to build
ifeq ($(CIRCLE_BUILD_NUM),)
	BUILD_NUM := dev
else
	BUILD_NUM := $(CIRCLE_BUILD_NUM)
	QUAY := docker login -p "$$QUAY_PASSWD" -u "$$QUAY_USER" -e "unused@unused" quay.io
endif

IMAGE := $(REGISTRY)/$(APP):$(BUILD_NUM)

all: test build

build: ## build it
	go build

build-linux: ##  build for linux
	GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w"

build-docker: build-linux ## build the container
	docker build -t $(IMAGE) .

coveralls: deps-coverage ## run coveralls
	gotestcover -v -race  -coverprofile=coverage.out $$(go list ./... | grep -v /vendor/)
	goveralls -repotoken $$COVERALLS_TOKEN -service=circleci -coverprofile=coverage.out

deploy: ## push the image to quay
	make push-circle

deps: _gvt-install
	find  ./vendor/* -maxdepth 0 -type d -exec rm -rf "{}" \;
	gvt rebuild

deps-circle:
	bash scripts/install-go.sh

deps-coverage:
	go get github.com/pierrre/gotestcover
	go get github.com/mattn/goveralls

deps-status: ## check status of deps with gostatus
	go get -u github.com/shurcooL/gostatus
	go list -f '{{join .Deps "\n"}}' . | gostatus -stdin -v

fix_circle_go: # ensure go 1.6 is setup
	scripts/install-go.sh

push: ## push the container the the registry
	docker push $(IMAGE)

push-circle:
	make build-docker
	$(QUAY)
	make push

test: ## TESTS!
	go test -race -v $$(go list ./... | grep -v /vendor/)

update-deps: _gvt-install ## update the deps
	bash scripts/refresh.sh

_gvt-install:
	go get -u github.com/FiloSottile/gvt

help: ## print list of tasks and descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

.PHONY: all test
