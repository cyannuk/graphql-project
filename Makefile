gen:
	go generate ./config/... ./domain/...

bin-gen:
	go-bindata -fs -nomemcopy -pkg data -o data/assets_gen.go db/migrations/

gql-gen: gen
	gqlgen

build:
	go build -o ./bin/ -v -ldflags="-w -s" ./main/

it-tests:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/ -v -ldflags="-w -s" ./main/
	docker compose -f docker-compose.yml up -d

deps:
	go mod download
	go install -ldflags="-w -s" github.com/99designs/gqlgen@latest
	go install -ldflags="-w -s" github.com/go-bindata/go-bindata/...@latest
	go install -ldflags="-w -s" github.com/amacneil/dbmate/v2/...@latest

.PHONY: build deps gen bin-gen gql-gen
.DEFAULT_GOAL := build
