all: clean format golint build

include Makefile.common

TARGET ?= logstash-exporter

GOLINT              := $(FIRST_GOPATH)/bin/golangci-lint
DOCKER_REPO         := wing924
DOCKER_IMAGE_NAME   := logstash-exporter

test:
	@echo ">> running tests"
	GO111MODULE=$(GO111MODULE) $(GO) test -short $(pkgs)

format:
	@echo ">> formatting code"
	GO111MODULE=$(GO111MODULE) $(GO) fmt $(pkgs)

golint: golangci-lint
	@echo ">> linting code"
	GO111MODULE=$(GO111MODULE) $(GOLINT) run

build: promu
	@echo ">> building binaries"
	GO111MODULE=$(GO111MODULE) $(PROMU) build --prefix $(PREFIX)

crossbuild: promu
	@echo ">> cross-building binaries"
	GO111MODULE=$(GO111MODULE) $(PROMU) crossbuild

tarball: promu
	@echo ">> building release tarball"
	GO111MODULE=$(GO111MODULE) $(PROMU) tarball --prefix $(PREFIX) $(BIN_DIR)

tarballs: promu
	@echo ">> building release tarballs"
	GO111MODULE=$(GO111MODULE) $(PROMU) crossbuild tarballs $(BIN_DIR)

clean:
	@echo ">> Cleaning up"
	@find . -type f -name '*~' -exec rm -fv {} \;
	@rm -fv $(TARGET)
	@rm -rfv .build

.PHONY: all clean format golint build test

.PHONY: golangci-lint
golangci-lint: $(GOLINT)
