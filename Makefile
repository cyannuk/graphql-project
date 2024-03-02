generate:
	go generate ./...

build:
	go build -o ./bin/ -v -ldflags="-w -s" ./main/

dependencies:
	go mod download

.PHONY: build dependencies gen
.DEFAULT_GOAL := build
