GO_PKG_DIRS  := $(subst $(shell go list -e .),.,$(shell go list ./... | grep -v /vendor/))

all: build http

http:
	zenfo-http

build: fmt lint vet db
	go clean -i
	go install ./...
	zenfo-build

db:
	dropdb zenfo
	createdb zenfo
	psql zenfo < zenfo.psql

vet:
	go vet $(GO_PKG_DIRS)

fmt:
	echo $(GO_PKG_DIRS)
	gofmt -s -w $(GO_PKG_DIRS)

lint:
	golint -set_exit_status $(GO_PKG_DIRS)
