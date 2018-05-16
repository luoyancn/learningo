# Makefile for build golang projects
#

GOPATH = $(shell pwd)

BUILD = go install -ldflags \
		"-X oceanstack/common.Version=0.1 \
		-X oceanstack/common.Buildstamp=`date '+%Y-%m-%dT%H:%M:%S'` \
		-X oceanstack/common.Githash=`git rev-parse HEAD`"

oceanserver:
	cd $(GOPATH)/src/oceanstack/cmd/oceanserver && $(BUILD)
oceancli:
	cd $(GOPATH)/src/oceanstack/cmd/oceancli && $(BUILD)
ocean:
	cd $(GOPATH)/src/oceanstack/cmd/ocean && $(BUILD)
oceanengine:
	cd $(GOPATH)/src/oceanstack/cmd/oceanengine && $(BUILD)

.PHONY: clean
clean:
	rm -rf bin/*
