package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/signature"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	MINING_DIFFICULTY = 5
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

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previousHash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}

func (b *Block) UnmarshalJSON(data []byte) error {
	var d struct {
		Timestamp    int64
		Nonce        int
		PreviousHash string
		Transactions []*Transaction
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	previousHash := [32]byte([]byte(d.PreviousHash))

	b.timestamp = d.Timestamp
	b.nonce = d.Nonce
	b.previousHash = previousHash
	b.transactions = d.Transactions
	return nil
}

// Blockchain はブロックチェーンを表す
type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
	mu                sync.Mutex
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

// AddTransaction はトランザクションのプールにトランザクションを追加する。 トランザクションの検証に失敗した場合、エラーを返します。
func (bc *Blockchain) AddTransaction(sender, recipient string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *signature.Signature) error {
	t := NewTransaction(sender, recipient, value)
	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return nil
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		balance := bc.GetTotalAmount(sender)
		if balance < value {
			return fmt.Errorf("残高不足です。残高:%f", balance)
		}
		bc.transactionPool = append(bc.transactionPool, t)
		return nil
	}

	return fmt.Errorf("トランザクションの検証に失敗しました。")
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
		transactions = append(transactions, NewTransaction(t.senderBlockchainAddress, t.senderBlockchainAddress, t.value))
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

// Mining は、マイニングを行うメソッド
func (bc *Blockchain) Mining() error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	err := bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	if err != nil {
		return err
	}
	nonce := bc.ProofOfWork()
	log.Println(nonce)
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	return nil
}

// GetTotalAmount は対象ブロックチェーンアドレスが所持している額を取得
func (bc *Blockchain) GetTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			if t.recipientBlockchainAddress == blockchainAddress {
				totalAmount += t.value
			}
			if t.senderBlockchainAddress == blockchainAddress {
				totalAmount -= t.value
			}
		}
	}
	return totalAmount
}

func (bc *Blockchain) ValidChain() bool {
	chain := bc.chain
	if len(chain) <= 1 {
		return true
	}
	block := chain[0]
	for _, b := range chain[1:] {
		if b.previousHash != block.Hash() {
			return false
		}

		if !bc.ValidProof(b.nonce, b.previousHash, b.transactions, MINING_DIFFICULTY) {
			log.Println(b.nonce)
			return false
		}
		block = b
	}
	return true
}

func (bc *Blockchain) VerifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey, s *signature.Signature, t *Transaction) bool {
	m, _ := json.Marshal(struct {
		SenderBlockchainAddress    string
		RecipientBlockchainAddress string
		Value                      float32
	}{
		SenderBlockchainAddress:    t.senderBlockchainAddress,
		RecipientBlockchainAddress: t.recipientBlockchainAddress,
		Value:                      t.value,
	})

	h := sha256.Sum256(m)
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)

}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Chain []*Block `json:"chain"`
	}{
		Chain: bc.chain,
	})
}

func (bc *Blockchain) UnmarshalJSON(data []byte) error {
	var d struct {
		Chain []*Block
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	bc.chain = d.Chain
	return nil
}

func (bc *Blockchain) Print() {
	for _, b := range bc.chain {
		b.Print()
	}
}

// Transaction はトランザクションを表す
type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

// NewTransaction はトランザクションを生成する
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{
		senderBlockchainAddress:    sender,
		recipientBlockchainAddress: recipient,
		value:                      value,
	}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockchainAddress    string  `json:"senderBlockchainAddress"`
		RecipientBlockchainAddress string  `json:"recipientBlockchainAddress"`
		Value                      float32 `json:"value"`
	}{
		SenderBlockchainAddress:    t.senderBlockchainAddress,
		RecipientBlockchainAddress: t.recipientBlockchainAddress,
		Value:                      t.value,
	})
}

func (t *Transaction) Print() {
	fmt.Printf("%s %s %s\n", strings.Repeat("=", 5), "Transaction", strings.Repeat("=", 5))
	fmt.Printf("sender: %s\n", t.senderBlockchainAddress)
	fmt.Printf("recipient: %s\n", t.recipientBlockchainAddress)
	fmt.Printf("value: %f\n", t.value)
}

type TransactionRequest struct {
	SenderPublicKey            *string  `json:"sender_public_key"`
	SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderPublicKey == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.Signature == nil ||
		tr.Value == nil {
		return false
	}
	return true
}
