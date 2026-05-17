BINARY=pi-router
PREFIX=/usr/local

.PHONY: build install clean

build:
	go build -o $(BINARY) ./cmd/pi-router

install: build
	install -m 755 $(BINARY) $(PREFIX)/bin/$(BINARY)

clean:
	rm -f $(BINARY)
