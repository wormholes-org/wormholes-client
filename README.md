# introduction

- ## 安装

```go
go get github.com/wormholes-org/wormholesclient
```

- ## 客户端

    - ### 创建客户端

      用Go初始化以太坊客户端是和区块链交互所需的基本步骤。导入eth_client包并通过调用接收区块链服务提供者rawurl来初始化它

      ```
      worm := client.NewClient(priKey, rawurl)
      ```

      #### 示例

      ```
      package main
      import (
          "github.com/wormholesclient/client"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          _ = worm
      }
      ```



- ## 签名

    - ### 签名买家

      ```
      SignBuyer(amount, nftAddress, exchanger, blockNumber, seller string) ([]byte, error)
      ```

      进行NFT交易时，对买家信息进行签名

      **参数**

      > - *amount:		   买家购买NFT的金额， 格式为十六进制字符串*
      > - *nftAddress: 交易的NFT地址， 格式为十六进制字符串，填写该字段时，表示交易已铸造的nft，不填写时，表示惰性交易，nft未铸造*
      > - *exchanger :      交易发生所在的交易所， 格式为十进制字符串*
      > - *blockNumber: 块高度， 代表此交易在此高度之前发生才有效， 格式为十六进制字符串*
      > - *seller: 卖家地址， 格式是十六进制字符串*

      **返回**

      `[]byte`  - 交易成功，返回签名的byte数组；签名失败，返回nil

      `error`   - 交易成功，返回nil；交易失败，返回相应的错误

      **示例**

      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, "")
         	buyer, err := worm.Wallet.SignBuyer("0xde0b6b3a7640000", "0x0000000000000000000000000000000000000002", "0x8b07aff2327a3B7e2876D899caFac99f7AE16B10", "0x487", "")
      	if err != nil {
      		log.Fatalln("签名失败")
      	}
      
      	fmt.Println("sign ", string(buyer))
      	//{ "price":"0xde0b6b3a7640000", "nft_address":"0x0000000000000000000000000000000000000002", "exchanger":"0x8b07aff2327a3B7e2876D899caFac99f7AE16B10", "block_number":"0x487", "sig":"0x24355436e991443b8ed3fb83e8c2fa02f8e2bfc0f716c320f836ee7d756e3c712e7e2510b994d1cb7be85d6643233abc81c23929ce7c1c1effd93db261aac5211b" }
      }
      ```

    - ### 签名卖家

      进行NFT交易时，对卖家信息进行签名

        - #### nft已铸造

          ```
          SignSeller1(amount, nftAddress, exchanger, blockNumber string) ([]byte, error)
          ```

          **参数**

          > - *amount:		   买家购买NFT的金额， 格式为十六进制字符串*
          > - *nftAddress: 交易的NFT地址， 格式为十六进制字符串*
          > - *exchanger:       交易发生所在的交易所， 格式为十进制字符串*
          > - *blockNumber: 块高度， 代表此交易在此高度之前发生才有效， 格式为十六进制字符串*

          **返回**

          `[]byte`  - 交易成功，返回签名的byte数组；签名失败，返回nil

          `error`   - 交易成功，返回nil；交易失败，返回相应的错误

          **示例**

          ```
          package main
          import (
              "github.com/wormholesclient/client"
              "fmt"
          )
          
          const (
              priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
          )
          
          func main() {
              worm := client.NewClient(priKey, "")
             	seller1, err := worm.Wallet.SignSeller1("0x38D7EA4C68000", "0x0000000000000000000000000000000000000003", "0x8b07aff2327a3B7e2876D899caFac99f7AE16B10", "0x65d")
          	if err != nil {
          		log.Fatalln("签名失败")
          	}
          
          	fmt.Println("sign ", string(seller1))
          	//seller1: { "price":"0x38D7EA4C68000", "nft_address":"0x0000000000000000000000000000000000000003", "exchanger":"0x8b07aff2327a3B7e2876D899caFac99f7AE16B10", "block_number":"0x65d", "sig":"0x94e88fb5686551dfc3006c608423983a248df8502cbbcaeb2c3352f267a25e531d5fc745bea5f7f564b7399fb70d87026bbf9952f1403e9d4dae4aa14b091cff1c" }
          }
          ```

        - #### nft未铸造

          ```
          SignSeller2(amount, royalty, metaURL, exclusiveFlag, exchanger, blockNumber string) ([]byte, error)
          ```

          **参数**

          > - *amount:		   NFT交易的金额， 格式为十六进制字符串*
          > - *royalty: 版税，十六进制字符串*
          > - *metaURL: NFT元数据地址*
          > - *exclusiveFlag: ”0”:非独占,”1”:独占*
          > - *exchanger:       交易发生所在的交易所， 格式为十进制字符串*
          > - *blockNumber: 块高度， 代表此交易在此高度之前发生才有效， 格式为十六进制字符串*

          **返回**

          `[]byte`  - 交易成功，返回签名的byte数组；签名失败，返回nil

          `error`   - 交易成功，返回nil；交易失败，返回相应的错误
          
          **示例**
          
          ```
          package main
          import (
              "github.com/wormholesclient/client"
              "fmt"
          )
          
          const (
              priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
          )
          
          func main() {
              worm := client.NewClient(priKey, "")
             	seller2, err := worm1.Wallet.SignSeller2("0x38D7EA4C68000", "0xa", "/ipfs/qqqqqqqqqq", "0", exchangeAddress, "0x7be")
          	if err != nil {
          		log.Fatalln("签名失败")
          	}
          
          	fmt.Println("sign ", string(seller2))
          	//seller2: {"price":"0x38D7EA4C68000","royalty":"0xa","meta_url":"/ipfs/qqqqqqqqqq","exclusive_flag":"0","exchanger":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","block_number":"0x7be","sig":"0x84c0c293298557e38fa5064a6fb3b9e6930fa46b234fcd0a923cd677369f5aad3f014a164b21077f713e25b4e986673f614f6ce824561fbda2b4e67e018fac6f1b"}
          ```

    - ### 交易所签名

      ```
      SignExchanger(exchangerOwner, to, blockNumber string) ([]byte, error)
      ```

      进行NFT交易时，对交易所信息进行签名

      **参数**

      > - *exchangerOwner: 授权交易所， 格式为十六进制字符串*
      > - *to:      被授权交易所， 格式为十六进制字符串*
      > - *blockNumber: 块高度， 代表此交易在此高度之前发生才有效， 格式为十六进制字符串*

      **返回**

      `[]byte`  - 交易成功，返回签名的byte数组；签名失败，返回nil

      `error`   - 交易成功，返回nil；交易失败，返回相应的错误

      **示例**

      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, "")
         	exchangeAuth, err := worm2.Wallet.SignExchanger(exchangeAddress, exchangeAddress1, "0x26")
      	if err != nil {
      		log.Fatalln("签名失败")
      	}
      
      	fmt.Println("sign ", string(exchangeAuth))
      	//exchangerAuth:	{"exchanger_owner":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","to":"0xB685EB3226d5F0D549607D2cC18672b756fd090c","block_number":"0x26","sig":"0x8c1706b407f50ed5cec8a392eac5f66f0338e9cf4eb71a465dc264ac7e315d2068f6061dfec02ee6b6f7f1150d1594c829436c36bc49c806ee5f5b4ad04e43631c"}
      ```

- ## NFT接口调用

    - ### Mint

      ```
      Mint(royalty uint32, metaURL string, exchanger string) (string, error)
      ```

      实现NFT的铸造，在wormholes链上创建一个NFT

      **参数**

      > - *royalty:         版税，格式为整数类型*
      > - metaURL:       NFT元数据地址，格式为为字符串*
      > - *exchanger:  铸造NFT时的交易所，格式为字符串，填写该字段时，该交易所独占此NFT，不填写时，没有交易    所独占此NFT*

      **返回**

      `string`  - 交易成功，返回交易的hash；交易失败，返回nil

      `error`   - 交易成功，返回nil；交易失败，返回相应的错误

      **示例**

      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.Mint(10,"/ipfs/ddfd90be9408b4","0x8b07aff2327a3B7e2876D899caFac99f7AE16B10")
          fmt.Println(rs) //0x8fa2d4b70013407012d002fa395939cb0d322553e4848aaae78d4fad638bef55
      }
      ```
      
    - ### Transfer

      ```
      Transfer(nftAddress, to string) (string, error)
      ```
    
      变更NFT的所有权，变更后的拥有者拥有该NFT

      **参数**

      > - *nftAddress:   NFT地址， 格式为十进制字符串，当为SNFT时，长度可以少于42(包含0x),代表合成的SNFT*
      > - *to:                  新拥有者地址*
    
      **返回**

      `string`  - 交易成功，返回交易的hash；交易失败，返回nil

      `error`   - 交易成功，返回nil；交易失败，返回相应的错误

      **示例**

      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.Transfer("0x0000000000000000000000000000000000000001", "0xcdcefddfd90be9408b4965341567182ac8f8a91a")
          fmt.Println(rs) //0x5e8dd659b0ceb95ab53ce32d37daa8688accab601ce58c75e706f08bb47617f4
      }
      ```
      
    - ### Author
    
      ```
      Author(nftAddress, to string) (string, error)
      ```

      NFT授权给某个交易所，使得该交易所拥有NFT的售卖权
    
      **参数**

      > - *nftAddress:   NFT地址， 格式为十进制字符串，当为SNFT时，长度可以少于42(包含0x),代表合成的SNFT*
      > - to:                 被授权者*
    
      **返回**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil

      `error`   - 交易成功，返回nil；交易失败，返回相应的错误

      **示例**

      ```
      package main
      import (
           "github.com/wormholesclient/client"
           "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
      
          rs, _ := worm.NFT.Author("0x0000000000000000000000000000000000000001","0x8b07aff2327a3B7e2876D899caFac99f7AE16B10")
          fmt.Println(rs) //0x2657d46f0c4ef16cadbc6842896c1b50f41333d6a247ee43e5085da5d7e3feff
      }
      ```
      
    - ### AuthorRevoke
    
      ```
      AuthorRevoke(nftAddress, to string) (string, error)
      ```

      将NFT授权给某个交易所，使得该交易所拥有NFT的售卖权

      **参数**
    
      > - *nftAddress:   NFT地址， 格式为十进制字符串，当为SNFT时，长度可以少于42(包含0x),代表合成的SNFT*
      > - *to:                 被取消授权者地址*
    
      **返回**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误

      **示例**

      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.AuthorRevoke("0x0000000000000000000000000000000000000001", "0x8b07aff2327a3B7e2876D899caFac99f7AE16B10")
          fmt.Println(rs) //0xe043dc7d8505d01f6cd949b7a7cc4ed685a9e1b640195801c3c6265b7d11efee
      }
      ```
      
    - ### AccountAuthor
    
      ```
      AccountAuthor(to string) (string, error)
      ```
      
      将账户下的所有NFT授权给某个交易所，使得该交易所拥有所有NFT的售卖权

      **参数**

      > - *to:                被授权者地址*

      **返回**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil

      `error`   - 交易成功，返回nil；交易失败，返回相应的错误

      **示例**

      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
      priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.AccountAuthor("0x8b07aff2327a3B7e2876D899caFac99f7AE16B10")
          fmt.Println(rs) //0x6b42237b9dad13211d89f1e6c66cf947bb371f407a4621ffcf7fd73e385f6fea
      }
      ```
      
    - ### AccountAuthorRevoke
    
      ```
      AccountAuthorRevoke(to string) (string, error)
      ```
    
      取消所有被授权的NFT
    
      **参数**
    
      > - *to:                被取消授权者地址*

      **返回**

      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误

      **示例**

      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
      priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
    
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.AccountAuthorRevoke("0x8b07aff2327a3B7e2876D899caFac99f7AE16B10")
          fmt.Println(rs) //0x1dee05dff7ea39874ed8401c91288ae627b56ae1df6dc4c26a856fafab0447c5
      }
      ```
      
    - ### SNFTToERB
    
      ```
      SNFTToERB(nftAddress string) (string, error)
      ```
      
      将账户挖到的某个SNFT碎片兑换成ERB
    
      **参数**
    
      > - *nftAddress:        被兑换的snft地址， 格式为十进制字符串，长度可以少于42(包含0x),代表合成的SNFT*
    
      **返回**

      `string`  - 交易成功，返回交易的hash；交易失败，返回nil

      `error`   - 交易成功，返回nil；交易失败，返回相应的错误

      **示例**
    
      ```
      package main
      import (
         "github.com/wormholesclient/client"
         "fmt"
      )
      
      const (
      endpoint = "http://192.168.4.237:8574"
      priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
    
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.SNFTToERB("0x8000000000000000000000000000000000000004")
          fmt.Println(rs) //0x77ff920a3a649378e4c7a58644bece643e379113b2bc99257b894a29e220e157
      }
      ```
      
    - ### TokenPledge
    
      ```
      TokenPledge() (string, error)
      ```
    
      当用户想要成为矿工时，需要先做ERB质押交易，以质押成为矿工
    
      **返回**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误
    
      **示例**
    
      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
    
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.TokenPledge()
          fmt.Println(rs) //0x6ceb02802455ab959964866410f37a2f0fcd78e7e64e87d6c9d8102de7f9974b
      }
      ```
      
    - ### TokenRevokesPledge
    
      ```
      TokenRevokesPledge() (string, error)
      ```
    
      当用户不想做矿工，或者不再想质押如此多的ERB时，可以做ERB撤销质押
    
      **返回**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误
    
      **示例**
    
      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
    
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.TokenRevokesPledge()
          fmt.Println(rs) //0xcbd19c386d8b5944d4a88017680239651edefc527e4ba2c8762ab0df2333a7ca
      }
      ```
      
    - ### Open
    
      ```
        Open(feeRate uint32, name, url string) (string, error)
      ```
    
      当用户想要开设交易所时，可以创建一个交易所
    
      参数
    
      > - *feeRate                  交易所抽成， 格式为整数类型*
      > - *name                      交易所名称， 格式为字符串*
      > - **url                           交易所服务器地址， 格式为字符串*
    
      **返回**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误
    
      **示例**

      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.Open(10,"wormholes","www.kang123456.com")
          fmt.Println(rs) //0xcccd6c9499224e7d216f3bd230447900550b07345841eebd2e62b613f7bd924f
      }
      ```
      
    - ### Close
    
      ```
       close() (string, error)
      ```
    
      该交易用于关闭交易所
    
      **参数**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误
    
      **示例**
    
      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.Close()
          fmt.Println(rs) //0x7d5b3653b67d298e7cce82b92fa720224e93941315528ef19a36c4daf1efe929
      }
      ```
      
    - ### InsertNFTBlock
    
      ```
      InsertNFTBlock(dir, startIndex string, number uint64, royalty uint32, creator string) (string, error)
      ```
    
      此交易是用来注入可供矿工挖取的NFT碎片, 只能有官方特定账户做此交易
    
        **参数**
    
      > - dir                    SNFT所在路径地址， 格式为字符串
      > - startIndex        SNFT碎片的开始编号， 格式为十六进制字符串
      > - number           注入的SNFT碎片数量， 格式为十进制字符串
      > - royalty             版税， 格式为整数类型
      > - creator            创作者， 格式是十六进制字符串
    
      **返回**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误
    
      **示例**
    
      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
    
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.InsertNFTBlock("wormholes2", "0x640001", 6553600, 20, "0xab7624f47fd7dadb6b8e255d06a2f10af55990fe")
          fmt.Println(rs) //0x61cd018d6e70af47c6204fea18db5b33fdecc92162cca66b0089783733809e84
      }
      ```
      
    - ### TransactionNFT
    
      ```
      TransactionNFT(buyer []byte, to string) (string, error)
      ```
    
      此交易用于买卖已经铸造过的NFT，该交易发起方可以是交易所，或者卖方
    
      **参数**
    
      > - *buyer:               买家*
      > - *to                      买家地址， 格式是十六进制字符串*
    
      **返回**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误
    
      **示例**
    
      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
    
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.TransactionNFT(buyer,"0x44d952db5dfB4CBb54443554F4bB9cbeBee2194c")
          fmt.Println(rs) //0xc9c4e6652ba411a0435d2e3187f019329b084734f19ae6699ee7f1fa9a92123b
      }
      ```
      
    - ### BuyerInitiatingTransaction
    
      ```
      BuyerInitiatingTransaction(seller1 []byte) (string, error)
      ```
    
      此交易用于买卖已经铸造过的NFT，该交易发起方是买方
    
      **参数**
    
      > - *seller1:             卖家*
      > - *to                      卖家地址， 格式是十六进制字符串*
    
      **返回**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误
    
      **示例**
    
      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.BuyerInitiatingTransaction(seller1, "0x18e534cb9ec13b846ed143e6bd1bb0881188c085")
          fmt.Println(rs) //0xfb9cf0100340c9bf965fc0f8ef44bb8a75af58175deab0dcff3979a97a8ebefa
      }
      ```
    
    - ### FoundryTradeBuyer
    
      ```
       FoundryTradeBuyer(seller2 []byte) (string, error)
      ```
    
      用于买卖未铸造过的NFT，该交易发起方是买方
    
      **参数**
    
      > - *seller2:             卖家*
      > - *to:                    卖家地址*
    
      **返回**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误
    
      **示例**
    
      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
    
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.FoundryTradeBuyer(seller2, "0x44d952db5dfb4cbb54443554f4bb9cbebee2194c")
          fmt.Println(rs) //0x4634d6bbc36b9444914a259c2acf0410af0b99122baef30d7a8701a496bc3b6c
      }
      ```
      
    - ### FoundryExchange
    
      ```
       FoundryExchange(buyer, seller2 []byte, to string) (string, error)
      ```

      用于买卖未铸造过的NFT，该交易发起方是交易所，或者卖方

      **参数**

      > - *buyer:             买家*
      > - *seller2:           卖家*
      > - *to:                    买家地址*

      **返回**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误
    
      **示例**
    
      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.FoundryTradeBuyer(buyer,seller2 "0x44d952db5dfb4cbb54443554f4bb9cbebee2194c")
          fmt.Println(rs) //0x70853466fdf5dc4476fab34b79f9be2e66f0448789937094de0b0aa5f3345e8c
      }
      ```
    
    - ### NftExchangeMatch

      ```
       NftExchangeMatch(buyer, exchangerAuth []byte, to string) (string, error)
      ```
    
      用于买卖已经铸造过的NFT，该交易发起方是交易所，该交易是当交易所A授权给另一交易所B，由交易所B发起交易时使用

      **参数**

      > - *buyer                         买家*
      > - *exchangerAuth:       授权交易所A*
      > - *to                               买家， 格式是十六进制字符串*

      **返回**

      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误
    
      **示例**
    
      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.NftExchangeMatch(buyer,exchangAuth,"0x44d952db5dfB4CBb54443554F4bB9cbeBee2194c")
          fmt.Println(rs) //0xf11e024297b89e6dfd02bc2da4680cea353ea6956c3ea9084afa40d58477932f
      }
      ```
      
    - ### FoundryExchangeInitiated

      ```
       FoundryExchangeInitiated(buyer, seller2, exchangerAuthor []byte, to string) (string, error)
      ```

      用于买卖未铸造过的NFT，该交易发起方是交易所，该交易是当交易所A授权给另一交易所B，由交易所B发起交易时使用
    
      **参数**

      > - *buyer                         买家*
      > - *Seller2                       卖家*
      > - *exchangerAuth:       授权交易所A*
      > - *to                               买家， 格式是十六进制字符串*
    
      **返回**
    
      `string`  - 交易成功，返回交易的hash；交易失败，返回nil
    
      `error`   - 交易成功，返回nil；交易失败，返回相应的错误
    
      **示例**
    
      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.NftExchangeMatch(buyer,seller2,exchangAuth,  "0x44d952db5dfB4CBb54443554F4bB9cbeBee2194c")
          fmt.Println(rs) //0xc9cc570057faf1edd83f48833520f9d546e4972083ee705152b5f35630f1588d
      }
      ```
      
    - ### FtDoesNotAuthorizeExchanges
    
      ```
       FtDoesNotAuthorizeExchanges(buyer, seller1 []byte, to string) (string, error)
      ```
    
      用于买卖已经铸造过的NFT，该交易发起方是交易所，该交易是当NFT未授权给交易所时使用
    
      **参数**
    
      > - *buyer               买家*
      > - *seller1:             卖家*
      > - *to                     买家， 格式是十六进制字符串*

      **返回**

      `string`  - 交易成功，返回交易的hash；交易失败，返回nil

      `error`   - 交易成功，返回nil；交易失败，返回相应的错误

      **示例**
    
      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.FtDoesNotAuthorizeExchanges(buyer,seller1,"0x44d952db5dfb4cbb54443554f4bb9cbebee2194c")
          fmt.Println(rs) //0x95615a6c7a164537257492c112a9fcd99907315893706a1b104456d9e3aa8af6
      }
      ```
    
    - ### AdditionalPledgeAmount

      ```
       AdditionalPledgeAmount(value int64) (string, error)
      ```

      用于交易所增加质押ERB的金额

      **参数**

      > - Values:         追加的金额
    
      **返回**

      `string`  - 交易成功，返回交易的hash；交易失败，返回nil

      `error`   - 交易成功，返回nil；交易失败，返回相应的错误

      **示例**

      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.AdditionalPledgeAmount(100)
          fmt.Println(rs) //0x25f2ed8cf5f1041be9e71d483a32b01fd3f7820ec59e0c060830214c53fea5f9
      }
      ```
      
    - ### RevokesPledgeAmount
    
      ```
       RevokesPledgeAmount(value int64) (string, error)
      ```

      用于交易所增加质押ERB的金额
    
      **参数**

      > - Values:          减少的金额

      **返回**

      `string`  - 交易成功，返回交易的hash；交易失败，返回nil

      `error`   - 交易成功，返回nil；交易失败，返回相应的错误

      **示例**

      ```
      package main
      import (
          "github.com/wormholesclient/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "7c6786275d6011adb6288587757653d3f9061275bafc2c35ae62efe0bc4973e9"
      )
      
      func main() {
          worm := client.NewNFT(priKey, endpoint)
          rs, _ := worm.NFT.RevokesPledgeAmount(100)
          fmt.Println(rs) //0xd2c7f943f0f5364b0928c518e7b6de7491c0e8efb6abf912a17e6860f70ebec1
      }
      ```