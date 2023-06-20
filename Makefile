build:
	@echo "Building blockchain server..."
	@go build -o bin/blockchain ./cmd/blockchain
	@echo "Done"

start: build
	@echo "Starting blockchain server..."
	@bin/blockchain

test:
	@go test -race ./...