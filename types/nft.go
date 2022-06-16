package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Account struct {
	Nonce   uint64
	Balance *big.Int
	// *** modify to support nft transaction 20211220 begin ***
	//NFTCount uint64		// number of nft who account have
	// *** modify to support nft transaction 20211220 end ***
	Root           common.Hash // merkle root of the storage trie
	CodeHash       []byte
	PledgedBalance *big.Int
	// *** modify to support nft transaction 20211215 ***
	//Owner common.Address
	// whether the account has a NFT exchanger
	ExchangerFlag    bool
	BlockNumber      *big.Int
	ExchangerBalance *big.Int
	VoteWeight       *big.Int
	// The ratio that exchanger get.
	FeeRate       uint32
	ExchangerName string
	ExchangerURL  string
	// ApproveAddress have the right to handle all nfts of the account
	ApproveAddressList []common.Address
	// NFTBalance is the nft number that the account have
	NFTBalance uint64
	AccountNFT
}

type AccountNFT struct {
	//Account
	Name                  string
	Symbol                string
	Price                 *big.Int
	Direction             uint8 // 0:未交易,1:买入,2:卖出
	Owner                 common.Address
	NFTApproveAddressList common.Address
	//Auctions map[string][]common.Address
	// MergeLevel is the level of NFT merged
	MergeLevel uint8

	Creator   common.Address
	Royalty   uint32
	Exchanger common.Address
	MetaURL   string
}

type BeneficiaryAddress struct {
	Address    common.Address
	NftAddress common.Address
}

type BeneficiaryAddressList []*BeneficiaryAddress
