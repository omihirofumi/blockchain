package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type WalletServer struct {
	port           uint16
	blockchainAddr string
	infoLog        *log.Logger
	errorLog       *log.Logger
}

func NewWalletServer(port uint16, blockchainAddr string) *WalletServer {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return &WalletServer{
		port:           port,
		blockchainAddr: blockchainAddr,
		infoLog:        infoLog,
		errorLog:       errorLog,
	}
}

func (ws *WalletServer) serve() error {
	ws.infoLog.Printf("Listening on port :%d ...", ws.port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", ws.port),
		Handler: ws.routes(),
	}

	return srv.ListenAndServe()
}
