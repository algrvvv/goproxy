format:
	@goimports -w .

build:
	@mkdir -p bin/
	@go build -o bin/goproxy main.go

run: build
	@./bin/goproxy

dev: format build run
