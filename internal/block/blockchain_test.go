package block

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBlockchain(t *testing.T) {
	bc := NewBlockChain()
	require.Equal(t, 0, bc.chain[0].nonce)
	require.Equal(t, 1, len(bc.chain))

	b1 := bc.CreateBlock(1, bc.LastBlock().Hash())
	require.Equal(t, 1, b1.nonce)
	require.Equal(t, bc.LastBlock().Hash(), b1.previousHash)
	b2 := bc.CreateBlock(2, bc.LastBlock().Hash())
	require.Equal(t, 2, b2.nonce)
	require.Equal(t, bc.LastBlock().Hash(), b2.previousHash)

	require.Equal(t, 3, len(bc.chain))
}
