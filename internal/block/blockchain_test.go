package block

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBlockchain(t *testing.T) {
	t.Parallel()
	bc := NewBlockChain("my address")
	require.Equal(t, 0, bc.chain[0].nonce)
	require.Equal(t, "my address", bc.blockchainAddress)
	require.Equal(t, 1, len(bc.chain))

	wantHash := bc.LastBlock().Hash()
	b1 := bc.CreateBlock(1, bc.LastBlock().Hash())
	require.Equal(t, 1, b1.nonce)
	require.Equal(t, wantHash, b1.previousHash)
	wantHash = bc.LastBlock().Hash()
	b2 := bc.CreateBlock(2, bc.LastBlock().Hash())
	require.Equal(t, 2, b2.nonce)
	require.Equal(t, wantHash, b2.previousHash)

	require.Equal(t, 3, len(bc.chain))

	bc.AddTransaction("from1", "to1", 100)
	bc.AddTransaction("from2", "to2", 200)
	bc.AddTransaction("from3", "to3", 300)
	wantTransaction := NewTransaction("from1", "to1", 100)

	require.Equal(t, wantTransaction, bc.transactionPool[0])
	require.Equal(t, 3, len(bc.transactionPool))
}

func TestAddTransaction(t *testing.T) {
	t.Parallel()
	ts := NewTransaction("from", "to", 100)
	require.Equal(t, "from", ts.senderBlockchainAddr)
	require.Equal(t, "to", ts.recipientBlockchainAddr)
	require.Equal(t, float32(100), ts.value)
}

func TestMining(t *testing.T) {
	t.Parallel()
	bc := NewBlockChain("my_address")
	bc.AddTransaction("A", "B", 1.0)
	result := bc.Mining()
	require.Equal(t, true, result)
	want := &Transaction{MINING_SENDER, "my_address", MINING_REWARD}
	require.Equal(t, want, bc.chain[1].transactions[1])
}

func TestGetTotalAmount(t *testing.T) {
	t.Parallel()
	bc := NewBlockChain("first_blockchain_address")
	bc.AddTransaction("first_blockchain_address", "A", 1.0)
	bc.AddTransaction("first_blockchain_address", "B", 1.0)
	bc.Mining()
	require.Equal(t, float32(1.0), bc.GetTotalAmount("A"))
	require.Equal(t, float32(1.0), bc.GetTotalAmount("B"))
	require.Equal(t, float32(MINING_REWARD-2), bc.GetTotalAmount("first_blockchain_address"))
	bc.AddTransaction("A", "first_blockchain_address", 1.0)
	bc.AddTransaction("B", "first_blockchain_address", 1.0)
	bc.Mining()
	require.Equal(t, float32(0.0), bc.GetTotalAmount("A"))
	require.Equal(t, float32(0.0), bc.GetTotalAmount("B"))
	require.Equal(t, float32(MINING_REWARD*2), bc.GetTotalAmount("first_blockchain_address"))

}
