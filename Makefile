 #
 # Copyright (c) 2018, Juniper Networks, Inc.
 # All rights reserved.
 #

PREFIX := /usr/local
VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git rev-parse --short HEAD)
ifdef GOBIN
PATH := $(GOBIN):$(PATH)
else
PATH := $(subst :,/bin:,$(GOPATH))/bin:$(PATH)
endif

WEDGE := wedge$(shell go tool dist env | grep -q 'GOOS=.windows.' && echo .exe)

LDFLAGS := $(LDFLAGS) -X main.commit=$(COMMIT) -X main.branch=$(BRANCH)
ifdef VERSION
	LDFLAGS += -X main.version=$(VERSION)
endif

all:
	$(MAKE) deps
	$(MAKE) wedge

deps:
	go get github.com/sparrc/gdm
	gdm restore

wedge:
	go build -i -o $(WEDGE) -ldflags "$(LDFLAGS)" ./cmd/wedge/wedge.go

go-install:
	go install -ldflags "-w -s $(LDFLAGS)" ./cmd/wedge

install: wedge
	mkdir -p $(DESTDIR)$(PREFIX)/bin/
	sudo cp $(WEDGE) $(DESTDIR)$(PREFIX)/bin/

test:
	go test ./...

clean:
	-rm -f wedge

.PHONY: deps wedge install test clean
