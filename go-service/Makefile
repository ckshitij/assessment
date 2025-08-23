init:
	go mod tidy

build: init
	go build -o report_srv cmd/server/main.go

install-dep: init
	go mod download && go install honnef.co/go/tools/cmd/staticcheck@latest && go install github.com/securego/gosec/v2/cmd/gosec@latest && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

start: fmt build
	./report_srv

test: init
	go test -cover ./...

lint: init
	golangci-lint run ./...

security: init
	staticcheck ./... && gosec ./...

fmt:
	go fmt ./...