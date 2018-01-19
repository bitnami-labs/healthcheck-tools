ROOT_PKG=github.com/bitnami-labs/healthcheck-tools
ROOT_PKG_DIR=${GOPATH}/src/$(ROOT_PKG)

SELF_DIR:=$(dir $(lastword $(MAKEFILE_LIST)))

# since go1.8 people can use go without having to define a GOPATH env
# this is the default value the go tooling would assume.
GOPATH?=~/go

EXTGOTOOLS=github.com/golang/protobuf/protoc-gen-go/...

godep-save: symlink
	cd $(ROOT_PKG_DIR) && godep save $$(scripts/gopkgs) $(EXTGOTOOLS)

godep-restore: symlink
	cd $(ROOT_PKG_DIR) && godep restore $$(scripts/gopkgs)
