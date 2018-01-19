ROOT_PKG=github.com/bitnami-labs/healthcheck-tool
TOOLS := $(shell ls ./cmd)

all:
	make $(TOOLS)


$(TOOLS):
	make -C ./cmd/$@

build:
	@$(MAKE) -s $(addprefix build-, $(TOOLS))

build-%:
	make -C cmd/$(*F) build

test:
	@$(MAKE) -s $(addprefix test-, $(TOOLS))

test-%:
	make -C cmd/$(*F) test

lint:
	@$(MAKE) -s $(addprefix lint-, $(TOOLS))

lint-%:
	make -C cmd/$(*F) test

clean:
	@$(MAKE) -s $(addprefix clean-, $(TOOLS))

clean-%:
	make -C cmd/$(*F) clean

godep-save:
	@$(MAKE) -s $(addprefix godep-save-, $(TOOLS))

godep-save-%:
	make -C cmd/$(*F) godep-save

godep-restore:
	@$(MAKE) -s $(addprefix godep-restore-, $(TOOLS))

godep-restore-%:
	make -C cmd/$(*F) godep-restore

get-build-deps:
	@echo "+ Downloading build dependencies"
	@go get github.com/tools/godep
	@go get github.com/golang/lint/golint
