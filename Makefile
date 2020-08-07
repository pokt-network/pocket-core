#!/usr/bin/make -f

GOVERSION=go1.13.15
CURRENTVERSION=$(go version | awk '${print $3}')
INSTALL_PATH?=${GOPATH}/bin/pocket
BUILD_TAGS?='cleveldb'
export GO111MODULE = on


all: test_short install

## Build
build: go.sum
	@echo "--> Building pocket core 🏗"
	@CGO_ENABLED=1 go build -mod=readonly -tags $(BUILD_TAGS) ./...

build_race: go.sum
	@echo "--> Building pocket core 🏗"
	@CGO_ENABLED=1 go build -mod=readonly -race -tags $(BUILD_TAGS) ./...

install: go.sum
	@echo "--> Building pocket core 🏗"
	@CGO_ENABLED=1 go build -tags $(BUILD_TAGS) -o $(INSTALL_PATH) ./app/cmd/pocket_core/main.go
	@echo "--> Done 🚀🌕"

install_unsafe:
	@echo "--> ️This is building an unverified golang version; use at own risk! ⚠"
	@echo "--> Building pocket core 🏗"
	CGO_ENABLED=1 go build -tags $(BUILD_TAGS) -o $(INSTALL_PATH) ./app/cmd/pocket_core/main.go
	@echo "--> Done 🚀🌕"

go.sum: check
	@echo "--> Ensure dependencies have not been modified 🔍"
	@go mod verify
	@go mod tidy
.PHONY: go.sum 

check: 
ifneq ($(GOVERSION), $(CURRENTVERSION))
	@echo "Go version does not match, please install ${GOVERSION} 💥"
	exit 1
endif

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache 📦"
	@go mod download
.PHONY: go-mod-cache

### Testing
test:
	@echo "--> running all tests 📝"
	@go test -mod=readonly -p 1 ./...

test_short:
	@echo "--> running short tests 📝"
	@go test -mod=readonly -short -p 1 ./...

test_race:
	@echo "--> running tests with race 📝"
	@go test -mod=readonly -race -p 1 ./...

test_short_race:
	@echo "--> running short tests with race 📝"
	@go test -mod=readonly -short -race -p 1 ./...

.PHONY:	test \
test_race
