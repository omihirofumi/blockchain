package signature

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Signature は署名を表す
type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%064x%064x", s.R, s.S)
}

func string2BigIntTuple(s string) (big.Int, big.Int, error) {
	if len(s) != 128 {
		return big.Int{}, big.Int{}, fmt.Errorf("%s", "invalid publickey")
	}
	bx, _ := hex.DecodeString(s[:64])
	by, _ := hex.DecodeString(s[64:])

	var bix big.Int
	var biy big.Int

	bix.SetBytes(bx)
	bix.SetBytes(by)
	return bix, biy, nil
}

func SignatureFromString(s string) (*Signature, error) {
	x, y, err := string2BigIntTuple(s)
	if err != nil {
		return nil, err
	}
	return &Signature{R: &x, S: &y}, nil
}

func PublicKeyFromString(s string) (*ecdsa.PublicKey, error) {
	x, y, err := string2BigIntTuple(s)
	if err != nil {
		return nil, err
	}

	return &ecdsa.PublicKey{elliptic.P256(), &x, &y}, nil
}

func PrivateKeyFromString(s string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, _ := hex.DecodeString(s[:])
	var bi big.Int
	bi.SetBytes(b)
	return &ecdsa.PrivateKey{PublicKey: *publicKey, D: &bi}
}
