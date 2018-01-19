all: build

include ../scripts/go.mk

build: symlink
	cd $(ROOT_PKG_DIR)/ssl-checker && go build -i

test: build
	cd $(ROOT_PKG_DIR)/ssl-checker && go test ./...

generate:
	cd $(ROOT_PKG_DIR)/ssl-checker && go generate ./...

clean:
	cd $(ROOT_PKG_DIR)/ssl-checker && go clean ./...

lint:
	golint
	go vet
