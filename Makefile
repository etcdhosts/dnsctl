BUILD_VERSION   	:= $(shell cat version)
BUILD_DATE      	:= $(shell date "+%F %T")
COMMIT_SHA1     	:= $(shell git rev-parse HEAD)

all: clean
	bash .cross_compile.sh

release: all
	gh release create ${BUILD_VERSION} -t "Bump ${BUILD_VERSION}" ./dist/*

install:
	go install -trimpath -ldflags	"-X 'main.version=${BUILD_VERSION}' \
               						-X 'main.buildDate=${BUILD_DATE}' \
               						-X 'main.commitID=${COMMIT_SHA1}'"

clean:
	rm -rf dist

.PHONY: all release clean install
