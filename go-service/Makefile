init:
	go mod tidy
build: init
	go build -o report_srv cmd/server/main.go

start: build
	./report_srv