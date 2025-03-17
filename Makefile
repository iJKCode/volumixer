

.PHONY: serve-debug
serve-debug: protobuf
	go build -race -o out/server-debug cmd/server/main.go
	out/server-debug

.PHONY: test
test: protobuf
	go test ./pkg/...

.PHONY: test-coverage
test-coverage: protobuf
	go test -coverpkg=./pkg/... ./pkg/...

.PHONY: protobuf
protobuf:
	buf generate

.PHONY: deps
deps:
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	#go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
