package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/labstack/echo/v4"
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/block"
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
		w := wallet.NewWallet()
		bc := block.NewBlockChain(w.BlockchainAddress())
		m, _ := bc.MarshalJSON()
		mc.Set(&memcache.Item{Key: KEY_BLOCKCHAIN, Value: m})
		return m
	}
	return item.Value
}

func (bs *BlockchainServer) GetBlockchain(c echo.Context) error {
	return c.JSONBlob(http.StatusOK, bs.GetBlockchainBlob())
}
