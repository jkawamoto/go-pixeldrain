#
# Makefile
#
# Copyright (c) 2018 Junpei Kawamoto
#
# This software is released under the MIT License.
#
# http://opensource.org/licenses/mit-license.php
#

#
# Makefile
#
VERSION = snapshot
GHRFLAGS =
.PHONY: build release get-deps test

default: build

build:
	goxc -d=pkg -pv=$(VERSION) -n=pd -os="linux,darwin,windows" -wd=cmd/pd

release:
	ghr  -u jkawamoto $(GHRFLAGS) $(VERSION) pkg/$(VERSION)

test:
	go test -v ./...
