package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/block"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/signature"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/wallet"
	"io"
	"net/http"
	"strconv"
)

var myWallet *wallet.Wallet

// GetWallet は、Walletを生成して返す
func (ws *WalletServer) GetWallet(c echo.Context) error {
	w := wallet.NewWallet()
	return c.JSON(http.StatusOK, w)
}

// CreateTransaction は、トランザクションを生成する
func (ws *WalletServer) CreateTransaction(c echo.Context) error {
	tr := &wallet.TransactionRequest{}
	if err := c.Bind(tr); err != nil {
		return err
	}
	if !tr.Validate() {
		ws.errorLog.Printf("missing field(s):%v", tr)
		return ws.badRequest("missing field(s)")
	}

	publicKey, err := signature.PublicKeyFromString(*tr.SenderPublicKey)
	if err != nil {
		ws.errorLog.Println(err)
		return ws.badRequest(err.Error())
	}
	privateKey := signature.PrivateKeyFromString(*tr.SenderPrivateKey, publicKey)
	value, err := strconv.ParseFloat(*tr.Value, 32)
	if err != nil {
		ws.errorLog.Println(err)
		return ws.badRequest(err.Error())
	}
	value32 := float32(value)

	transaction := wallet.NewTransaction(privateKey, publicKey, *tr.SenderBlockchainAddress, *tr.RecipientBlockchainAddress, value32)
	sg, err := transaction.GenerateSignature()
	if err != nil {
		ws.errorLog.Println(err)
		return ws.badRequest(err.Error())
	}
	signatureStr := sg.String()

	btr := &block.TransactionRequest{
		SenderPublicKey:            tr.SenderPublicKey,
		SenderBlockchainAddress:    tr.SenderBlockchainAddress,
		RecipientBlockchainAddress: tr.RecipientBlockchainAddress,
		Value:                      &value32,
		Signature:                  &signatureStr,
	}

	body, err := json.Marshal(btr)
	if err != nil {
		ws.errorLog.Println(err)
		return err
	}
	url := fmt.Sprintf("%s/transactions", ws.blockchainAddr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		ws.errorLog.Println(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ws.errorLog.Println(err)
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	ws.infoLog.Println(string(respBody))
	return c.JSON(http.StatusOK, string(respBody))
}
