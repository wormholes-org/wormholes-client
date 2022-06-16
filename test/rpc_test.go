package test

import (
	"context"
	"encoding/json"
	"github.com/wormholes-org/wormholesclient/client"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGetAccountInfo(t *testing.T) {
	nft := client.NewClient("http://192.168.4.237:8574", "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9")
	ctx := context.Background()
	blockNumber, _ := nft.NFT.BlockNumber(ctx)
	fmt.Println("blockNumber ", blockNumber)
	rs, _ := nft.NFT.GetBlockBeneficiaryAddressByNumber(ctx, int64(blockNumber))
	rss1, _ := json.Marshal(*rs)
	fmt.Println(string(rss1))

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(*rs))
	fmt.Println("Address", (*rs)[n].Address.String())
	fmt.Println("NFTAddress", (*rs)[n].NftAddress.String())
	rs1, _ := nft.NFT.GetAccountInfo(ctx, (*rs)[n].NftAddress, int64(blockNumber))
	rss1, _ = json.Marshal(*rs1)
	fmt.Println(string(rss1))
}
