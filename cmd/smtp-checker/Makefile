TOOL=smtp-checker

all: lint build

include ../../scripts/go.mk

build:
	cd $(ROOT_PKG_DIR)/cmd/$(TOOL) && go build -i

test: build
	cd $(ROOT_PKG_DIR)/cmd/$(TOOL) && go test ./...

generate:
	cd $(ROOT_PKG_DIR)/cmd/$(TOOL) && go generate ./...

clean:
	cd $(ROOT_PKG_DIR)/cmd/$(TOOL) && go clean ./...

lint:
	golint
	go vet
