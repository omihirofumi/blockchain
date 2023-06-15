package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (bs *BlockchainServer) routes() http.Handler {
	e := echo.New()

	e.GET("/", bs.HelloWorld)

	return e
}
