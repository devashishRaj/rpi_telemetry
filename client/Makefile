build-arm:
	GOARCH=arm GOARM=6 GOOS=linux go build -o bin/main-rpi4b-arm main.go

build-arm64:
	GOARCH=arm64 GOOS=linux go build -o bin/main-aarch64 main.go

build-all: build-arm build-arm64

tidy:
	go mod tidy

.PHONY: build-arm build-arm64 build-all check

# Add a dependency for the build targets to run the tidy target
build-arm: tidy
build-arm64: tidy
build-all: tidy

