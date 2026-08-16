BINARY := cauldron
VERSION ?= dev

.PHONY: build test vet fmt check clean

build:
	go build -ldflags "-X github.com/CauldronUp/cauldron/internal/cli.Version=$(VERSION)" -o $(BINARY) ./cmd/cauldron

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -l -w .

check: fmt vet test

clean:
	rm -f $(BINARY) $(BINARY).exe
