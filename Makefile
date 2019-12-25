export GO111MODULE=on

#######################
### Current Build Properties
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
BUILD_TAGS = netgo

BUILD_FLAGS = -tags "$(BUILD_TAGS)" -ldflags \
    '-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
    -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
    -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(BUILD_TAGS)"'

all: install test

########################################
### Install & Build

build: go.sum build-plasmad build-plasmacli

build-plasmad: go.sum
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o -mod=readonly -o build/plasmad.exe ./cmd/plasmad
else
	go build $(BUILD_FLAGS) -o -mod=readonly -o build/plasmad ./cmd/plasmad
endif

build-plasmacli: go.sum
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o -mod=readonly -o build/plasmacli.exe ./cmd/plasmacli
else
	go build $(BUILD_FLAGS) -o -mod=readonly -o build/plasmacli ./cmd/plasmacli
endif

install: go.sum install-plasmad install-plasmacli

install-plasmad: go.sum
	go install $(BUILD_FLAGS) -mod=readonly ./cmd/plasmad

install-plasmacli: go.sum
	go install $(BUILD_FLAGS) -mod=readonly ./cmd/plasmacli

########################################
### Dependencies & Maintenance

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

clean: 
	rm -rf build/ coverage.txt

########################################
###

test: test-unit

test-unit: 
	go test -mod=readonly -race -coverprofile=coverage.txt -covermode=atomic -v ./...

# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: all build build-plasmad build-plasmacli install install-plasmad install-plasmacli go.sum test test-unit
