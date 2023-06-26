package main

import (
	"github.com/labstack/echo/v4"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/block"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/signature"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/wallet"
	"net/http"
)

const (
	KEY_BLOCKCHAIN = "blockchain"
)

var bcCache *block.Blockchain

func (bs *BlockchainServer) HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// GetBlockchain は、ブロックチェーンを返す。
func (bs *BlockchainServer) GetBlockchain() *block.Blockchain {
	if bcCache == nil {
		bs.infoLog.Println("blockchain created")
		w := wallet.NewWallet()
		bs.infoLog.Printf("mining address is %s\n", w.BlockchainAddress())
		bc := block.NewBlockChain(w.BlockchainAddress())
		bcCache = bc
		return bc
	}
	return bcCache
}

func (bs *BlockchainServer) GetChain(c echo.Context) error {
	return c.JSON(http.StatusOK, bs.GetBlockchain())
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
	bc := bs.GetBlockchain()
	err = bc.AddTransaction(*tr.SenderBlockchainAddress, *tr.RecipientBlockchainAddress,
		*tr.Value, publicKey, sg)
	if err != nil {
		bs.errorLog.Println(err)
		return bs.errResponse(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, bc)
}

func (bs *BlockchainServer) Mining(c echo.Context) error {
	bc := bs.GetBlockchain()
	err := bc.Mining()
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

// GetTotalAmount は、対象アドレスの総額を計算
func (bs *BlockchainServer) GetTotalAmount(c echo.Context) error {
	bcAddr := c.Param("blockchainAddress")
	bc := bs.GetBlockchain()
	amount := bc.GetTotalAmount(bcAddr)
	payload := struct {
		Amount float32 `json:"amount"`
	}{
		Amount: amount,
	}
	return c.JSON(http.StatusOK, payload)
}

// VerifyChain は、チェーンが不正ではないか検証する。
func (bs *BlockchainServer) VerifyChain(c echo.Context) error {
	bc := bs.GetBlockchain()
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	if bc.ValidChain() {
		payload.Error = false
		payload.Message = "blockchain is valid!"
	} else {
		payload.Error = true
		payload.Message = "block chain invalid.."
	}
	return c.JSON(http.StatusOK, payload)
}
