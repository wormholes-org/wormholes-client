package client

import (
	"context"
	"encoding/json"
	"github.com/wormholesclient/tools"
	types2 "github.com/wormholesclient/types"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
)

type WormClient struct {
	NFT    NFT
	Wallet Wallet
}

type NFT struct {
	c      *rpc.Client
	priKey string
}

type Wallet struct {
	priKey string
}

//NewClient creates a new NFT for the given URL and priKey.
//when the rawurl is  nil, Initialize the wallet, can sign buyer, seller, exchange information.
//when the rawurl is not nil, Initialize the NFT, can carry out nft related transactions.
func NewClient(priKey, rawurl string) *WormClient {
	if rawurl == "" {
		return &WormClient{
			Wallet: Wallet{
				priKey: priKey,
			},
		}
	} else {
		client, err := rpc.Dial(rawurl)
		if err != nil {
			log.Fatalf("failed to connect to Ethereum node: %v", err)
			return nil
		}
		return &WormClient{
			NFT: NFT{
				c:      client,
				priKey: priKey,
			},
			Wallet: Wallet{
				priKey: priKey,
			},
		}
	}
}

func (nft *NFT) CloseConnect() {
	nft.c.Close()
}

// ChainID retrieves the current chain ID for transaction replay protection.
func (nft *NFT) ChainID(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	err := nft.c.CallContext(ctx, &result, "eth_chainId")
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&result), err
}

// BlockNumber returns the most recent block number
func (nft *NFT) BlockNumber(ctx context.Context) (uint64, error) {
	var result hexutil.Uint64
	err := nft.c.CallContext(ctx, &result, "eth_blockNumber")
	return uint64(result), err
}

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

// TransactionInBlock returns a single transaction at index in the given block.
func (nft *NFT) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	var json *rpcTransaction
	err := nft.c.CallContext(ctx, &json, "eth_getTransactionByBlockHashAndIndex", blockHash, hexutil.Uint64(index))
	if err != nil {
		return nil, err
	}
	if json == nil {
		return nil, ethereum.NotFound
	} else if _, r, _ := json.tx.RawSignatureValues(); r == nil {
		return nil, fmt.Errorf("server returned transaction without signature")
	}
	if json.From != nil && json.BlockHash != nil {
		setSenderFromServer(json.tx, *json.From, *json.BlockHash)
	}
	return json.tx, err
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (nft *NFT) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var result hexutil.Uint64
	err := nft.c.CallContext(ctx, &result, "eth_getTransactionCount", account, "pending")
	return uint64(result), err
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (nft *NFT) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var hex hexutil.Big
	if err := nft.c.CallContext(ctx, &hex, "eth_gasPrice"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

// SendTransaction injects a signed transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the TransactionReceipt method to get the
// contract address after the transaction has been mined.
func (nft *NFT) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	return nft.c.CallContext(ctx, nil, "eth_sendRawTransaction", hexutil.Encode(data))
}

// NetworkID returns the network ID (also known as the chain ID) for this chain.
func (nft *NFT) NetworkID(ctx context.Context) (*big.Int, error) {
	version := new(big.Int)
	var ver string
	if err := nft.c.CallContext(ctx, &ver, "net_version"); err != nil {
		return nil, err
	}
	if _, ok := version.SetString(ver, 10); !ok {
		return nil, fmt.Errorf("invalid net_version result %q", ver)
	}
	return version, nil
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (nft *NFT) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	var r *types.Receipt
	err := nft.c.CallContext(ctx, &r, "eth_getTransactionReceipt", txHash)
	if err == nil {
		if r == nil {
			return nil, ethereum.NotFound
		}
	}
	return r, err
}

func (nft *NFT) GetAccountInfo(ctx context.Context, address common.Address, block int64) (*types2.Account, error) {
	blockNrOrHash := rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(block))
	var r *types2.Account
	err := nft.c.CallContext(ctx, &r, "eth_getAccountInfo", address, blockNrOrHash)
	if err == nil {
		if r == nil {
			return nil, ethereum.NotFound
		}
	}

	return r, err
}

func (nft *NFT) GetBlockBeneficiaryAddressByNumber(ctx context.Context, block int64) (*types2.BeneficiaryAddressList, error) {
	blockNumber := rpc.BlockNumber(block)
	var r *types2.BeneficiaryAddressList
	err := nft.c.CallContext(ctx, &r, "eth_getBlockBeneficiaryAddressByNumber", blockNumber, true)
	if err == nil {
		if r == nil {
			return nil, ethereum.NotFound
		}
	}

	return r, err
}

func (w *Wallet) sign(data []byte, priKey string) ([]byte, error) {
	//am := core.StartClefAccountManager("/home/user1/azh/data/node15/keystore", true, false, "") //获取account
	key, err := crypto.HexToECDSA(priKey)
	if err != nil {
		return nil, err
	}

	signature, err := crypto.Sign(tools.SignHash(data), key)
	if err != nil {
		return nil, err
	}

	signature[64] += 27

	return signature, nil
}

