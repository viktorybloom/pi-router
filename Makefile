BINARY := pi-router
CMD := ./cmd/pi-router
BUILD_DIR := build

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build
build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) $(CMD)

.PHONY: pi
pi:
	GOOS=linux GOARCH=arm64 \
	go build $(LDFLAGS) \
	-o $(BUILD_DIR)/$(BINARY)-arm64 \
	$(CMD)

.PHONY: run
run:
	go run $(CMD)

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
