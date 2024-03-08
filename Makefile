generate:
	go generate ./...

build:
	go build -o ./bin/ -v -ldflags="-w -s" ./main/

dependencies:
	go mod download
	go install -ldflags="-w -s" github.com/99designs/gqlgen@latest

.PHONY: build dependencies generate
.DEFAULT_GOAL := build
