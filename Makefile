# Makefile for build golang projects
#

GOPATH = $(shell pwd)

BUILD = go install -ldflags \
		"-X oceanstack/common.Version=0.1 \
		-X oceanstack/common.Buildstamp=`date '+%Y-%m-%dT%H:%M:%S'` \
		-X oceanstack/common.Githash=`git rev-parse HEAD`"

oceanserver:prepare
	cd $(GOPATH)/src/oceanstack/cmd/oceanserver && $(BUILD)
oceancli:
	cd $(GOPATH)/src/oceanstack/cmd/oceancli && $(BUILD)
ocean:
	cd $(GOPATH)/src/oceanstack/cmd/ocean && $(BUILD)
oceanengine:prepare
	cd $(GOPATH)/src/oceanstack/cmd/oceanengine && $(BUILD)

prepare:set_proxy
	protoc --go_out=plugins=grpc:. src/oceanstack/rpc/*.proto
set_proxy:
	git config --global http.proxy http://localhost:8123 && export http_proxy=http://localhost:8123 && cd src && govendor sync

.PHONY: clean
clean:
	rm -rf bin/*
