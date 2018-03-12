#
# Makefile
#
VERSION = snapshot
GHRFLAGS =
.PHONY: build release

default: build

build:
	goxc -d=pkg -pv=$(VERSION) -n=pd -os="linux,darwin,windows"

release:
	ghr  -u jkawamoto $(GHRFLAGS) v$(VERSION) pkg/$(VERSION)
