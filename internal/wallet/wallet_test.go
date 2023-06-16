package wallet

import (
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestNewWallet(t *testing.T) {
	w, err := NewWallet()
	require.NoError(t, err)
	log.Println(w.PublicKeyStr())
	log.Println(w.PrivateKeyStr())
	log.Println(w.BlockchainAddress())
}
