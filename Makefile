BIN := ./bin
CMD := ./cmd
INTERNAL_SRC := ./internal/**/*.go
BLOCKCHAIN := $(BIN)/blockchain
BLOCKCHAIN_SRC := $(CMD)/blockchain
WALLET := $(BIN)/wallet
WALLET_SRC := $(CMD)/wallet

.PHONY: blockchain
blockchain: $(BLOCKCHAIN)

$(BLOCKCHAIN): $(BLOCKCHAIN_SRC) $(INTERNAL_SRC)
	@echo "Building blockchain server..."
	@go build -o bin/blockchain ./cmd/blockchain
	@echo "Done"

start_b: $(BLOCKCHAIN)
	@echo "Starting blockchain server..."
	@bin/blockchain

.PHONY: wallet
wallet: $(WALLET)

$(WALLET): $(WALLET_SRC) $(INTERNAL_SRC)
	@echo "Building wallet server..."
	@go build -o bin/wallet ./cmd/wallet
	@echo "Done"

start_w: $(WALLET)
	@echo "Starting wallet server..."
	@bin/wallet

test:
	@go test -race ./...