package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/wormholes-org/wormholes-client/tools"
	types2 "github.com/wormholes-org/wormholes-client/types"
	"log"
	"math/big"
)

type Wallet struct {
	priKey string
}

type Wormholes struct {
	Wallet
	c *rpc.Client
}

// NewClient creates a new wormclient for the given URL and priKey.
//when the rawurl is  nil, Initialize the wallet, can sign buyer, seller, exchange information.
//when the rawurl is not nil, Initialize the NFT, can carry out nft related transactions.
func NewClient(priKey, rawurl string) *Wormholes {
	if rawurl == "" {
		return &Wormholes{
			Wallet{priKey: priKey},
			nil,
		}
	} else {
		client, err := rpc.Dial(rawurl)
		if err != nil {
			log.Fatalf("failed to connect to Ethereum node: %v", err)
			return &Wormholes{}
		}
		return &Wormholes{
			Wallet{
				priKey: priKey,
			},
			client,
		}
	}
}

func (worm *Wormholes) CloseConnect() {
	worm.c.Close()
}

// ChainID retrieves the current chain ID for transaction replay protection.
func (worm *Wormholes) ChainID(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	err := worm.c.CallContext(ctx, &result, "eth_chainId")
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&result), err
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
//
// Note that loading full blocks requires two requests. Use HeaderByNumber
// if you don't need all transactions or uncle headers.
func (worm *Wormholes) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return worm.getBlock(ctx, "eth_getBlockByNumber", toBlockNumArg(number), true)
}

type rpcBlock struct {
	Hash         common.Hash      `json:"hash"`
	Transactions []rpcTransaction `json:"transactions"`
	UncleHashes  []common.Hash    `json:"uncles"`
}

func (worm *Wormholes) getBlock(ctx context.Context, method string, args ...interface{}) (*types.Block, error) {
	var raw json.RawMessage
	err := worm.c.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	} else if len(raw) == 0 {
		return nil, ethereum.NotFound
	}
	// Decode header and transactions.
	var head *types.Header
	var body rpcBlock
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	// Quick-verify transaction and uncle lists. This mostly helps with debugging the server.
	if head.UncleHash == types.EmptyUncleHash && len(body.UncleHashes) > 0 {
		return nil, fmt.Errorf("server returned non-empty uncle list but block header indicates no uncles")
	}
	if head.UncleHash != types.EmptyUncleHash && len(body.UncleHashes) == 0 {
		return nil, fmt.Errorf("server returned empty uncle list but block header indicates uncles")
	}
	if head.TxHash == types.EmptyRootHash && len(body.Transactions) > 0 {
		return nil, fmt.Errorf("server returned non-empty transaction list but block header indicates no transactions")
	}
	if head.TxHash != types.EmptyRootHash && len(body.Transactions) == 0 {
		return nil, fmt.Errorf("server returned empty transaction list but block header indicates transactions")
	}
	// Load uncles because they are not included in the block response.
	var uncles []*types.Header
	if len(body.UncleHashes) > 0 {
		uncles = make([]*types.Header, len(body.UncleHashes))
		reqs := make([]rpc.BatchElem, len(body.UncleHashes))
		for i := range reqs {
			reqs[i] = rpc.BatchElem{
				Method: "eth_getUncleByBlockHashAndIndex",
				Args:   []interface{}{body.Hash, hexutil.EncodeUint64(uint64(i))},
				Result: &uncles[i],
			}
		}
		if err := worm.c.BatchCallContext(ctx, reqs); err != nil {
			return nil, err
		}
		for i := range reqs {
			if reqs[i].Error != nil {
				return nil, reqs[i].Error
			}
			if uncles[i] == nil {
				return nil, fmt.Errorf("got null header for uncle %d of block %x", i, body.Hash[:])
			}
		}
	}
	// Fill the sender cache of transactions in the block.
	txs := make([]*types.Transaction, len(body.Transactions))
	for i, tx := range body.Transactions {
		if tx.From != nil {
			setSenderFromServer(tx.tx, *tx.From, body.Hash)
		}
		txs[i] = tx.tx
	}
	return types.NewBlockWithHeader(head).WithBody(txs, uncles), nil
}

