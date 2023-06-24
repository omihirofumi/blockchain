package main

import (
	"github.com/labstack/echo/v4"
)

func (bs *BlockchainServer) errResponse(statusCode int, msg string) error {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Error = true
	payload.Message = msg
	return echo.NewHTTPError(statusCode, payload)
}
