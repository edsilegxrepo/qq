SRC = ./
ifeq ($(OS),Windows_NT)
	BINARY = qq.exe
else
	BINARY = qq
endif
DESTDIR = ~/.local/bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS = -X github.com/JFryy/qq/cli.Version=$(VERSION)
GO_OPTS = -trimpath -buildmode=pie -s -w

all: build

build:
	go build -ldflags "$(LDFLAGS)" $(GO_OPTS) -o bin/$(BINARY) $(SRC)

test: build
	./tests/test.sh
	go test ./... -v -cover

lint:
	gofumpt -l -w .
	go vet ./...
	golangci-lint run ./... --no-config
	govulncheck ./...
	gosec ./...

clean:
	rm -f bin/$(BINARY) qq_test_binary coverage.out coverage.html
	go clean -testcache

install: build test
	mkdir -p $(DESTDIR)
	cp bin/$(BINARY) $(DESTDIR)

docker-push:
	docker buildx build --platform linux/amd64,linux/arm64 . -t jfryy/qq:latest --push

.PHONY: all build test clean install docker-push
