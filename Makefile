all: style lint

include Makefile.common

TARGET ?= logstash-exporter

DOCKER_REPO         := wing924
DOCKER_IMAGE_NAME   := logstash-exporter

crossbuild: promu
	@echo ">> cross-building binaries"
	GO111MODULE=$(GO111MODULE) $(PROMU) crossbuild

tarball: build common-tarball

tarballs: crossbuild
	@echo ">> building release tarballs"
	GO111MODULE=$(GO111MODULE) $(PROMU) crossbuild tarballs

clean:
	@echo ">> Cleaning up"
	@find . -type f -name '*~' -exec rm -fv {} \;
	@rm -fv $(TARGET) $(TARGET)-*.tar.gz
	@rm -rfv .build .tarballs
