package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type BlockchainServer struct {
	port     uint16
	mcPort   uint16
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewBlockchainServer(port uint16, mcPort uint16) *BlockchainServer {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return &BlockchainServer{
		port:     port,
		mcPort:   mcPort,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

func (bs *BlockchainServer) serve() error {
	bs.infoLog.Printf("Listening on port :%d ...", bs.port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", bs.port),
		Handler: bs.routes(),
	}

	return srv.ListenAndServe()
}
