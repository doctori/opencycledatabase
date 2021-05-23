# Bump these on release
VERSION_MAJOR ?= 1
VERSION_MINOR ?= 0
VERSION_BUILD ?= 0

VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)
VERSION_PACKAGE = $(REPOPATH/pkg/version)

SHELL := /bin/bash
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
ORG := github.com/doctori
PROJECT := opencycledatabase
REGISTRY?=doctori/opencycledatabase
REPOPATH ?= $(ORG)/$(PROJECT)
VERSION_PACKAGE = $(REPOPATH)/pkg/version

GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GO_LDFLAGS := '-extldflags "-static"
GO_LDFLAGS += -X $(VERSION_PACKAGE).version=$(VERSION)
GO_LDFLAGS += -w -s # Drop debugging symbols.
GO_LDFLAGS += '

ENTRYPOINT = $(REPOPATH)/cmd/opencycledatabase

dist/ocd : $(GO_FILES)
	GOARCH=$(GOARCH) GOOS=linux CGO_ENABLED=0 go build -ldflags $(GO_LDFLAGS) -o $@ $(ENTRYPOINT)

.PHONY: clean
clean : 
	rm -rf dist/*