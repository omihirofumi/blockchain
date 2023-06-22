package main

import "log"

func main() {
	ws := NewWalletServer(5002, "http://localhost:5001")

	log.Fatal(ws.serve())
}
