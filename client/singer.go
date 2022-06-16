package client

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type senderFromServer struct {
	addr      common.Address
	blockhash common.Hash
}

var errNotCached = errors.New("sender not cached")

func setSenderFromServer(tx *types.Transaction, addr common.Address, block common.Hash) {
	// Use types.Sender for side-effect to store our signer into the cache.
	types.Sender(&senderFromServer{addr, block}, tx)
}

func (s *senderFromServer) Equal(other types.Signer) bool {
	os, ok := other.(*senderFromServer)
	return ok && os.blockhash == s.blockhash
}

func (s *senderFromServer) Sender(tx *types.Transaction) (common.Address, error) {
	if s.addr == (common.Address{}) {
		return common.Address{}, errNotCached
	}
	return s.addr, nil
}

func (s *senderFromServer) ChainID() *big.Int {
	panic("can't sign with senderFromServer")
}
func (s *senderFromServer) Hash(tx *types.Transaction) common.Hash {
	panic("can't sign with senderFromServer")
}
func (s *senderFromServer) SignatureValues(tx *types.Transaction, sig []byte) (R, S, V *big.Int, err error) {
	panic("can't sign with senderFromServer")
}
