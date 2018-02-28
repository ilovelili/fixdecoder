VERSION = $(shell grep 'version =' version.go | sed -E 's/.*"(.+)"$$/\1/')

default: all

all: build

deps:
	go get -d -v -u github.com/tidwall/gjson

build: deps
	go build -o fixdecoder

test:
	go test
		
version:
	@echo $(VERSION)

.PTHONY: all deps build version test