export GO111MODULE=on

all: install test

########################################
### Install & Build

build: go.sum 
ifeq ($(OS),Windows_NT)
	go build -o -mod=readonly -o build/plasmad.exe ./cmd/plasmad
	go build -o -mod=readonly -o build/plasmacli.exe ./cmd/plasmacli
else
	go build -o -mod=readonly -o build/plasmad ./cmd/plasmad
	go build -o -mod=readonly -o build/plasmacli ./cmd/plasmacli
endif

install: go.sum
	go install -mod=readonly ./cmd/plasmad
	go install -mod=readonly ./cmd/plasmacli

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
.PHONY: all build install go.sum test test-unit
