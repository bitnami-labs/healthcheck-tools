TOOL=ssl-checker

all: build

include ../scripts/go.mk

build:
	cd $(ROOT_PKG_DIR)/$(TOOL) && go build -i

test: build
	cd $(ROOT_PKG_DIR)/$(TOOL) && go test ./...

generate:
	cd $(ROOT_PKG_DIR)/$(TOOL) && go generate ./...

clean:
	cd $(ROOT_PKG_DIR)/$(TOOL) && go clean ./...

lint:
	golint
	go vet
