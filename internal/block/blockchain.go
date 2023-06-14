package block

import (
	"crypto/sha256"
	"encoding/json"
	"time"
)

// Block はブロックチェーンにおける各ブロックを表す
type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transaction  []string
}

// NewBlock はブロックを生成する
func NewBlock(nonce int, previousHash [32]byte) *Block {
	return &Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

// Blockchain はブロックチェーンを表す
type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

// NewBlockChain はgenesisブロックを追加したブロックチェーンを生成する
func NewBlockChain() *Blockchain {
	// genesis block
	b := &Block{}
	bc := &Blockchain{}
	bc.CreateBlock(0, b.Hash())
	return bc
}

// CreateBlock はブロックを生成しブロックチェーンへ追加する
func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

// LastBlock は前回のブロックを返す
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}
