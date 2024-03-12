generate:
	go generate ./...
	# go-bindata -fs -nomemcopy -pkg data -o data/assets_gen.go db/migrations/
	# gqlgen

build:
	go build -o ./bin/ -v -ldflags="-w -s" ./main/

dependencies:
	go mod download
	go install -ldflags="-w -s" github.com/99designs/gqlgen@latest
	go install -ldflags="-w -s" github.com/go-bindata/go-bindata/...@latest
	go install -ldflags="-w -s" github.com/amacneil/dbmate/v2/...@latest

.PHONY: build dependencies generate
.DEFAULT_GOAL := build
