package main

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/labstack/echo/v4"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/block"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/signature"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/wallet"
	"net/http"
)

const (
	KEY_BLOCKCHAIN = "blockchain"
)

func (bs *BlockchainServer) HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// GetBlockchainBlob はmemcachedを探して、ブロックチェーンのバイトを返す。
func (bs *BlockchainServer) GetBlockchainBlob() []byte {
	mc := memcache.New(fmt.Sprintf("%s:%d", "127.0.0.1", bs.mcPort))
	item, err := mc.Get(KEY_BLOCKCHAIN)
	if err != nil {
		bs.infoLog.Println("blockchain created")
		w := wallet.NewWallet()
		bs.infoLog.Printf("mining address is %s\n", w.BlockchainAddress())
		bc := block.NewBlockChain(w.BlockchainAddress())
		m, _ := bc.MarshalJSON()
		mc.Set(&memcache.Item{Key: KEY_BLOCKCHAIN, Value: m})
		return m
	}
	return item.Value
}

func (bs *BlockchainServer) GetBlockchain() (*block.Blockchain, error) {
	bc := &block.Blockchain{}
	bb := bs.GetBlockchainBlob()
	err := json.Unmarshal(bb, &bc)
	if err != nil {
		return nil, err
	}
	return bc, nil
}

func (bs *BlockchainServer) GetChain(c echo.Context) error {
	return c.JSONBlob(http.StatusOK, bs.GetBlockchainBlob())
}

func (bs *BlockchainServer) CreateTransactions(c echo.Context) error {
	tr := &block.TransactionRequest{}
	if err := c.Bind(tr); err != nil {
		bs.errorLog.Println(err)
		return bs.errResponse(http.StatusInternalServerError, err.Error())
	}

	if !tr.Validate() {
		bs.errorLog.Printf("missing field(s): %v\n", tr)
		return bs.errResponse(http.StatusBadRequest, "missing field(s)")
	}

	publicKey, err := signature.PublicKeyFromString(*tr.SenderPublicKey)
	if err != nil {
		bs.errorLog.Println(err)
		return bs.errResponse(http.StatusBadRequest, err.Error())
	}
	sg, err := signature.SignatureFromString(*tr.Signature)
	if err != nil {
		bs.errorLog.Println(err)
		return bs.errResponse(http.StatusBadRequest, err.Error())
	}
	bc, err := bs.GetBlockchain()
	if err != nil {
		bs.errorLog.Println(err)
		return bs.errResponse(http.StatusInternalServerError, err.Error())
	}
	err = bc.AddTransaction(*tr.SenderBlockchainAddress, *tr.RecipientBlockchainAddress,
		*tr.Value, publicKey, sg)
	if err != nil {
		bs.errorLog.Println(err)
		return bs.errResponse(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, bc)
}

func (bs *BlockchainServer) Mining(c echo.Context) error {
	bc, err := bs.GetBlockchain()
	if err != nil {
		bs.errorLog.Println(err)
		return bs.errResponse(http.StatusInternalServerError, "mining failed")
	}
	err = bc.Mining()
	if err != nil {
		bs.errorLog.Println(err)
		return bs.errResponse(http.StatusInternalServerError, "mining failed")
	}
	var payload struct {
		Status bool `json:"status"`
	}
	payload.Status = true
	return c.JSON(http.StatusOK, payload)
}

func (bs *BlockchainServer) GetTotalAmount(c echo.Context) error {
	bcAddr := c.Param("blockchainAddress")
	bc, err := bs.GetBlockchain()
	if err != nil {
		bs.errorLog.Println(err)
		return bs.errResponse(http.StatusInternalServerError, err.Error())
	}
	amount := bc.GetTotalAmount(bcAddr)
	payload := struct {
		Amount float32 `json:"amount"`
	}{
		Amount: amount,
	}
	return c.JSON(http.StatusOK, payload)
}
