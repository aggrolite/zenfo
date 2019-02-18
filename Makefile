GO_PKG_DIRS  := $(subst $(shell go list -e .),.,$(shell go list ./... | grep -v /vendor/))

all: build http

http:
	zenfo-http -dbname zenfo -dbuser postgres

build: fmt lint vet db
	go clean -i
	go install ./...
	zenfo-build -dbname zenfo -dbuser postgres

db:
	dropdb zenfo --if-exists
	createdb zenfo
	psql zenfo < zenfo.psql

vet:
	go vet $(GO_PKG_DIRS)

fmt:
	echo $(GO_PKG_DIRS)
	gofmt -s -w $(GO_PKG_DIRS)

lint:
	golint -set_exit_status $(GO_PKG_DIRS)

dbuild:
	docker build -t zenfo:v1 .

docker:
	docker run --rm -it \
	-e POSTGRES_PASSWORD=secret \
	-p 8081:8081 \
	-v $$(pwd):/code \
	-w /code \
	zenfo:v1 \
	/bin/bash
