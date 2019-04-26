export GO111MODULE=on

all: install test

########################################
### Install & Build

build: go.sum 
ifeq ($(OS),Windows_NT)
	go build -o -mod=readonly -o build/plasmad.exe ./server/plasmad
	go build -o -mod=readonly -o build/plasmacli.exe ./client/plasmacli
else
	go build -o -mod=readonly -o build/plasmad ./server/plasmad
	go build -o -mod=readonly -o build/plasmacli ./client/plasmacli
endif

install: go.sum
	go install -mod=readonly ./server/plasmad
	go install -mod=readonly ./client/plasmacli

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
.PHONY: build build-plasmad build-plasmacli install go.sum