// BlockNumber returns the most recent block number
func (worm *Wormholes) BlockNumber(ctx context.Context) (uint64, error) {
	var result hexutil.Uint64
	err := worm.c.CallContext(ctx, &result, "eth_blockNumber")
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

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

// TransactionInBlock returns a single transaction at index in the given block.
func (worm *Wormholes) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	var json *rpcTransaction
	err := worm.c.CallContext(ctx, &json, "eth_getTransactionByBlockHashAndIndex", blockHash, hexutil.Uint64(index))
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
func (worm *Wormholes) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var result hexutil.Uint64
	err := worm.c.CallContext(ctx, &result, "eth_getTransactionCount", account, "pending")
	return uint64(result), err
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (worm *Wormholes) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var hex hexutil.Big
	if err := worm.c.CallContext(ctx, &hex, "eth_gasPrice"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

// SendTransaction injects a signed transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the TransactionReceipt method to get the
// contract address after the transaction has been mined.
func (worm *Wormholes) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	return worm.c.CallContext(ctx, nil, "eth_sendRawTransaction", hexutil.Encode(data))
}

// NetworkID returns the network ID (also known as the chain ID) for this chain.
func (worm *Wormholes) NetworkID(ctx context.Context) (*big.Int, error) {
	version := new(big.Int)
	var ver string
	if err := worm.c.CallContext(ctx, &ver, "net_version"); err != nil {
		return nil, err
	}
	if _, ok := version.SetString(ver, 10); !ok {
		return nil, fmt.Errorf("invalid net_version result %q", ver)
	}
	return version, nil
}

// Balance returns the wei balance of the given account in the pending state.
func (worm *Wormholes) Balance(ctx context.Context, account string) (*big.Int, error) {
	var accounts common.Address
	accounts = common.HexToAddress(account)
	var result hexutil.Big
	err := worm.c.CallContext(ctx, &result, "eth_getBalance", accounts, "pending")
	return (*big.Int)(&result), err
}

// BalanceAt returns the wei balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (worm *Wormholes) BalanceAt(ctx context.Context, account string, blockNumber *big.Int) (*big.Int, error) {
	var accounts common.Address
	accounts = common.HexToAddress(account)
	var result hexutil.Big
	err := worm.c.CallContext(ctx, &result, "eth_getBalance", accounts, toBlockNumArg(blockNumber))
	return (*big.Int)(&result), err
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (worm *Wormholes) TransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	txHashs := common.HexToHash(txHash)
	var r *types.Receipt
	err := worm.c.CallContext(ctx, &r, "eth_getTransactionReceipt", txHashs)
	if err == nil {
		if r == nil {
			return nil, ethereum.NotFound
		}
	}
	return r, err
}

func (worm *Wormholes) GetAccountInfo(ctx context.Context, address string, block int64) (*types2.Account, error) {
	var addresss common.Address
	addresss = common.HexToAddress(address)
	blockNrOrHash := rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(block))
	var r *types2.Account
	err := worm.c.CallContext(ctx, &r, "eth_getAccountInfo", addresss, blockNrOrHash)
	if err == nil {
		if r == nil {
			return nil, ethereum.NotFound
		}
	}
	return r, err
}

func (worm *Wormholes) GetStakers(ctx context.Context, number uint64) (*types2.ValidatorList, error) {
	var al *types2.ValidatorList
	nu := rpc.BlockNumber(number)
	err := worm.c.CallContext(ctx, &al, "eth_getValidator", nu)
	if err != nil {
		return nil, err
	}
	return al, err
}

