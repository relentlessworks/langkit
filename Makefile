.PHONY: build test vet clean run

BINARY=langkit
CMD=cmd/langkit

build:
	CGO_ENABLED=0 go build -o $(BINARY) ./$(CMD)

test:
	go test -v ./...

vet:
	go vet ./...

run: build
	./$(BINARY) -no-auth

clean:
	rm -f $(BINARY)

all: vet test build
