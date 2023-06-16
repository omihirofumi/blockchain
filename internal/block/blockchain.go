package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 1
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 10.0
)

// Block はブロックチェーンにおける各ブロックを表す
type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

// NewBlock はブロックを生成する
func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions,
	}
}

// Hash はブロックのハッシュ値を計算する
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(struct {
		Timestamp    int64
		Nonce        int
		PreviousHash [32]byte
		Transactions []*Transaction
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
	return sha256.Sum256([]byte(m))
}

func (b *Block) Print() {
	fmt.Printf("%s %s %s\n", strings.Repeat("=", 5), "Block", strings.Repeat("=", 5))
	for _, t := range b.transactions {
		t.Print()
	}
	fmt.Printf("nonce: %d\n", b.nonce)
	fmt.Println(strings.Repeat("=", 15))
}

// Blockchain はブロックチェーンを表す
type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

// NewBlockChain はgenesisブロックを追加したブロックチェーンを生成する
func NewBlockChain(blockchainAddress string) *Blockchain {
	// genesis block
	b := &Block{}
	// このブロックチェーン作成者のアドレス
	bc := &Blockchain{blockchainAddress: blockchainAddress}
	bc.CreateBlock(0, b.Hash())
	return bc
}

// CreateBlock はブロックを生成しブロックチェーンへ追加する
func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

// LastBlock は前回のブロックを返す
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// AddTransaction はトランザクションのプールにトランザクションを追加する
func (bc *Blockchain) AddTransaction(sender, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

// ValidProof はナンスの検証をする
func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

// CopyTransactionPool は現時点でのトランザクションプールをコピーする
func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions, NewTransaction(t.senderBlockchainAddr, t.senderBlockchainAddr, t.value))
	}
	return transactions
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	return true
}

// GetTotalAmount は対象ブロックチェーンアドレスが所持している額を取得
func (bc *Blockchain) GetTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			if t.recipientBlockchainAddr == blockchainAddress {
				totalAmount += t.value
			}
			if t.senderBlockchainAddr == blockchainAddress {
				totalAmount -= t.value
			}
		}
	}
	return totalAmount
}

func (bc *Blockchain) Print() {
	for _, b := range bc.chain {
		b.Print()
	}
}

// Transaction はトランザクションを表す
type Transaction struct {
	senderBlockchainAddr    string
	recipientBlockchainAddr string
	value                   float32
}

// NewTransaction はトランザクションを生成する
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{
		senderBlockchainAddr:    sender,
		recipientBlockchainAddr: recipient,
		value:                   value,
	}
}

func (t *Transaction) Print() {
	fmt.Printf("%s %s %s\n", strings.Repeat("=", 5), "Transaction", strings.Repeat("=", 5))
	fmt.Printf("sender: %s\n", t.senderBlockchainAddr)
	fmt.Printf("recipient: %s\n", t.recipientBlockchainAddr)
	fmt.Printf("value: %f\n", t.value)
}
