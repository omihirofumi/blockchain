bin := ./bin
cmd := ./cmd
blockchain := $(bin)/blockchain

.PHONY: build
build: $(blockchain)

$(blockchain): ./cmd/blockchain
	@echo "Building blockchain server..."
	@go build -o bin/blockchain ./cmd/blockchain
	@echo "Done"

start: $(blockchain)
	@echo "Starting blockchain server..."
	@bin/blockchain

test:
	@go test -race ./...