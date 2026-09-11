APP_NAME := interpreter-go
BUILD_DIR := bin
BUILD_OUTPUT := $(BUILD_DIR)/$(APP_NAME)
CODECRAFTERS_OUTPUT := /tmp/codecrafters-build-interpreter-go

.PHONY: all build build-codecrafters run test test-v test-dispatch test-app fmt vet clean

all: fmt vet test build

build:
	mkdir -p $(BUILD_DIR)
	go build -v -o $(BUILD_OUTPUT) ./app/cmd

build-codecrafters:
	go build -v -o $(CODECRAFTERS_OUTPUT) ./app/cmd

run: build
	./$(BUILD_OUTPUT)

test:
	go test ./...

test-v:
	go test ./... -v

test-dispatch:
	go test ./app/dispatch -v

test-app:
	go test ./app -v

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(CODECRAFTERS_OUTPUT)

generate:
	go generate ./...