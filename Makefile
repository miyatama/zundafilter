export
CGO_LDFLAGS := $(shell mecab-config --libs)
CGO_CFLAGS := -I$(shell mecab-config --inc-dir)

# version
VERSION:=$(shell cat VERSION)

# output directory
BINDER:=bin

# data directory
DATADIR:=data

# package info
ROOT_PACKAGE:=$(shell go list .)
COMMAND_PACKAGES:=$(shell go list ./cmd/...)

# output binary
BINARIES:=$(COMMAND_PACKAGES:$(ROOT_PACKAGE)/cmd/%=$(BINDER)/%)

# go file
GO_FILES:=$(shell find . -type f -name '*.go' -print)

# build LDFLAGS
GO_LDFLAGS_VERSION:=-X '${ROOT_PACKAGE}.VERSION=${VERSION}' 
GO_LDFLAGS_SYMBOL:=
ifdef RELEASE
	GO_LDFLAGS_SYMBOL:=-w -s
endif
GO_LDFLAGS_STATIC:=
ifdef RELEASE
	GO_LDFLAGS_STATIC:=-extldflags '-static'
endif
GO_LDFLAGS:=$(GO_LDFLAGS_VERSION) $(GO_LDFLAGS_SYMBOL) $(GO_LDFLAGS_STATIC)

# build tags
GO_BUILD_TAGS:=debug
ifdef RELEASE
	GO_BUILD_TAGS:=release
endif

GO_BUILD_RACE:=-race
ifdef RELEASE
	GO_BUILD_RACE:=
endif

GO_BUILD_STATIC:=
ifdef RELEASE
	GO_BUILD_STATIC:=-a -installsuffix netgo
	GO_BUILD_TAGS:=$(GO_BUILD_TAGS),netgo
endif
GO_BUILD:=-tags=$(GO_BUILD_TAGS) $(GO_BUILD_RACE) $(GO_BUILD_STATIC) -ldflags "$(GO_LDFLAGS)"

.PHONY: build
build: $(BINARIES)

.PHONY: test
test: 
	@go test -v ./zunda_mecab
	@go test -v ./filters

.PHONY: clean
clean:
	@$(RM) -fr $(GOPB_FILES) $(BINARIES) $(BINDER)

$(BINARIES): $(GO_FILES) VERSION 
	@go build -o $@ $(GO_BUILD) $(@:$(BINDER)/%=$(ROOT_PACKAGE)/cmd/%)
	@cp -r $(DATADIR) $(BINDER)/$(DATADIR)
