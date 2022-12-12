package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/wormholes-org/wormholes-client/tools"
	types2 "github.com/wormholes-org/wormholes-client/types"
	"golang.org/x/xerrors"
	"log"
	"math/big"
	"strings"
)

// NormalTransaction
//	Parameter Description
//  to 			Account address
//  value		transaction amount
//  data
func (worm *Wormholes) NormalTransaction(to string, value int64, data string) (string, error) {
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("NormalTransaction() priKeyToAddress err ", err)
		return "", err
	}

	toAddr := common.HexToAddress(to)
	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("NormalTransaction() suggestGasPrice err ", err)
		return "", err
	}

	wei, _ := new(big.Int).SetString("1000000000000000000", 10)
	charge := new(big.Int).Mul(big.NewInt(value), wei)
	tx := types.NewTransaction(nonce, toAddr, charge, gasLimit, gasPrice, []byte(data))
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("NormalTransaction() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("NormalTransaction() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("NormalTransaction() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// Mint NFT user minting
//	Users can use this transaction to create an NFT on the wormholes chain
//
//	Parameter Description
//	royalty: 10,																					Royalty, formatted as an integer
//	metaURL: "/ipfs/ddfd90be9408b4",	NFT metadata address
//	exchanger:"0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4",							The exchange when the NFT is minted, the format is a string. When this field is filled, the exchange will exclusively own the NFT. If it is not filled in, no exchange will exclusively own the NFT
func (worm *Wormholes) Mint(royalty uint32, metaURL string, exchanger string) (string, error) {
	if exchanger != "" {
		err := tools.CheckAddress("Mint() exchanger", exchanger)
		if err != nil {
			return "", err
		}
	}

	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		return "", err
	}

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(60000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("Mint() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:      types2.Mint,
		Royalty:   royalty,
		MetaURL:   metaURL,
		Exchanger: exchanger,
		Version:   types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("Mint() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)

	tx := types.NewTransaction(nonce, account, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("Mint() networkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("Mint() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("Mint() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// Transfer NFT transfer
// 	Change ownership of NFTs
//
//	Parameter Description
//	wormAddress: "0x8000000000000000000000000000000000000001",  worm address, the format is a decimal string, when it is SNFT, the length can be less than 42 (including 0x), representing the synthesized SNFT
//	to:         "0x814920c33b1a037F91a16B126282155c6F92A10F",  Target NFT user address
func (worm *Wormholes) Transfer(wormAddress, to string) (string, error) {
	err := tools.CheckHex("Transfer() wormAddress", wormAddress)
	if err != nil {
		return "", err
	}
	err = tools.CheckAddress("Transfer() to", to)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		return "", err
	}

	toAddr := common.HexToAddress(to)

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(50000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("Transfer() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:       types2.Transfer,
		NFTAddress: wormAddress,
		Version:    types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("Transfer() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)

	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, toAddr, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("Transfer() networkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("Transfer() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("Transfer() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// Author Authorize an NFT to an exchange
//
//	Parameter Description
//	wormAddress: "0x0000000000000000000000000000000000000001",	Authorized worm address, the format is a decimal string, when it is SNFT, the length can be less than 42 (including 0x), representing the synthesized SNFT
//	to:         "0x814920c33b1a037F91a16B126282155c6F92A10F",	Licensee's address
func (worm *Wormholes) Author(wormAddress, to string) (string, error) {
	err := tools.CheckHex("Author() wormAddress", wormAddress)
	if err != nil {
		return "", err
	}
	err = tools.CheckAddress("Author() to", to)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		return "", err
	}

	toAddr := common.HexToAddress(to)

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(50000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("Author() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:       types2.Author,
		NFTAddress: wormAddress,
		Version:    types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("Author failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)

	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, toAddr, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("Author() networkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("Author signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("Author sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// AuthorRevoke Cancel the authorization of an NFT
//
//	Parameter Description
//	wormAddress: "0x0000000000000000000000000000000000000002",	Authorized worm address, the format is a decimal string, when it is SNFT, the length can be less than 42 (including 0x), representing the synthesized SNFT
//	to:         "0x814920c33b1a037F91a16B126282155c6F92A10F",	Licensee's address
func (worm *Wormholes) AuthorRevoke(wormAddress, to string) (string, error) {
	err := tools.CheckHex("AuthorRevoke() wormAddress", wormAddress)
	if err != nil {
		return "", err
	}
	err = tools.CheckAddress("AuthorRevoke() to", to)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		return "", err
	}

	toAddr := common.HexToAddress(to)

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(50000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("AuthorRevoke suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:       types2.AuthorRevoke,
		NFTAddress: wormAddress,
		Version:    types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("AuthorRevoke() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, toAddr, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("AuthorRevoke() networkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("AuthorRevoke() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("AuthorRevoke() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// AccountAuthor
//	Authorize all NFTs under an account to the exchange
//	Parameter Description
//	to:     "0x814920c33b1a037F91a16B126282155c6F92A10F",							Licensee's address
func (worm *Wormholes) AccountAuthor(to string) (string, error) {
	err := tools.CheckAddress("AccountAuthor() to", to)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		return "", err
	}

	toAddr := common.HexToAddress(to)

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(50000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("AccountAuthor() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.AccountAuthor,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("AccountAuthor() ailed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, toAddr, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("AccountAuthor() networkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("AccountAuthor() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("AccountAuthor sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// AccountAuthorRevoke
//	Cancel all NFT authorizations under an account
//
//	Parameter Description
//	to:     "0x814920c33b1a037F91a16B126282155c6F92A10F",							Licensee's address
func (worm *Wormholes) AccountAuthorRevoke(to string) (string, error) {
	err := tools.CheckAddress("AccountAuthorRevoke() to", to)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		return "", err
	}

	toAddr := common.HexToAddress(to)

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(50000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("AccountAuthorRevoke() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.AccountAuthorRevoke,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("AccountAuthorRevoke() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, toAddr, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("AccountAuthorRevoke() networkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("AccountAuthorRevoke() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("AccountAuthorRevoke() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// SNFTToERB
//	Convert NFT fragments mined by miners to ERB
//
//	Parameter Description
//	wormAddress: "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",	The converted sworm address, in the format of a decimal string, the length can be less than 42 (including 0x), representing the synthesized SNFT
//
//	The exchange price corresponding to the synthesis level
//	0: 100000000000000000
// 	1: 150000000000000000
//	2: 225000000000000000
//	3: 300000000000000000
func (worm *Wormholes) SNFTToERB(wormAddress string) (string, error) {
	err := tools.CheckHex("SNFTToERB() wormAddress", wormAddress)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		return "", err
	}

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(50000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("SNFTToERB() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:       types2.SNFTToERB,
		NFTAddress: wormAddress,
		Version:    types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("SNFTToERB() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, account, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("SNFTToERB() networkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SNFTToERB() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SNFTToERB() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// TokenPledge
//	When a user wants to become a miner, he needs to do an ERB pledge transaction first to pledge the ERB needed to become a miner
func (worm *Wormholes) TokenPledge(proxySign []byte, proxyAddress string) (string, error) {
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("TokenPledge() priKeyToAddress err ", err)
		return "", err
	}

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(70000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("TokenPledge() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:         types2.TokenPledge,
		ProxyAddress: proxyAddress,
		ProxySign:    string(proxySign),
		Version:      types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("TokenPledge() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	wei, _ := new(big.Int).SetString("1000000000000000000", 10)
	pledge := new(big.Int).Mul(big.NewInt(100000), wei)
	tx := types.NewTransaction(nonce, account, pledge, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("TokenPledge() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("TokenPledge() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("TokenPledge() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// TokenRevokesPledge
//	When the user does not want to be a miner, or no longer wants to pledge so much ERB, he can do ERB to revoke the pledge
func (worm *Wormholes) TokenRevokesPledge() (string, error) {
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("TokenRevokesPledge() priKeyToAddress err ", err)
		return "", err
	}

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(50000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("TokenRevokesPledge() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.TokenRevokesPledge,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("TokenRevokesPledge() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	wei, _ := new(big.Int).SetString("1000000000000000000", 10)
	pledge := new(big.Int).Mul(big.NewInt(100000), wei)

	tx := types.NewTransaction(nonce, account, pledge, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("TokenRevokesPledge() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("TokenRevokesPledge() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("TokenRevokesPledge() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// Open
//	This transaction can be initiated when a user wants to open an exchange
//
//	Parameter Description
//	feeRate:   10,																		  The exchange rate, the format is an integer type
//	name:      "wormholes",										 Exchange name, formatted as a string
//	url:       "www.kang123456.com",		Exchange server address, formatted as a string
func (worm *Wormholes) Open(feeRate uint32, name, url string) (string, error) {
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("Open() priKeyToAddress err ", err)
		return "", err
	}

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(60000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("Open() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.Open,
		FeeRate: feeRate,
		Name:    name,
		Url:     url,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("Open() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	wei, _ := new(big.Int).SetString("1000000000000000000", 10)
	amount := new(big.Int).Mul(big.NewInt(100), wei)

	tx := types.NewTransaction(nonce, account, amount, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("open() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("open() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("open() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// Close
//	When the user does not want to continue to open an exchange, he can initiate this transaction to close the opened exchange
func (worm *Wormholes) Close() (string, error) {
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("close() priKeyToAddress err ", err)
		return "", err
	}

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(60000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("close() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.Close,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("close() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, account, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("close networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("close() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("close() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// InsertNFTBlock
// Deprecated: use VoteOfficialNFT instead.
//	This transaction is used to inject NFT fragments that can be mined by miners. Only official accounts can do this transaction
//
//	Parameter Descriptiom
//	dir:        "wormholes",  													The path address where sworm is located, the format is a string
//	startIndex: "0x640001",	 														The start number of the sworm fragment, formatted as a hexadecimal string
//	number:     6553600,														The number of sworm shards injected, formatted as a decimal string
//	royalty:    20,																			Royalty, formatted as an integer
//	creator:    "0xab7624f47fd7dadb6b8e255d06a2f10af55990fe",	creator, format is a hex string
func (worm *Wormholes) InsertNFTBlock(dir, startIndex string, number uint64, royalty uint32, creator string) (string, error) {
	err := tools.CheckAddress("InsertNFTBlock() creator", creator)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("InsertNFTBlock() priKeyToAddress err ", err)
		return "", err
	}

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("InsertNFTBlock() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:       types2.InsertNFTBlock,
		Dir:        dir,
		StartIndex: startIndex,
		Number:     number,
		Royalty:    royalty,
		Creator:    creator,
		Version:    types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("InsertNFTBlock() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, account, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("InsertNFTBlock() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("InsertNFTBlock() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("InsertNFTBlock() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// TransactionNFT
//	For buying and selling NFTs that have been minted, the transaction originator can be an exchange or a seller
//
//	Parameter Description
//	buyer: { "price":"0xde0b6b3a7640000", "worm_address":"0x0000000000000000000000000000000000000002", "exchanger":"0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4", "block_number":"0x487", "sig":"0x24355436e991443b8ed3fb83e8c2fa02f8e2bfc0f716c320f836ee7d756e3c712e7e2510b994d1cb7be85d6643233abc81c23929ce7c1c1effd93db261aac5211b" }																				买家
//	to:     "0x5051B76579BC966A9480dd6E72B39A4C89c1154C",				Buyer's address
func (worm *Wormholes) TransactionNFT(buyer []byte, to string) (string, error) {
	err := tools.CheckAddress("TransactionNFT() to", to)
	if err != nil {
		return "", err
	}

	var buyers types2.Buyer
	err = json.Unmarshal(buyer, &buyers)
	if err != nil {
		return "", xerrors.New("the formate of buyer is wrong")
	}

	err = tools.CheckHex("buyers.BlockNumber", buyers.BlockNumber)
	if err != nil {
		return "", err
	}

	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("TransactionNFT() priKeyToAddress err ", err)
		return "", err
	}

	ctx := context.Background()
	nonce, err := worm.PendingNonceAt(ctx, account)

	toAddr := common.HexToAddress(to)

	gasLimit := uint64(100000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("TransactionNFT() suggestGasPrice err ", err)
		return "", err
	}

	//msg := buyer.Amount + buyer.NFTAddress + buyer.Exchanger + buyer.BlockNumber
	//
	//sig, err := tools.Sign([]byte(msg), wormholes.Buyer.PriKey)
	//if err != nil {
	//	return "", err
	//}
	//buyer.Sig = hexutil.Encode(sig)

	transaction := types2.Transaction{
		Type:    types2.TransactionNFT,
		Buyer:   &buyers,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("TransactionNFT() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(buyers.Amount)
	fmt.Println(value)
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("TransactionNFT() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("TransactionNFT() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("TransactionNFT sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// BuyerInitiatingTransaction
//	Used to buy and sell NFTs that have been minted, the transaction initiator is the buyer
//
//	Parameter Description
//	seller1: { "price":"0x38D7EA4C68000", "worm_address":"0x0000000000000000000000000000000000000003", "exchanger":"0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4", "block_number":"0x65d", "sig":"0x94e88fb5686551dfc3006c608423983a248df8502cbbcaeb2c3352f267a25e531d5fc745bea5f7f564b7399fb70d87026bbf9952f1403e9d4dae4aa14b091cff1c" }
func (worm *Wormholes) BuyerInitiatingTransaction(seller1 []byte) (string, error) {
	var seller1s types2.Seller1
	err := json.Unmarshal(seller1, &seller1s)
	if err != nil {
		return "", xerrors.New("the formate of seller1 is wrong")
	}

	err = tools.CheckHex("seller1s.BlockNumber", seller1s.BlockNumber)
	if err != nil {
		return "", err
	}
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("BuyerInitiatingTransaction() priKeyToAddress err ", err)
		return "", err
	}

	ctx := context.Background()
	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(100000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("BuyerInitiatingTransaction() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.BuyerInitiatingTransaction,
		Seller1: &seller1s,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("BuyerInitiatingTransaction() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(seller1s.Amount)
	tx := types.NewTransaction(nonce, account, value, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("BuyerInitiatingTransaction networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("BuyerInitiatingTransaction signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("BuyerInitiatingTransaction sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// FoundryTradeBuyer
//	For buying and selling unminted NFTs, the transaction originator is the buyer
//
//	Parameter Description
//	seller2: { "price":"0x38D7EA4C68000", "royalty":"0xa", "meta_url":"/ipfs/qqqqqqqqqq", "exclusive_flag":"0", "exchanger":"0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4", "block_number":"0x703", "sig":"0xb08cf8b2f2d4b2635a85d1c7a816f01c24ac2a90ab49bdbe0e52e0a8f07eea5521eb80554df2c403423bdf49f412a7811b10a16005832a1bc171f5dfd3c983121c" }
func (worm *Wormholes) FoundryTradeBuyer(seller2 []byte) (string, error) {
	var seller2s types2.Seller2
	err := json.Unmarshal([]byte(seller2), &seller2s)
	if err != nil {
		return "", xerrors.New("the formate of seller2 is wrong")
	}

	err = tools.CheckFlag("seller2s.ExclusiveFlag", seller2s.ExclusiveFlag)
	if err != nil {
		return "", err
	}

	err = tools.CheckHex("seller2s.BlockNumber", seller2s.BlockNumber)
	if err != nil {
		return "", err
	}

	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("FoundryTradeBuyer() priKeyToAddress err ", err)
		return "", err
	}

	ctx := context.Background()
	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(101000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("FoundryTradeBuyer() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.FoundryTradeBuyer,
		Seller2: &seller2s,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("FoundryTradeBuyer() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(seller2s.Amount)
	tx := types.NewTransaction(nonce, account, value, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("FoundryTradeBuyer() failed to format wormholes dataNetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("FoundryTradeBuyer() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("FoundryTradeBuyer() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// FoundryExchange
//	For buying and selling unminted NFTs, the transaction originator is the exchange, or the seller
//
//	Parameter Description
//	buyer:   {"price":"0xde0b6b3a7640000","exchanger":"0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4","block_number":"0x7c6","sig":"0xd4d2319bd9c4c1664ceb8cdb4d417fc22a6b4083845d5390154f4d268b07bc81755b0f728f989554142ca8124fe543b93a526f92664d7cc905ec361721ef130a1b"}
//	seller2: {"price":"0x38D7EA4C68000","royalty":"0xa","meta_url":"/ipfs/qqqqqqqqqq","exclusive_flag":"0","exchanger":"0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4","block_number":"0x7be","sig":"0x84c0c293298557e38fa5064a6fb3b9e6930fa46b234fcd0a923cd677369f5aad3f014a164b21077f713e25b4e986673f614f6ce824561fbda2b4e67e018fac6f1b"}
//	to:      "0x5051B76579BC966A9480dd6E72B39A4C89c1154C",  Buyer's address
func (worm *Wormholes) FoundryExchange(buyer, seller2 []byte, to string) (string, error) {
	err := tools.CheckAddress("to", to)
	if err != nil {
		return "", err
	}

	var buyers types2.Buyer
	err = json.Unmarshal(buyer, &buyers)
	if err != nil {
		return "", xerrors.New("the formate of buyer is wrong")
	}

	err = tools.CheckHex("buyers.BlockNumber", buyers.BlockNumber)
	if err != nil {
		return "", err
	}

	var seller2s types2.Seller2
	err = json.Unmarshal(seller2, &seller2s)
	if err != nil {
		return "", xerrors.New("the formate of seller2 is wrong")
	}

	err = tools.CheckFlag("seller2s.ExclusiveFlag", seller2s.ExclusiveFlag)
	if err != nil {
		return "", err
	}

	err = tools.CheckHex("seller2s.BlockNumber", seller2s.BlockNumber)
	if err != nil {
		return "", err
	}

	if buyers.Amount < seller2s.Amount {
		return "", xerrors.New("buyer`s amount must be greater then seller`s amount")
	}
	if seller2s.Exchanger != buyers.Exchanger {
		return "", xerrors.New("buyer`s exchanger and seller`s exchanger and transaction`s exchanger aren`t same")
	}

	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("FoundryExchange() priKeyToAddress err ", err)
		return "", err
	}

	toAddr := common.HexToAddress(to)

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(140000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("FoundryExchange() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.FoundryExchange,
		Buyer:   &buyers,
		Seller2: &seller2s,
		Version: types2.WormHolesVersion,
	}
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("FoundryExchange() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(buyers.Amount)
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("FoundryExchange() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("FoundryExchange() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("FoundryExchange() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// NftExchangeMatch
//	It is used to buy and sell NFTs that have been minted. The transaction originator is the exchange. This transaction is used when exchange A authorizes another exchange B, and exchange B initiates the transaction.
//
//	Parameter Description
//	{"price":"0xde0b6b3a7640000","worm_address":"0x0000000000000000000000000000000000000004","exchanger":"0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4","block_number":"0x930","sig":"0xfa6cac0a88e4792a45b7f743a1f3737d70e4f100e3f8b10a404617fcbaa706130f617e785edc0cc5796758ca2dba82ea422a18b6624b63b4b2ee412713d243651c"}
//	{"exchanger_owner":"0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4","to":"0xEaE404DCa7c22A15A59f63002Df54BBb8D90c5FB","block_number":"0x92b","sig":"0x972099c287a8da54bb13e7134fcd7edcf96122f1dc949ab987961072011e57662ccb9482ed3738fcdefa613a4d7f58b02fffdf4702943e48bc93af3be7af34191c"}
//	to            "0x5051B76579BC966A9480dd6E72B39A4C89c1154C",	Buyer's address
func (worm *Wormholes) NftExchangeMatch(buyer, seller, exchangerAuth []byte, to string) (string, error) {
	err := tools.CheckAddress("NftExchangeMatch() to", to)
	if err != nil {
		return "", err
	}
	toAddr := common.HexToAddress(to)

	var buyers types2.Buyer
	err = json.Unmarshal(buyer, &buyers)
	if err != nil {
		return "", xerrors.New("the formate of buyer is wrong")
	}

	err = tools.CheckHex("buyers.BlockNumber", buyers.BlockNumber)
	if err != nil {
		return "", err
	}

	var sellers types2.Seller1
	err = json.Unmarshal(seller, &sellers)
	if err != nil {
		return "", xerrors.New("the formate of sellers is wrong")
	}

	err = tools.CheckHex("sellers.BlockNumber", sellers.BlockNumber)
	if err != nil {
		return "", err
	}

	var exchangeAuths types2.ExchangerAuth
	err = json.Unmarshal(exchangerAuth, &exchangeAuths)
	if err != nil {
		return "", xerrors.New("the formate of exchangerAuth is wrong")
	}

	err = tools.CheckHex("exchangeAuths.BlockNumber", exchangeAuths.BlockNumber)
	if err != nil {
		return "", err
	}

	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("NftExchangeMatch() priKeyToAddress err ", err)
		return "", err
	}

	ctx := context.Background()
	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(140000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("NftExchangeMatch() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:          types2.NftExchangeMatch,
		Buyer:         &buyers,
		Seller1:       &sellers,
		ExchangerAuth: &exchangeAuths,
		Version:       types2.WormHolesVersion,
	}
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("NftExchangeMatch() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(buyers.Amount)
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("NftExchangeMatch() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("NftExchangeMatch signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("NftExchangeMatch sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// FoundryExchangeInitiated
//	It is used to buy and sell unminted NFTs. The transaction originator is the exchange. The transaction is used when exchange A authorizes another exchange B, and exchange B initiates the transaction
//
//	Parameter Description
//	buyer:       {"price":"0xde0b6b3a7640000","exchanger":"0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4","block_number":"0x2b","sig":"0xd41864b0f26a605e92d89b0afe508962f89384f7d77dbdca6efd23e7138b84790330a2cdcde2c5cd8653e7f753f244acdea57781e58e713ef93c1568fa8a79cd1c"}
//	seller2:      {"price":"0x38D7EA4C68000","royalty":"0xa","meta_url":"/ipfs/qqqqqqqqqq","exclusive_flag":"0","exchanger":"0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4","block_number":"0x24","sig":"0x836f3e13f001f89d106ddb1e386c5749767b094d54311d950204e9a2594af02a1a9b4d50a425c4e7dfa173088519db7ac5d18ba6acf620fe08036bbf8c2be4e41b"}
//	exchangerAuth:	{"exchanger_owner":"0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4","to":"0xEaE404DCa7c22A15A59f63002Df54BBb8D90c5FB","block_number":"0x26","sig":"0x8c1706b407f50ed5cec8a392eac5f66f0338e9cf4eb71a465dc264ac7e315d2068f6061dfec02ee6b6f7f1150d1594c829436c36bc49c806ee5f5b4ad04e43631c"}
//	to:            "0x5051B76579BC966A9480dd6E72B39A4C89c1154C",	Buyer's address
func (worm *Wormholes) FoundryExchangeInitiated(buyer, seller2, exchangerAuth []byte, to string) (string, error) {
	err := tools.CheckAddress("FoundryExchangeInitiated() to", to)
	if err != nil {
		return "", err
	}

	var buyers types2.Buyer
	err = json.Unmarshal(buyer, &buyers)
	if err != nil {
		return "", xerrors.New("the formate of buyer is wrong")
	}

	err = tools.CheckHex("buyers.BlockNumber", buyers.BlockNumber)
	if err != nil {
		return "", err
	}

	var seller2s types2.Seller2
	err = json.Unmarshal(seller2, &seller2s)
	if err != nil {
		return "", xerrors.New("the formate of seller2 is wrong")
	}

	err = tools.CheckFlag("seller2s.ExclusiveFlag", seller2s.ExclusiveFlag)
	if err != nil {
		return "", err
	}

	err = tools.CheckHex("seller2s.BlockNumber", seller2s.BlockNumber)
	if err != nil {
		return "", err
	}

	//sellerMsg := seller2s.Amount +
	//	seller2s.Royalty +
	//	seller2s.MetaURL +
	//	seller2s.ExclusiveFlag +
	//	seller2s.Exchanger +
	//	seller2s.BlockNumber
	//
	//addr, _ := tools.RecoverAddress(sellerMsg, seller2s.Sig)
	//fmt.Println("---------------seller", addr.String())

	if buyers.Amount < seller2s.Amount {
		return "", xerrors.New("buyer`s amount must be greater then seller`s amount")
	}
	if seller2s.Exchanger != buyers.Exchanger {
		return "", xerrors.New("buyer`s exchanger and seller`s exchanger and transaction`s exchanger aren`t same")
	}

	var exchangerAuths types2.ExchangerAuth
	err = json.Unmarshal([]byte(exchangerAuth), &exchangerAuths)
	if err != nil {
		return "", xerrors.New("the formate of exchangerAuthor is wrong")
	}

	err = tools.CheckHex("exchangeAuths.BlockNumber", exchangerAuths.BlockNumber)
	if err != nil {
		return "", err
	}

	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("FoundryExchangeInitiated() priKeyToAddress err ", err)
		return "", err
	}

	toAddr := common.HexToAddress(to)

	ctx := context.Background()

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(170000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("FoundryExchangeInitiated() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:          types2.FoundryExchangeInitiated,
		Buyer:         &buyers,
		Seller2:       &seller2s,
		ExchangerAuth: &exchangerAuths,
		Version:       types2.WormHolesVersion,
	}
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("FoundryExchangeInitiated() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(buyers.Amount)
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("FoundryExchangeInitiated() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("FoundryExchangeInitiated() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("FoundryExchangeInitiated() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// NFTDoesNotAuthorizeExchanges
//	Used to buy and sell NFTs that have been minted, the transaction originator is the exchange, and the transaction is used when the NFT is not authorized to the exchange
//
//	Parameter Description
//	buyer:  {"price":"0xde0b6b3a7640000","worm_address":"0x0000000000000000000000000000000000000002","exchanger":"0x5051B76579BC966A9480dd6E72B39A4C89c1154C","block_number":"0x11b","sig":"0x158f0ba9dedac427a7746e78aef44ff64c5affa749e56e28793bec6af2a1ff2804a5fd1cce251c84e08674333424a99c8b7497a92f30ed74ceddfc482940ebaa1c"}
//	seller1: {"price":"0xde0b6b3a7640000","worm_address":"0x0000000000000000000000000000000000000002","exchanger":"0x5051B76579BC966A9480dd6E72B39A4C89c1154C","block_number":"0x113","sig":"0x1c8559524220b49e6b9548be405331228d8f26ced8ce12e81b672443fe28067327eef62ce2b3826e2e9ec10f8b2cf5d8a2b2519a0e95f288ea3f098fdea6ab6b1c"}
//	to:      "0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4",		Buyer's address
func (worm *Wormholes) NFTDoesNotAuthorizeExchanges(buyer, seller1 []byte, to string) (string, error) {
	err := tools.CheckAddress("FtDoesNotAuthorizeExchanges() to", to)
	if err != nil {
		return "", err
	}
	var buyers types2.Buyer
	err = json.Unmarshal(buyer, &buyers)
	if err != nil {
		return "", xerrors.New("the formate of buyer is wrong")
	}

	err = tools.CheckHex("buyers.BlockNumber", buyers.BlockNumber)
	if err != nil {
		return "", err
	}

	var seller1s types2.Seller1
	err = json.Unmarshal(seller1, &seller1s)
	if err != nil {
		return "", xerrors.New("the formate of buyer is wrong")
	}

	err = tools.CheckHex("seller1s.BlockNumber", seller1s.BlockNumber)
	if err != nil {
		return "", err
	}

	if buyers.Amount < seller1s.Amount {
		return "", xerrors.New("buyer`s amount must be greater then seller`s amount")
	}
	if seller1s.Exchanger != buyers.Exchanger {
		return "", xerrors.New("buyer`s exchanger and seller`s exchanger and transaction`s exchanger aren`t same")
	}

	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("FtDoesNotAuthorizeExchanges() priKeyToAddress err ", err)
		return "", err
	}

	toAddr := common.HexToAddress(to)

	ctx := context.Background()
	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(130000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("FtDoesNotAuthorizeExchanges() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.FtDoesNotAuthorizeExchanges,
		Buyer:   &buyers,
		Seller1: &seller1s,
		Version: types2.WormHolesVersion,
	}
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("FtDoesNotAuthorizeExchanges() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(buyers.Amount)
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("FtDoesNotAuthorizeExchanges() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("FtDoesNotAuthorizeExchanges() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("FtDoesNotAuthorizeExchanges() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// AdditionalPledgeAmount
//	The amount used by the exchange to increase the pledged ERB
//
//	Parameter Description
//	value:  100,		Append amount, format is hex string
func (worm *Wormholes) AdditionalPledgeAmount(value int64) (string, error) {
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("AdditionalPledgeAmount() priKeyToAddress err ", err)
		return "", err
	}

	ctx := context.Background()
	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(55000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("AdditionalPledgeAmount() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.AdditionalPledgeAmount,
		Version: types2.WormHolesVersion,
	}
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("AdditionalPledgeAmount() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	additional := big.NewInt(value)
	tx := types.NewTransaction(nonce, account, additional, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("AdditionalPledgeAmount() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("AdditionalPledgeAmount() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("AdditionalPledgeAmount() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// RevokesPledgeAmount
//	Amount used for exchanges to reduce the amount of staked ERB
//
//	Parameter Description
//	value:  100,		Amount to decrease, format is hexadecimal string
func (worm *Wormholes) RevokesPledgeAmount(value int64) (string, error) {
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("RevokesPledgeAmount() priKeyToAddress err ", err)
		return "", err
	}

	ctx := context.Background()
	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(55000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("RevokesPledgeAmount() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.RevokesPledgeAmount,
		Version: types2.WormHolesVersion,
	}
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("RevokesPledgeAmount() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	revokes := big.NewInt(value)
	tx := types.NewTransaction(nonce, account, revokes, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("RevokesPledgeAmount() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("RevokesPledgeAmount() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("RevokesPledgeAmount() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// VoteOfficialNFT
//	This transaction is used to inject NFT fragments that can be mined by miners. Only official accounts can do this transaction
//
//	Parameter Descriptiom
//	dir:        "wormholes",  													The path address where sworm is located, the format is a string
//	startIndex: "0x640001",	 														The start number of the sworm fragment, formatted as a hexadecimal string
//	number:     6553600,														The number of sworm shards injected, formatted as a decimal string
//	royalty:    20,																			Royalty, formatted as an integer
//	creator:    "0xab7624f47fd7dadb6b8e255d06a2f10af55990fe",	creator, format is a hex string
func (worm *Wormholes) VoteOfficialNFT(dir, startIndex string, number uint64, royalty uint32, creator string) (string, error) {
	err := tools.CheckAddress("VoteOfficialNFT() creator", creator)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("VoteOfficialNFT() priKeyToAddress err ", err)
		return "", err
	}

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(60000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("VoteOfficialNFT() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:       types2.VoteOfficialNFT,
		Dir:        dir,
		StartIndex: startIndex,
		Number:     number,
		Royalty:    royalty,
		Creator:    creator,
		Version:    types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("VoteOfficialNFT() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, account, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("VoteOfficialNFT() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("VoteOfficialNFT() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("VoteOfficialNFT() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// VoteOfficialNFTByApprovedExchanger
//	This transaction is used to inject NFT fragments that can be mined by miners. Only official accounts can do this transaction
//
//	Parameter Descriptiom
//	dir:        "wormholes",  													The path address where sworm is located, the format is a string
//	startIndex: "0x640001",	 														The start number of the sworm fragment, formatted as a hexadecimal string
//	number:     6553600,														The number of sworm shards injected, formatted as a decimal string
//	royalty:    20,																			Royalty, formatted as an integer
//  exchanger:	{"exchanger_owner":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","to":"0xB685EB3226d5F0D549607D2cC18672b756fd090c","block_number":"0x0","sig":"0xae18a165e51e322d04d2862b6e2760d0493b58870f9afe3c6d15b6e44145c293075662043611501c89d3e4b299a21fe1f8581def86cce4dd43b20c47960ac2481c"}
//	creator:    "0xab7624f47fd7dadb6b8e255d06a2f10af55990fe",	creator, format is a hex string
func (worm *Wormholes) VoteOfficialNFTByApprovedExchanger(dir, startIndex string, number uint64, royalty uint32, creator string, exchangerAuth []byte) (string, error) {
	err := tools.CheckAddress("VoteOfficialNFTByApprovedExchanger() creator", creator)
	if err != nil {
		return "", err
	}

	var exchangeAuths types2.ExchangerAuth
	err = json.Unmarshal(exchangerAuth, &exchangeAuths)
	if err != nil {
		return "", xerrors.New("the formate of exchangerAuth is wrong")
	}

	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("VoteOfficialNFTByApprovedExchanger() priKeyToAddress err ", err)
		return "", err
	}

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(60000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("VoteOfficialNFTByApprovedExchanger() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:          types2.VoteOfficialNFTByApprovedExchanger,
		Dir:           dir,
		StartIndex:    startIndex,
		Number:        number,
		Royalty:       royalty,
		Creator:       creator,
		ExchangerAuth: &exchangeAuths,
		Version:       types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("VoteOfficialNFTByApprovedExchanger() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, account, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("VoteOfficialNFTByApprovedExchanger() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("VoteOfficialNFTByApprovedExchanger() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("VoteOfficialNFTByApprovedExchanger() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// UnforzenAccount
//	change revenue model
func (worm *Wormholes) UnforzenAccount() (string, error) {
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("VoteOfficialNFTByApprovedExchanger() priKeyToAddress err ", err)
		return "", err
	}

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(50000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.UnforzenAccount,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("VoteOfficialNFTByApprovedExchanger() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, account, nil, gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("VoteOfficialNFTByApprovedExchanger() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("VoteOfficialNFTByApprovedExchanger() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("VoteOfficialNFTByApprovedExchanger() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

//AccountDelegate
//Delegate large accounts to small accounts
// Parameter Description
// proxyAddress:		0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4
func (worm *Wormholes) AccountDelegate(proxySign []byte, proxyAddress string) (string, error) {
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		log.Println("AccountDelegate() priKeyToAddress err ", err)
		return "", err
	}

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(70000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("AccountDelegate() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:         types2.AccountDelegate,
		ProxyAddress: proxyAddress,
		ProxySign:    string(proxySign),
		Version:      types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("AccountDelegate() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, account, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("AccountDelegate() networkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("AccountDelegate() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("AccountSign() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

func (worm *Wormholes) RecoverCoefficient() (string, error) {
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(worm.priKey)
	if err != nil {
		return "", err
	}

	nonce, err := worm.PendingNonceAt(ctx, account)

	gasLimit := uint64(60000)
	gasPrice, err := worm.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("RecoverCoefficient() suggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.RecoverCoefficient,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("RecoverCoefficient() failed to format wormholes data")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)

	tx := types.NewTransaction(nonce, account, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := worm.NetworkID(ctx)
	if err != nil {
		log.Println("RecoverCoefficient() networkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("RecoverCoefficient() signTx err ", err)
		return "", err
	}
	err = worm.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("RecoverCoefficient() sendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

var _ APIs = &Wormholes{}
