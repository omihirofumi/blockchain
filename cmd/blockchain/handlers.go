package main

import (
	"github.com/labstack/echo/v4"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/block"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/wallet"
	"net/http"
)

func (bs *BlockchainServer) HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (bs *BlockchainServer) GetBlockchain(c echo.Context) error {
	w := wallet.NewWallet()
	bc := block.NewBlockChain(w.BlockchainAddress())
	bc.Print()

	m, err := bc.MarshalJSON()
	if err != nil {
		return err
	}
	bs.infoLog.Println(string(m))
	return c.JSONBlob(http.StatusOK, m)
}
