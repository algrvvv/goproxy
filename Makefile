format:
	@goimports -w .

build:
	@mkdir -p bin/
	@go build -o bin/checks cmd/checks/main.go

run: build
	@./bin/checks

dev: format build run
