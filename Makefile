BINARY := minesweeper-go
PREFIX ?= /usr/local

.PHONY: all build install uninstall test vet fmt clean run

all: build

build:
	go build -o $(BINARY) .

install: build
	install -d $(PREFIX)/bin
	install -m 0755 $(BINARY) $(PREFIX)/bin/$(BINARY)

uninstall:
	rm -f $(PREFIX)/bin/$(BINARY)

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

run:
	go run .

clean:
	rm -f $(BINARY)
