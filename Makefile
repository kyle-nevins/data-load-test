# Makefile for data-load-test

GO ?= go
DEP ?= dep

GOPATH := $(CURDIR)/vendor:$(GOPATH)

APPDIR = $(CURDIR)/app

BINDIR = $(CURDIR)/bin
BINNAME = data-load-test
CONFDIR = $(CURDIR)/conf

SHELL = bash

.PHONY: all check cf-broker cf-push

all: build check cf-broker cf-push

build:
	@mkdir -p $(BINDIR)
	@$(DEP) ensure
	@$(GO) build -o $(BINDIR)/$(BINNAME) $(CURDIR)/server.go

check:
	@which $(GO) >/dev/null
	@which $(DEP) >/dev/null

cf-broker: all
	@cp $(BINDIR)/$(BINNAME) $(APPDIR)
	@cp $(CONFDIR)/{manifest.yml,Procfile} $(APPDIR)
