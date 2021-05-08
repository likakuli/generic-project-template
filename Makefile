DATE := $(shell date +%Y.%m.%d-%H%M)
LATEST_COMMIT := $(shell git log --pretty=format:'%h' -n 1)
BRANCH := $(shell git branch |grep -v "no branch"| grep \*|cut -d ' ' -f2)
BUILT_ON_IP := $(shell [ $$(uname) = Linux ] && hostname -i || hostname )
RUNTIME_VER := $(shell go version)

BUILT_ON_OS := $(shell uname -a)
ifeq ($(BRANCH),)
BRANCH := master
endif

COMMIT_CNT := $(shell git rev-list HEAD | wc -l | sed 's/ //g' )
BUILD_NUMBER := ${BRANCH}-${COMMIT_CNT}
DIR := github.com/likakuli/generic-project-template/cmd/apiserver/app
COMPILE_LDFLAGS := -X "${DIR}.BuildDate=${DATE}" \
                          -X "${DIR}.LatestCommit=${LATEST_COMMIT}" \
                          -X "${DIR}.BuildNumber=${BUILD_NUMBER}" \
                          -X "${DIR}.BuiltOnIP=${BUILT_ON_IP}" \
                          -X "${DIR}.BuiltOnOs=${BUILT_ON_OS}" \
						  -X "${DIR}.RuntimeVer=${RUNTIME_VER}"

.PHONY: generic-project-template
generic-project-template:
	@go build -ldflags '${COMPILE_LDFLAGS}' -o generic-project-template ./cmd/apiserver/apiserver.go
