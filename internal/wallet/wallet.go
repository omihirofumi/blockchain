package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/signature"
	"golang.org/x/crypto/ripemd160"
)

// Wallet は仮想通貨のウォレットを表す
type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

// NewWallet はウォレットを生成する
func NewWallet() *Wallet {
	// 参考：https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
	// 0 - Having a private ECDSA key
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	// 1 - Take the corresponding public key generated with it (33 bytes, 1 byte 0x02 (y-coord is even), and 32 bytes corresponding to X coordinate)
	publicKey := &privateKey.PublicKey

	// 2 - Perform SHA-256 hashing on the public key
	h2 := sha256.New()
	h2.Write(publicKey.X.Bytes())
	h2.Write(publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	// 3 - Perform RIPEMD-160 hashing on the result of SHA-256
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	//4 - Add version byte in front of RIPEMD-160 hash (0x00 for Main Network)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	// 5 - Perform SHA-256 hash on the extended RIPEMD-160 result
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	// 6 - Perform SHA-256 hash on the result of the previous SHA-256 hash
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	// 7 - Take the first 4 bytes of the second SHA-256 hash. This is the address checksum
	chsum := digest6[:4]

	// 8 - Add the 4 checksum bytes from stage 7 at the end of extended RIPEMD-160 hash from stage 4. This is the 25-byte binary Bitcoin Address.
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4)
	copy(dc8[21:], chsum)

	// 9 - Convert the result from a byte string into a base58 string using Base58Check encoding. This is the most commonly used Bitcoin Address format
	address := base58.Encode(dc8)

	return &Wallet{
		privateKey:        privateKey,
		publicKey:         &privateKey.PublicKey,
		blockchainAddress: address,
	}
}

// PrivateKey はプライベートキーを返す
func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

// PrivateKeyStr はプライベートキーの文字列を返す
func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

// PublicKey はプライベートキーを返す
func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

// PublicKeyStr はプライベートキーの文字列を返す
func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

// BlockchainAddress はブロックチェーンアドレスを返す
func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"privateKey"`
		PublicKey         string `json:"publicKey"`
		BlockchainAddress string `json:"blockchainAddress"`
	}{
		PrivateKey:        w.PrivateKeyStr(),
		PublicKey:         w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),
	})
}

// Transaction トランザクションを表す
type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

// NewTransaction はトランザクションを生成する
func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey,
	sender string, recipient string, value float32) *Transaction {
	return &Transaction{
		senderPrivateKey:           privateKey,
		senderPublicKey:            publicKey,
		senderBlockchainAddress:    sender,
		recipientBlockchainAddress: recipient,
		value:                      value,
	}
}

// GenerateSignature は署名を生成する
func (t *Transaction) GenerateSignature() (*signature.Signature, error) {
	m, err := json.Marshal(struct {
		SenderBlockchainAddress    string
		RecipientBlockchainAddress string
		Value                      float32
	}{
		SenderBlockchainAddress:    t.senderBlockchainAddress,
		RecipientBlockchainAddress: t.recipientBlockchainAddress,
		Value:                      t.value,
	})

	if err != nil {
		return nil, err
	}

	h := sha256.Sum256(m)
	r, s, err := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	if err != nil {
		return nil, err
	}
	return &signature.Signature{R: r, S: s}, nil
}

type TransactionRequest struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderPublicKey            *string `json:"sender_public_key"`
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	Value                      *string `json:"value"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderPublicKey == nil ||
		tr.SenderPrivateKey == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.Value == nil {
		return false
	}
	return true
}
