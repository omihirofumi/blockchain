package main

import "log"

func main() {
	bs := NewBlockchainServer(5001, 11211)

	log.Fatal(bs.serve())
}
