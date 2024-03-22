gen:
	go generate ./config/... ./domain/...

bin-gen:
	go-bindata -fs -nomemcopy -pkg data -o data/assets_gen.go db/migrations/

gql-gen: gen
	gqlgen

build:
	go build -o ./bin/ -v -ldflags="-w -s" ./cmd/service/

start-test-env:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/ ./cmd/service/
	docker compose -f docker-compose-tests.yml --env-file .env --env-file test.env up -d

stop-test-env:
	docker compose -f docker-compose-tests.yml down

integration-tests: start-test-env
	go clean -testcache
	go test -v ./tests/...

deps:
	go mod download
	go install -ldflags="-w -s" github.com/99designs/gqlgen@latest
	go install -ldflags="-w -s" github.com/go-bindata/go-bindata/...@latest
	go install -ldflags="-w -s" github.com/amacneil/dbmate/v2/...@latest

.PHONY: build deps gen bin-gen gql-gen start-test-env stop-test-env integration-tests
.DEFAULT_GOAL := build
