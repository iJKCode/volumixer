
.PHONY: serve-debug
serve-debug:
	go build -race -o out/server-debug cmd/server/main.go
	out/server-debug

.PHONY: test
test:
	go test ./pkg/...

.PHONY: test-coverage
test-coverage:
	go test -coverpkg=./pkg/... ./pkg/...
