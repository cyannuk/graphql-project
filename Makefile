gen:
	go generate ./...

bin-gen:
	go-bindata -fs -nomemcopy -pkg data -o data/assets_gen.go db/migrations/

gql-gen: gen
	gqlgen

build:
	go build -o ./bin/ -v -ldflags="-w -s" ./main/

deps:
	go mod download
	go install -ldflags="-w -s" github.com/99designs/gqlgen@latest
	go install -ldflags="-w -s" github.com/go-bindata/go-bindata/...@latest
	go install -ldflags="-w -s" github.com/amacneil/dbmate/v2/...@latest

.PHONY: build deps gen bin-gen gql-gen
.DEFAULT_GOAL := build
