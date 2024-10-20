.PHONY: server
server:
	go run -race main.go server --config config.yaml

.PHONY: client
client:
	go run -race main.go client --config config.yaml

.PHONY: compile
compile:
	cd pkg/proto; protoc \
		--go_out ./gen --go_opt paths=source_relative \
		wisdom.proto

	mockery

.PHONY: test
test:
	go clean -testcache
	go test -race ./...

.PHONY: build
build:
	docker build . --file Dockerfile --tag github.com/kolesnikovm/tcp-pow:latest

.PHONY: demo
demo: build
	cd demo; docker-compose up