.PHONY: server client generate test tidy

generate:
	buf generate

server:
	go run ./server

client:
	go run ./client

tidy:
	go mod tidy