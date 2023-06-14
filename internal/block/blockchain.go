package block

import (
	"time"
)

// Block はブロックチェーンにおける各ブロックを表す
type Block struct {
	timestamp    int64
	nonce        int
	previousHash string
	transaction  []string
}

// NewBlock はブロックを生成する
func NewBlock(nonce int, previousHash string) *Block {
	return &Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
	}
}

// Blockchain はブロックチェーンを表す
type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

// NewBlockChain はgenesisブロックを追加したブロックチェーンを生成する
func NewBlockChain() *Blockchain {
	bc := &Blockchain{}
	bc.CreateBlock(0, "genesis hash")
	return bc
}

// CreateBlock はブロックを生成しブロックチェーンへ追加する
func (bc *Blockchain) CreateBlock(nonce int, previousHash string) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}
