SHELL=/bin/bash
DIR=$(shell pwd)
ORANGE=\e[0;33m
NOCOLOR=\e[0m

.PHONY: build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o odr_linux_amd64 -ldflags="-s -w" main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o odr_darwin_amd64 -ldflags="-s -w" main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o odr_windows_amd64 -ldflags="-s -w" main.go
