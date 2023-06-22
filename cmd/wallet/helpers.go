package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (ws *WalletServer) badRequest(msg string) error {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Error = true
	payload.Message = msg
	return echo.NewHTTPError(http.StatusBadRequest, payload)
}
