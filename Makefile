.PHONY: build run test clean

build:
	mkdir -p bin
	go build -o bin/ebenezer-cli ./main.go

run:
	go run ./main.go $(opts)

test:
	go test -json -skip /pkg/test -v ./... $(args) 2>&1 | gotestfmt

test-match:
	make test args="-run $(case)"

test-ci:
	go test -json -skip /pkg/test -v ./... 2>&1 | gotestfmt

clean:
	go clean
	rm -f ebenezer-cli

install-local: build
	@echo "Installing ebenezer-cli locally..."
	rm -rf $(HOME)/.local/bin/ebenezer-cli
	cp bin/ebenezer-cli $(HOME)/.local/bin/ebenezer-cli