.PHONY: server client generate test tidy

generate:
	buf generate

mock:
	mockgen -destination=internal/mocks/mock_store.go -package=mocks github.com/Chetas1/grpc-blog-service/internal/store BlogStore

server:
	go run ./server

client:
	go run ./client

tidy:
	go mod tidy

test:
	go mod tidy
	go test ./... -coverprofile=coverage.out.tmp
	grep -v "_mock.go" coverage.out.tmp > coverage.out && rm coverage.out.tmp
	go tool cover -func=coverage.out