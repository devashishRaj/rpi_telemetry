build-arm:
	GOARCH=arm GOARM=6 GOOS=linux go build -o bin/main-rpi4b-arm main.go

build-arm64:
	GOARCH=arm64 GOOS=linux go build -o bin/main-aarch64 main.go

build-all: build-arm build-arm64