func (worm *Wormholes) GetValidators(ctx context.Context, block int64) (*types2.ValidatorList, error) {
	blocNumber := rpc.BlockNumber(block)
	var r *types2.ValidatorList
	err := worm.c.CallContext(ctx, &r, "eth_getValidator", blocNumber)
	if err == nil {
		if r == nil {
			return nil, ethereum.NotFound
		}
	}
	return r, nil
}

func (worm *Wormholes) GetBlockBeneficiaryAddressByNumber(ctx context.Context, block int64) (*types2.BeneficiaryAddressList, error) {
	blockNumber := rpc.BlockNumber(block)
	var r *types2.BeneficiaryAddressList
	err := worm.c.CallContext(ctx, &r, "eth_getBlockBeneficiaryAddressByNumber", blockNumber, true)
	if err == nil {
		if r == nil {
			return nil, ethereum.NotFound
		}
	}
	return r, err
}

func (worm *Wormholes) QueryMinerProxy(ctx context.Context, number int64, account string) (types2.MinerProxyList, error) {
	var result types2.MinerProxyList
	nu := fmt.Sprintf("0x%x", number)
	var accounts common.Address

	accounts = common.HexToAddress(account)

	err := worm.c.CallContext(ctx, &result, "eth_queryMinerProxy", nu, accounts)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (w *Wallet) Sign(data []byte, priKey string) ([]byte, error) {
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
//amount: The amount the buyer purchased the NFT, formatted as a hexadecimal string
//nftAddress: The NFT address of the transaction. The format is a hexadecimal string. When this field is filled in, it means that the transaction has minted nft. When not filled, it means lazy transaction, and the nft has not been minted
//exchanger: The exchange on which the transaction took place, formatted as a decimal string
//blockNumber: Block height, which means that this transaction is valid before this height, the format is a hexadecimal string
//seller: Seller's address, formatted as a hexadecimal string
func (w *Wallet) SignBuyer(amount, nftAddress, exchanger, blockNumber, seller string) ([]byte, error) {
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
// Signed Mint Seller
//	amount: The amount the buyer purchased the NFT, formatted as a hexadecimal string
//	nftAddress: The NFT address of the transaction, formatted as a hexadecimal string
//	exchanger:	The exchange on which the transaction took place, formatted as a decimal string
//	blockNumber: Block height, which means that this transaction is valid before this height, the format is a hexadecimal string
func (w *Wallet) SignSeller1(amount, nftAddress, exchanger, blockNumber string) ([]byte, error) {
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
// Signed Unminted Seller
//	amount: The amount of the NFT transaction, formatted as a hexadecimal string
//	royalty: royalty, hex string
//	metaURL: NFT metadata address
//	exclusiveFlag: "0": Inclusive, "1": Exclusive
//	exchanger:	The exchange on which the transaction took place, formatted as a decimal string
//	blockNumber: Block height, which means that this transaction is valid before this height, the format is a hexadecimal string
func (w *Wallet) SignSeller2(amount, royalty, metaURL, exclusiveFlag, exchanger, blockNumber string) ([]byte, error) {
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
// Signed by an authorized exchange
//	exchangerOwner: Authorize exchange, formatted as a hexadecimal string
// 	to: Authorized exchange, formatted as a hexadecimal string
//	block_number: Block height, which means that this transaction is valid before this height, the format is a hexadecimal string
func (w *Wallet) SignExchanger(exchangerOwner, to, blockNumber string) ([]byte, error) {
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

func (w *Wallet) SignDelegate(address, pledgeAcoount string) ([]byte, error) {
	key, err := crypto.HexToECDSA(w.priKey)
	if err != nil {
		return nil, err
	}

	msg := address + pledgeAcoount
	signature, err := crypto.Sign(tools.SignHash([]byte(msg)), key)
	if err != nil {
		return nil, err
	}

	signature[64] += 27
	return []byte(hexutil.Encode(signature)), nil
}
