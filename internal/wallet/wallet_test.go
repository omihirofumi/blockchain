package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestNewWallet(t *testing.T) {
	w := NewWallet()
	log.Println(w.PublicKeyStr())
	log.Println(w.PrivateKeyStr())
	log.Println(w.BlockchainAddress())
}

func TestNewTransaction(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := &privateKey.PublicKey
	ts := NewTransaction(privateKey, publicKey, "A", "B", 1.0)
	s, err := ts.GenerateSignature()
	require.NoError(t, err)
	log.Println(s.String())
}
