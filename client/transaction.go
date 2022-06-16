package client

import (
	"context"
	"encoding/json"
	"github.com/wormholes-org/wormholesclient/tools"
	types2 "github.com/wormholes-org/wormholesclient/types"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/xerrors"
	"log"
	"math/big"
	"strings"
)

// Mint NFT用户铸造
//	用户可以使用该交易在wormholes链上创建一个NFT
//
//	参数说明
//	royalty: 10,																					版税， 格式为整数类型
//	metaURL: "/ipfs/ddfd90be9408b4",	NFT元数据地址
//	exchanger:"0x8b07aff2327a3B7e2876D899caFac99f7AE16B10",							铸造NFT时的交易所，格式为字符串，填写该字段时，该交易所独占此NFT，不填写时，没有交易所独占此NFT

func (nft *NFT) Mint(royalty uint32, metaURL string, exchanger string) (string, error) {
	if exchanger != "" {
		err := tools.CheckAddress("exchanger", exchanger)
		if err != nil {
			return "", err
		}
	}

	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		return "", err
	}

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(100000)
	//gasPrice, err := nft.SuggestGasPrice(ctx)
	gasPrice := big.NewInt(50000000000)
	if err != nil {
		log.Println("AuthExchangerMint() err=", err)
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
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, account, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// Transfer NFT转移
// 	变更NFT的所有权
//
//	参数说明
//	nftAddress: "0x8000000000000000000000000000000000000001",  nft地址，格式为十进制字符串，当为SNFT时，长度可以少于42(包含0x),代表合成的SNFT
//	to:         "0xcdcefddfd90be9408b4965341567182ac8f8a91a",  目标NFT用户地址
func (nft *NFT) Transfer(nftAddress, to string) (string, error) {
	err := tools.CheckHex("nftAddress", nftAddress)
	if err != nil {
		return "", err
	}
	err = tools.CheckAddress("To", to)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		return "", err
	}

	toAddr := common.HexToAddress(to)

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("SuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:       types2.Transfer,
		NFTAddress: nftAddress,
		Version:    types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)

	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, toAddr, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// Author 将某个NFT授权给交易所
//
//	参数说明
//	nftAddress: "0x0000000000000000000000000000000000000001",	被授权nft地址， 格式为十进制字符串，当为SNFT时，长度可以少于42(包含0x),代表合成的SNFT
//	to:         "0x8b07aff2327a3B7e2876D899caFac99f7AE16B10",	被授权者地址
func (nft *NFT) Author(nftAddress, to string) (string, error) {
	err := tools.CheckHex("nftAddress", nftAddress)
	if err != nil {
		return "", err
	}
	err = tools.CheckAddress("to", to)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		return "", err
	}

	toAddr := common.HexToAddress(to)

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("SuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:       types2.Author,
		NFTAddress: nftAddress,
		Version:    types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)

	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, toAddr, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// AuthorRevoke 取消某个NFT的授权
//
//	参数说明
//	nftAddress: "0x0000000000000000000000000000000000000002",	被授权nft地址， 格式为十进制字符串，当为SNFT时，长度可以少于42(包含0x),代表合成的SNFT
//	to:         "0x8b07aff2327a3B7e2876D899caFac99f7AE16B10",	被授权者地址
func (nft *NFT) AuthorRevoke(nftAddress, to string) (string, error) {
	err := tools.CheckHex("nftAddress", nftAddress)
	if err != nil {
		return "", err
	}
	err = tools.CheckAddress("to", to)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		return "", err
	}

	toAddr := common.HexToAddress(to)

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("SuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:       types2.AuthorRevoke,
		NFTAddress: nftAddress,
		Version:    types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, toAddr, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// AccountAuthor
//	将某个账户下的所有NFT授权给交易所
//	参数说明
//	to:     "0x8b07aff2327a3B7e2876D899caFac99f7AE16B10",							被授权者地址
func (nft *NFT) AccountAuthor(to string) (string, error) {
	err := tools.CheckAddress("to", to)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		return "", err
	}

	toAddr := common.HexToAddress(to)

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("SuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.AccountAuthor,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, toAddr, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// AccountAuthorRevoke
//	取消某个账户下的所有NFT授权
//
//	参数说明
//	to:     "0x8b07aff2327a3B7e2876D899caFac99f7AE16B10",							被授权者地址
func (nft *NFT) AccountAuthorRevoke(to string) (string, error) {
	err := tools.CheckAddress("to", to)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		return "", err
	}

	toAddr := common.HexToAddress(to)

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("SuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.AccountAuthorRevoke,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, toAddr, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// SNFTToERB
//	将矿工挖出来的NFT碎片兑换成ERB
//
//	参数说明
//	nftAddress: "0x101060AEFE0d70fB40eda7F4a605c1315Be4A72F",	被兑换的snft地址， 格式为十进制字符串，长度可以少于42(包含0x),代表合成的SNFT
//
//	合成级别对应的兑换价格
//	0: 100000000000000000
// 	1: 150000000000000000
//	2: 225000000000000000
//	3: 300000000000000000
func (nft *NFT) SNFTToERB(nftAddress string) (string, error) {
	err := tools.CheckHex("nftAddress", nftAddress)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		return "", err
	}

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("SuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:       types2.SNFTToERB,
		NFTAddress: nftAddress,
		Version:    types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, account, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err ", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// TokenPledge
//	当用户想要成为矿工时，需要先做ERB质押交易，以质押成为矿工所需要的ERB
func (nft *NFT) TokenPledge() (string, error) {
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.TokenPledge,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	wei, _ := new(big.Int).SetString("1000000000000000000", 10)
	pledge := new(big.Int).Mul(big.NewInt(100000), wei)
	tx := types.NewTransaction(nonce, account, pledge, gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// TokenRevokesPledge
//	当用户不想做矿工，或者不再想质押如此多的ERB时，可以做ERB撤销质押
func (nft *NFT) TokenRevokesPledge() (string, error) {
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.TokenRevokesPledge,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	wei, _ := new(big.Int).SetString("1000000000000000000", 10)
	pledge := new(big.Int).Mul(big.NewInt(100000), wei)

	tx := types.NewTransaction(nonce, account, pledge, gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// Open
//	当用户想要开设交易所时，可以发起此交易
//
//	参数说
//	feeRate:   10,																		  交易所抽成，格式为整数类型
//	name:      "wormholes",										 交易所名称，格式为字符串
//	url:       "www.kang123456.com",		交易所服务器地址，格式为字符串
func (nft *NFT) Open(feeRate uint32, name, url string) (string, error) {
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
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
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	wei, _ := new(big.Int).SetString("1000000000000000000", 10)
	amount := new(big.Int).Mul(big.NewInt(100), wei)

	tx := types.NewTransaction(nonce, account, amount, gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// Close
//	用户不想继续开设交易所时，可以发起此交易来关闭开设的交易所
func (nft *NFT) Close() (string, error) {
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.Close,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, account, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// InsertNFTBlock
//	此交易是用来注入可供矿工挖取的NFT碎片, 只能有官方特定账户做此交易
//
//	交易格式
//	dir:        "wormholes",  													snft所在路径地址，格式为字符串
//	startIndex: "0x640001",	 														snft碎片的开始编号， 格式为十六进制字符串
//	number:     6553600,														注入的snft碎片数量， 格式为十进制字符串
//	royalty:    20,																			版税， 格式为整数类型
//	creator:    "0xab7624f47fd7dadb6b8e255d06a2f10af55990fe",	创作者， 格式是十六进制字符串
func (nft *NFT) InsertNFTBlock(dir, startIndex string, number uint64, royalty uint32, creator string) (string, error) {
	err := tools.CheckAddress("creator", creator)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
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
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	tx := types.NewTransaction(nonce, account, big.NewInt(0), gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// TransactionNFT
//	用于买卖已经铸造过的NFT，该交易发起方可以是交易所，或者卖方
//
//	格式说明
//	buyer: { "price":"0xde0b6b3a7640000", "nft_address":"0x0000000000000000000000000000000000000002", "exchanger":"0x8b07aff2327a3B7e2876D899caFac99f7AE16B10", "block_number":"0x487", "sig":"0x24355436e991443b8ed3fb83e8c2fa02f8e2bfc0f716c320f836ee7d756e3c712e7e2510b994d1cb7be85d6643233abc81c23929ce7c1c1effd93db261aac5211b" }																				买家
//	to:     "0x44d952db5dfB4CBb54443554F4bB9cbeBee2194c",				买家地址
func (nft *NFT) TransactionNFT(buyer []byte, to string) (string, error) {
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

	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	ctx := context.Background()
	nonce, err := nft.PendingNonceAt(ctx, account)

	toAddr := common.HexToAddress(to)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
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
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(buyers.Amount)
	fmt.Println(value)
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// BuyerInitiatingTransaction
//	用于买卖已经铸造过的NFT，该交易发起方是买方
//
//	参数说明
//	seller1: { "price":"0x38D7EA4C68000", "nft_address":"0x0000000000000000000000000000000000000003", "exchanger":"0x8b07aff2327a3B7e2876D899caFac99f7AE16B10", "block_number":"0x65d", "sig":"0x94e88fb5686551dfc3006c608423983a248df8502cbbcaeb2c3352f267a25e531d5fc745bea5f7f564b7399fb70d87026bbf9952f1403e9d4dae4aa14b091cff1c" }
func (nft *NFT) BuyerInitiatingTransaction(seller1 []byte) (string, error) {
	var seller1s types2.Seller1
	err := json.Unmarshal(seller1, &seller1s)
	if err != nil {
		return "", xerrors.New("the formate of seller1 is wrong")
	}

	err = tools.CheckHex("seller1s.BlockNumber", seller1s.BlockNumber)
	if err != nil {
		return "", err
	}
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	ctx := context.Background()
	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.BuyerInitiatingTransaction,
		Seller1: &seller1s,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(seller1s.Amount)
	tx := types.NewTransaction(nonce, account, value, gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// FoundryTradeBuyer
//	用于买卖未铸造过的NFT，该交易发起方是买方
//
//	交易格式
//	seller2: { "price":"0x38D7EA4C68000", "royalty":"0xa", "meta_url":"/ipfs/qqqqqqqqqq", "exclusive_flag":"0", "exchanger":"0x8b07aff2327a3B7e2876D899caFac99f7AE16B10", "block_number":"0x703", "sig":"0xb08cf8b2f2d4b2635a85d1c7a816f01c24ac2a90ab49bdbe0e52e0a8f07eea5521eb80554df2c403423bdf49f412a7811b10a16005832a1bc171f5dfd3c983121c" }
func (nft *NFT) FoundryTradeBuyer(seller2 []byte) (string, error) {
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

	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	ctx := context.Background()
	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.FoundryTradeBuyer,
		Seller2: &seller2s,
		Version: types2.WormHolesVersion,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(seller2s.Amount)
	tx := types.NewTransaction(nonce, account, value, gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// FoundryExchange
//	用于买卖未铸造过的NFT，该交易发起方是交易所，或者卖方
//
//	参数说明
//	buyer:   {"price":"0xde0b6b3a7640000","exchanger":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","block_number":"0x7c6","sig":"0xd4d2319bd9c4c1664ceb8cdb4d417fc22a6b4083845d5390154f4d268b07bc81755b0f728f989554142ca8124fe543b93a526f92664d7cc905ec361721ef130a1b"}
//	seller2: {"price":"0x38D7EA4C68000","royalty":"0xa","meta_url":"/ipfs/qqqqqqqqqq","exclusive_flag":"0","exchanger":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","block_number":"0x7be","sig":"0x84c0c293298557e38fa5064a6fb3b9e6930fa46b234fcd0a923cd677369f5aad3f014a164b21077f713e25b4e986673f614f6ce824561fbda2b4e67e018fac6f1b"}
//	to:      "0x44d952db5dfb4cbb54443554f4bb9cbebee2194c",  买家地址
func (nft *NFT) FoundryExchange(buyer, seller2 []byte, to string) (string, error) {
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
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	toAddr := common.HexToAddress(to)

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
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
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(buyers.Amount)
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// NftExchangeMatch
//	用于买卖已经铸造过的NFT，该交易发起方是交易所，该交易是当交易所A授权给另一交易所B，由交易所B发起交易时使用
//
//	参数说明
//	{"price":"0xde0b6b3a7640000","nft_address":"0x0000000000000000000000000000000000000004","exchanger":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","block_number":"0x930","sig":"0xfa6cac0a88e4792a45b7f743a1f3737d70e4f100e3f8b10a404617fcbaa706130f617e785edc0cc5796758ca2dba82ea422a18b6624b63b4b2ee412713d243651c"}
//	{"exchanger_owner":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","to":"0xB685EB3226d5F0D549607D2cC18672b756fd090c","block_number":"0x92b","sig":"0x972099c287a8da54bb13e7134fcd7edcf96122f1dc949ab987961072011e57662ccb9482ed3738fcdefa613a4d7f58b02fffdf4702943e48bc93af3be7af34191c"}
//	to            "0x44d952db5dfB4CBb54443554F4bB9cbeBee2194c",	买家地址
func (nft *NFT) NftExchangeMatch(buyer, exchangerAuth []byte, to string) (string, error) {
	err := tools.CheckAddress("to", to)
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

	var exchangeAuths types2.ExchangerAuth
	err = json.Unmarshal(exchangerAuth, &exchangeAuths)
	if err != nil {
		return "", xerrors.New("the formate of exchangerAuth is wrong")
	}

	err = tools.CheckHex("exchangeAuths.BlockNumber", exchangeAuths.BlockNumber)
	if err != nil {
		return "", err
	}

	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	ctx := context.Background()
	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:          types2.NftExchangeMatch,
		Buyer:         &buyers,
		ExchangerAuth: &exchangeAuths,
		Version:       types2.WormHolesVersion,
	}
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(buyers.Amount)
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// FoundryExchangeInitiated
//	用于买卖未铸造过的NFT，该交易发起方是交易所，该交易是当交易所A授权给另一交易所B，由交易所B发起交易时使用
//
//	参数说明
//	buyer:       {"price":"0xde0b6b3a7640000","exchanger":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","block_number":"0x2b","sig":"0xd41864b0f26a605e92d89b0afe508962f89384f7d77dbdca6efd23e7138b84790330a2cdcde2c5cd8653e7f753f244acdea57781e58e713ef93c1568fa8a79cd1c"}
//	seller2:      {"price":"0x38D7EA4C68000","royalty":"0xa","meta_url":"/ipfs/qqqqqqqqqq","exclusive_flag":"0","exchanger":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","block_number":"0x24","sig":"0x836f3e13f001f89d106ddb1e386c5749767b094d54311d950204e9a2594af02a1a9b4d50a425c4e7dfa173088519db7ac5d18ba6acf620fe08036bbf8c2be4e41b"}
//	exchangerAuth:	{"exchanger_owner":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","to":"0xB685EB3226d5F0D549607D2cC18672b756fd090c","block_number":"0x26","sig":"0x8c1706b407f50ed5cec8a392eac5f66f0338e9cf4eb71a465dc264ac7e315d2068f6061dfec02ee6b6f7f1150d1594c829436c36bc49c806ee5f5b4ad04e43631c"}
//	to:            "0x44d952db5dfB4CBb54443554F4bB9cbeBee2194c",	买家地址
func (nft *NFT) FoundryExchangeInitiated(buyer, seller2, exchangerAuth []byte, to string) (string, error) {
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

	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	toAddr := common.HexToAddress(to)

	ctx := context.Background()

	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
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
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(buyers.Amount)
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// FtDoesNotAuthorizeExchanges
//	用于买卖已经铸造过的NFT，该交易发起方是交易所，该交易是当NFT未授权给交易所时使用
//
//	参数说明
//	buyer:  {"price":"0xde0b6b3a7640000","nft_address":"0x0000000000000000000000000000000000000002","exchanger":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","block_number":"0x11b","sig":"0x158f0ba9dedac427a7746e78aef44ff64c5affa749e56e28793bec6af2a1ff2804a5fd1cce251c84e08674333424a99c8b7497a92f30ed74ceddfc482940ebaa1c"}
//	seller1: {"price":"0xde0b6b3a7640000","nft_address":"0x0000000000000000000000000000000000000002","exchanger":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","block_number":"0x113","sig":"0x1c8559524220b49e6b9548be405331228d8f26ced8ce12e81b672443fe28067327eef62ce2b3826e2e9ec10f8b2cf5d8a2b2519a0e95f288ea3f098fdea6ab6b1c"}
//	to:      "0x44d952db5dfb4cbb54443554f4bb9cbebee2194c",		买家地址
func (nft *NFT) FtDoesNotAuthorizeExchanges(buyer, seller1 []byte, to string) (string, error) {
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

	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	toAddr := common.HexToAddress(to)

	ctx := context.Background()
	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
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
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	value, _ := hexutil.DecodeBig(buyers.Amount)
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// AdditionalPledgeAmount
//	用于交易所增加质押ERB的金额
//
//	参数说明
//	value:  100,		追加金额，格式是十六进制字符串
func (nft *NFT) AdditionalPledgeAmount(value int64) (string, error) {
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	ctx := context.Background()
	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.AdditionalPledgeAmount,
		Version: types2.WormHolesVersion,
	}
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	additional := big.NewInt(value)
	tx := types.NewTransaction(nonce, account, additional, gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

// RevokesPledgeAmount
//	用于交易所减少质押ERB的金额
//
//	参数说明
//	value:  100,		减少金额，格式是十六进制字符串
func (nft *NFT) RevokesPledgeAmount(value int64) (string, error) {
	account, fromKey, err := tools.PriKeyToAddress(nft.priKey)
	if err != nil {
		log.Println("PriKeyToAddress err ", err)
		return "", err
	}

	ctx := context.Background()
	nonce, err := nft.PendingNonceAt(ctx, account)

	gasLimit := uint64(51000)
	gasPrice, err := nft.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("ASuggestGasPrice err ", err)
		return "", err
	}

	transaction := types2.Transaction{
		Type:    types2.RevokesPledgeAmount,
		Version: types2.WormHolesVersion,
	}
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Println("格式化wormholes数据失败")
		return "", err
	}

	tx_data := append([]byte("wormholes:"), data...)
	fmt.Println(string(tx_data))

	revokes := big.NewInt(value)
	tx := types.NewTransaction(nonce, account, revokes, gasLimit, gasPrice, tx_data)
	chainID, err := nft.NetworkID(ctx)
	if err != nil {
		log.Println("NetworkID err=", err)
		return "", err
	}
	log.Println("chainID=", chainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
	if err != nil {
		log.Println("SignTx err ", err)
		return "", err
	}
	err = nft.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("SendTransaction err ", err)
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

var _ APIs = &NFT{}
