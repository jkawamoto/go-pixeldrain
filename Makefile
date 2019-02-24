#
# Makefile
#
# Copyright (c) 2018-2019 Junpei Kawamoto
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
.PHONY: build release get-deps test client

default: build

client:
	swagger generate client -f https://raw.githubusercontent.com/jkawamoto/pixeldrain-swagger/master/swagger.yaml -t .
	sed -i -e "s/\*string/string/" models/standard_error.go

build:
	goxc -d=pkg -pv=$(VERSION) -n=pd -os="linux,darwin,windows" -wd=cmd/pd

release:
	ghr  -u jkawamoto $(GHRFLAGS) $(VERSION) pkg/$(VERSION)

test:
	go test -v ./...
