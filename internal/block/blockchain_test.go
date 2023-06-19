package block

import (
	"github.com/omihirofumi/crypto-demo-with-blockchain/internal/wallet"
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
}

func TestNewTransaction(t *testing.T) {
	t.Parallel()
	ts := NewTransaction("from", "to", 100)
	require.Equal(t, "from", ts.senderBlockchainAddress)
	require.Equal(t, "to", ts.recipientBlockchainAddress)
	require.Equal(t, float32(100), ts.value)
}

func TestBlockchain_VerifyTransactionSignature(t *testing.T) {
	t.Parallel()

	w1 := wallet.NewWallet()
	w2 := wallet.NewWallet()
	bc := NewBlockChain(w1.BlockchainAddress())
	wts := wallet.NewTransaction(w1.PrivateKey(), w1.PublicKey(), w1.BlockchainAddress(), w2.BlockchainAddress(), 1.0)
	s, err := wts.GenerateSignature()
	require.NoError(t, err)
	bts := NewTransaction(w1.BlockchainAddress(), w2.BlockchainAddress(), 1.0)
	actual := bc.VerifyTransactionSignature(w1.PublicKey(), s, bts)
	require.Equal(t, true, actual)
}

func TestBlockchain_AddTransaction(t *testing.T) {
	t.Parallel()

	w1 := wallet.NewWallet()
	w2 := wallet.NewWallet()
	bc := NewBlockChain(w1.BlockchainAddress())
	err := bc.AddTransaction(MINING_SENDER, w1.BlockchainAddress(), 1.0, nil, nil)
	require.NoError(t, err)
	err = bc.Mining()
	require.NoError(t, err)
	wts := wallet.NewTransaction(w1.PrivateKey(), w1.PublicKey(), w1.BlockchainAddress(), w2.BlockchainAddress(), 1.0)
	s, _ := wts.GenerateSignature()

	err = bc.AddTransaction(w1.BlockchainAddress(), w2.BlockchainAddress(), 1.0, w1.PublicKey(), s)
	require.NoError(t, err)
	wts = wallet.NewTransaction(w1.PrivateKey(), w1.PublicKey(), w1.BlockchainAddress(), w2.BlockchainAddress(), 100.0)
	s, _ = wts.GenerateSignature()
	err = bc.AddTransaction(w1.BlockchainAddress(), w2.BlockchainAddress(), 100.0, w1.PublicKey(), s)
	require.Error(t, err)
	wts = wallet.NewTransaction(w1.PrivateKey(), w1.PublicKey(), w1.BlockchainAddress(), w2.BlockchainAddress(), 100.0)
	s, _ = wts.GenerateSignature()
	err = bc.AddTransaction(w1.BlockchainAddress(), w1.BlockchainAddress(), 1.0, w1.PublicKey(), s)
	require.Error(t, err)

}
