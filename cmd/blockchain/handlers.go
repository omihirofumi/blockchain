package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (bs *BlockchainServer) HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