// SignBuyer
//amount: 买家购买NFT的金额， 格式为十六进制字符串
//nftAddress: 交易的NFT地址， 格式为十六进制字符串，填写该字段时，表示交易已铸造的nft，不填写时，表示惰性交易，nft未铸造
//exchanger: 交易发生所在的交易所， 格式为十进制字符串
//blockNumber: 块高度， 代表此交易在此高度之前发生才有效， 格式为十六进制字符串
//seller: 卖家地址， 格式是十六进制字符串
func (w *Wallet) SignBuyer(amount, nftAddress, exchanger, blockNumber, seller string) ([]byte, error) {
	//am := core.StartClefAccountManager("/home/user1/azh/data/node15/keystore", true, false, "") //获取account
	key, err := crypto.HexToECDSA(w.priKey)
	if err != nil {
		return nil, err
	}

	msg := amount + nftAddress + exchanger + blockNumber + seller
	signature, err := crypto.Sign(tools.SignHash([]byte(msg)), key)
	if err != nil {
		return nil, err
	}

	signature[64] += 27

	buyer := types2.Buyer{
		Amount:      amount,
		NFTAddress:  nftAddress,
		Exchanger:   exchanger,
		BlockNumber: blockNumber,
		Seller:      seller,
		Sig:         hexutil.Encode(signature),
	}

	result, err := json.Marshal(buyer)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SignSeller1
// 签名已铸币卖家
//	amount: nft交易金额，十六进制字符串
//	nftAddress: 交易的nft地址，十六进制字符串
//	exchanger:	交易发生所在的交易所， 格式为十进制字符串
//	blockNumber: 块高度， 代表此交易在此高度之前发生才有效， 格式为十六进制字符串
func (w *Wallet) SignSeller1(amount, nftAddress, exchanger, blockNumber string) ([]byte, error) {
	//am := core.StartClefAccountManager("/home/user1/azh/data/node15/keystore", true, false, "") //获取account
	key, err := crypto.HexToECDSA(w.priKey)
	if err != nil {
		return nil, err
	}

	msg := amount + nftAddress + exchanger + blockNumber
	signature, err := crypto.Sign(tools.SignHash([]byte(msg)), key)
	if err != nil {
		return nil, err
	}

	signature[64] += 27

	seller1 := types2.Seller1{
		Amount:      amount,
		NFTAddress:  nftAddress,
		Exchanger:   exchanger,
		BlockNumber: blockNumber,
		Sig:         hexutil.Encode(signature),
	}

	result, err := json.Marshal(seller1)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SignSeller2
// 签名未铸币卖家
//	amount: nft交易金额，十六进制字符串
//	royalty: 版税，十六进制字符串
//	metaURL: NFT元数据地址
//	exclusiveFlag: ”0”:非独占,”1”:独占
//	exchanger:	铸造NFT时的交易所，格式为字符串，填写该字段时，该交易所独占此NFT，不填写时，没有交易所独占此NFT
//	blockNumber: 块高度， 代表此交易在此高度之前发生才有效， 格式为十六进制字符串
func (w *Wallet) SignSeller2(amount, royalty, metaURL, exclusiveFlag, exchanger, blockNumber string) ([]byte, error) {
	//am := core.StartClefAccountManager("/home/user1/azh/data/node15/keystore", true, false, "") //获取account
	key, err := crypto.HexToECDSA(w.priKey)
	if err != nil {
		return nil, err
	}

	msg := amount + royalty + metaURL + exclusiveFlag + exchanger + blockNumber
	signature, err := crypto.Sign(tools.SignHash([]byte(msg)), key)
	if err != nil {
		return nil, err
	}

	signature[64] += 27

	seller2 := types2.Seller2{
		Amount:        amount,
		Royalty:       royalty,
		MetaURL:       metaURL,
		ExclusiveFlag: exclusiveFlag,
		Exchanger:     exchanger,
		BlockNumber:   blockNumber,
		Sig:           hexutil.Encode(signature),
	}

	result, err := json.Marshal(seller2)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SignExchanger
// 被授权交易所签名
//	exchangerOwner: 授权交易所， 格式为十六进制字符串
// 	to: 被授权交易所， 格式为十六进制字符串
//	block_number: 块高度， 代表授权在此高度之前有效， 格式为十六进制字符串
func (w *Wallet) SignExchanger(exchangerOwner, to, blockNumber string) ([]byte, error) {
	//am := core.StartClefAccountManager("/home/user1/azh/data/node15/keystore", true, false, "") //获取account
	key, err := crypto.HexToECDSA(w.priKey)
	if err != nil {
		return nil, err
	}

	msg := exchangerOwner + to + blockNumber
	signature, err := crypto.Sign(tools.SignHash([]byte(msg)), key)
	if err != nil {
		return nil, err
	}

	signature[64] += 27

	exchangeAuth := types2.ExchangerAuth{
		ExchangerOwner: exchangerOwner,
		To:             to,
		BlockNumber:    blockNumber,
		Sig:            hexutil.Encode(signature),
	}

	result, err := json.Marshal(exchangeAuth)
	if err != nil {
		return nil, err
	}
	return result, nil
}
