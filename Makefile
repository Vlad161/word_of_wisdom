PWD := $(PWD)
export PATH := $(PWD)/bin:$(PATH)

tools:
	cd tools && go generate -tags tools

.PHONY: test
test:
	@go test -race -count 1 ./...

generate:
	@go generate ./...

run_client:
	go run client/cmd/main.go

run_server:
	go run server/cmd/main.go

run: run_server run_client # make run -j2

run_docker_compose:
	@docker compose up --build

docker_build_client:
	@docker build -t word_of_wisdom_client -f client/Dockerfile .

docker_build_server:
	@docker build -t word_of_wisdom_server -f server/Dockerfile .