package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (ws *WalletServer) routes() http.Handler {
	e := echo.New()

	e.GET("/wallet", ws.GetWallet)
	e.POST("/transaction", ws.CreateTransaction)

	return e
}
