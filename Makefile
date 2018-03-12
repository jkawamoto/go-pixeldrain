#
# Makefile
#
VERSION = snapshot
GHRFLAGS =
.PHONY: build release get-deps test

default: build

build:
	goxc -d=pkg -pv=$(VERSION) -n=pd -os="linux,darwin,windows"

release:
	ghr  -u jkawamoto $(GHRFLAGS) v$(VERSION) pkg/$(VERSION)

get-deps:
	go get -d -t -v .

test:
	go test -v ./...
