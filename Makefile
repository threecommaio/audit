OUT := audit
VERSION ?= $(shell git describe --tags)
PKG=github.com/threecommaio/audit
BUILDTIME:=$(shell date +"%Y.%m.%d.%H%M%S")
COMMIT_HASH:=$(shell git log --pretty=format:'%h' -n 1)

FLAGS:=-X ${PKG}/pkg.Version=${VERSION}
FLAGS:=${FLAGS} -X ${PKG}/pkg.BuildTime=${BUILDTIME}
FLAGS:=${FLAGS} -X ${PKG}/pkg.CommitHash=${COMMIT_HASH}

all: build

build:
	@go build -i -o ${OUT} -ldflags="${FLAGS}"

install: build
	@go install

clean:
	-@rm ${OUT} ${OUT}-v*

.PHONY: build install clean