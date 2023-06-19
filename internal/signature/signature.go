package signature

import (
	"fmt"
	"math/big"
)

// Signature は署名を表す
type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%x%X", s.R, s.S)
}
