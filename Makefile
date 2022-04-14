# go
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOPATH ?= $(shell go env GOPATH)
GO_MAJOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
MIN_GO_MAJOR_VERSION = 1
MIN_GO_MINOR_VERSION = 18
GO_BINARY := $(shell which go)

# docker
DOCKER_BINARY := $(shell which docker)
DOCKER_COMPOSE_BINARY := $(shell which docker-compose)

# sys
UNAME := $(shell uname)
ARCH ?= $(shell go env GOARCH)

# make release option=major [prere=alpha.1 build=build.1]
# make release option=minor [prere=alpha.1 build=build.1]
# make release option=patch [prere=alpha.1 build=build.1]
# make release option=prere prere=alpha.2 [build=build.1]
# make release option=prere prere=alpha.3 build=build.1
.PHONY: release
release: health
	./hack/release.sh $(option) $(prere) $(build)


.PHONY: health
health: fmt test vet lint


.PHONY: test
test:
	go test ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golint ./...

.PHONY: fmt
fmt:
	gofmt -s -w .

