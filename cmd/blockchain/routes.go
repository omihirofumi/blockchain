package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (bs *BlockchainServer) routes() http.Handler {
	r := mux.NewRouter()

	// 仮
	r.HandleFunc("/", bs.HelloWorld)

	return r
}
