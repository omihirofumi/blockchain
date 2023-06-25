package main

import "log"

func main() {
	bs := NewBlockchainServer(5001)

	log.Fatal(bs.serve())
}
