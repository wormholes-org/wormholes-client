package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     []byte
	CodeHash []byte

	PledgedBalance     *big.Int
	PledgedBlockNumber *big.Int
	//Owner common.Address
	// whether the account has a NFT exchanger
	ExchangerFlag    bool
	BlockNumber      *big.Int
	ExchangerBalance *big.Int
	VoteBlockNumber  *big.Int
	VoteWeight       *big.Int
	Coefficient      uint8
	// The ratio that exchanger get.
	FeeRate       uint16
	ExchangerName string
	ExchangerURL  string
	// ApproveAddress have the right to handle all nfts of the account
	ApproveAddressList []common.Address
	// NFTBalance is the nft number that the account have
	//NFTBalance uint64
	// Indicates the reward method chosen by the miner
	//RewardFlag uint8 // 0:SNFT 1:ERB default:0
	AccountNFT
	Extra []byte
}
type AccountNFT struct {
	//Account
	Name   string
	Symbol string
	//Price                 *big.Int
	//Direction             uint8 // 0:un_tx,1:buy,2:sell
	Owner                 common.Address
	NFTApproveAddressList common.Address
	//Auctions map[string][]common.Address
	// MergeLevel is the level of NFT merged
	MergeLevel  uint8
	MergeNumber uint32
	//PledgedFlag           bool
	//NFTPledgedBlockNumber *big.Int

	Creator   common.Address
	Royalty   uint16
	Exchanger common.Address
	MetaURL   string
}

type BeneficiaryAddress struct {
	Address    common.Address
	NftAddress common.Address
}

type BeneficiaryAddressList []*BeneficiaryAddress

type ValidatorList struct {
	Validators []*Validator
}

type Validator struct {
	Addr    common.Address
	Balance *big.Int
	Proxy   common.Address
	Weight  []*big.Int
}

type MinerProxy struct {
	Address common.Address
	Proxy   common.Address
}

type MinerProxyList []*MinerProxy
