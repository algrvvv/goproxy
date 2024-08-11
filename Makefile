format:
	@goimports -w .

build:
	@mkdir -p bin/
	@go build -o bin/goproxy cmd/goproxy/main.go

run: build
	@./bin/goproxy

dev: format build run

releases:
	GOOS=linux GOARCH=amd64 go build -o bin/linux_goproxy cmd/goproxy/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/mac_arm_goproxy cmd/goproxy/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/windows_goproxy.exe cmd/goproxy/main.go
