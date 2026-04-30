.PHONY: build test lint ci clean

BINARY_NAME = redisrelay

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./...

lint:
	go vet ./...

ci: test lint

clean:
	rm -f $(BINARY_NAME)
