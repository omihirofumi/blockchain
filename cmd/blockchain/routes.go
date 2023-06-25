package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func (bs *BlockchainServer) routes() http.Handler {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5002"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", bs.HelloWorld)
	e.GET("/chain", bs.GetChain)
	e.GET("/amount/:blockchainAddress", bs.GetTotalAmount)
	e.GET("/mining", bs.Mining)
	e.POST("/transactions", bs.CreateTransactions)
	return e
}
