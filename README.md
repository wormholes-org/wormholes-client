# Introduction

[![Join the chat at https://gitter.im/wormholes-org/miner](https://badges.gitter.im/wormholes-org/miner.svg)](https://gitter.im/wormholes-org/miner?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

- ## Install

```go
go get github.com/wormholes-org/wormholes-client
```

- ## Client

    - ### Create a client

      Initializing the wormholes client with Go is the basic steps required to interact with the blockchain.
      Import the wormholes-org/wormholes-client package and initialize it by calling the receiving blockchain service provider
      rawurl

      ```
      worm := client.NewClient(priKey, rawurl)
      ```

      #### Example

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "88aaf04596c2c9e71c94c1ec5c160d4326346511b28324d6f19efa9716cb66fd"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          _ = worm
      }
      ```



- ## Signature

    - ### Sign buyer

      ```
      SignBuyer(amount, nftAddress, exchanger, blockNumber, seller string) ([]byte, error)
      ```

      When conducting NFT transactions, sign buyer information

      **Params**

      > - *amount:		   The amount the buyer purchased the NFT, formatted as a hexadecimal string*
      > - *nftAddress: The NFT address of the transaction. The format is a hexadecimal string. When this field is filled in, it means that the transaction has minted nft. When not filled, it means lazy transaction, and the nft has not been minted*
      > - *exchanger :      The exchange on which the transaction took place, formatted as a decimal string*
      > - *blockNumber: Block height, which means that this transaction is valid before this height, the format is a hexadecimal string*
      > - *seller: Seller's address, formatted as a hexadecimal string*

      **Return**

      `[]byte`  - If the transaction is successful, return the signed byte array; if the signature fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          priKey   = "88aaf04596c2c9e71c94c1ec5c160d4326346511b28324d6f19efa9716cb66fd"
      )
      
      func main() {
          worm := client.NewClient(priKey, "")
             buyer, err := worm.Wallet.SignBuyer("0xde0b6b3a7640000", "0x0000000000000000000000000000000000000002", "0x8b07aff2327a3B7e2876D899caFac99f7AE16B10", "0x487", "")
          if err != nil {
              log.Fatalln("Signing failed")
          }
      
          fmt.Println("sign ", string(buyer))
          //{ "price":"0xde0b6b3a7640000", "nft_address":"0x0000000000000000000000000000000000000002", "exchanger":"0x8b07aff2327a3B7e2876D899caFac99f7AE16B10", "block_number":"0x487", "sig":"0x24355436e991443b8ed3fb83e8c2fa02f8e2bfc0f716c320f836ee7d756e3c712e7e2510b994d1cb7be85d6643233abc81c23929ce7c1c1effd93db261aac5211b" }
      }
      ```

    - ### Sign seller

      When conducting NFT transactions, sign the seller information

        - #### nft cast

          ```
          SignSeller1(amount, nftAddress, exchanger, blockNumber string) ([]byte, error)
          ```

          **Params**

          > - *amount:		   The amount the buyer purchased the NFT, formatted as a hexadecimal string*
          > - *nftAddress: The NFT address of the transaction, formatted as a hexadecimal string*
          > - *exchanger:       The exchange on which the transaction took place, formatted as a decimal string*
          > - *blockNumber: Block height, which means that this transaction is valid before this height, the format is a hexadecimal string*

          **Return**

          `[]byte`  - If the transaction is successful, return the signed byte array; if the signature fails, return nil

          `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

          **Example**

          ```
          package main
          import (
              "github.com/wormholes-org/wormholes-client/client"
              "fmt"
          )
          
          const (
              priKey   = "11e4259f98e6a18772be5e1b2e2c9e5b12b4a9fe8e3cfa0853df59fa0825e861"
          )
          
          func main() {
              worm := client.NewClient(priKey, "")
                 seller1, err := worm.Wallet.SignSeller1("0x38D7EA4C68000", "0x0000000000000000000000000000000000000003", "0x8b07aff2327a3B7e2876D899caFac99f7AE16B10", "0x65d")
              if err != nil {
                  log.Fatalln("Signing failed")
              }
          
              fmt.Println("sign ", string(seller1))
              //seller1: { "price":"0x38D7EA4C68000", "nft_address":"0x0000000000000000000000000000000000000003", "exchanger":"0x8b07aff2327a3B7e2876D899caFac99f7AE16B10", "block_number":"0x65d", "sig":"0x94e88fb5686551dfc3006c608423983a248df8502cbbcaeb2c3352f267a25e531d5fc745bea5f7f564b7399fb70d87026bbf9952f1403e9d4dae4aa14b091cff1c" }
          }
          ```

        - #### nft uncast

          ```
          SignSeller2(amount, royalty, metaURL, exclusiveFlag, exchanger, blockNumber string) ([]byte, error)
          ```

          **Params**

          > - *amount:		   The amount of the NFT transaction, formatted as a hexadecimal string*
          > - *royalty: royalty, hex string*
          > - *metaURL: NFT metadata address*
          > - *exclusiveFlag: "0": Inclusive, "1": Exclusive*
          > - *exchanger:       The exchange on which the transaction took place, formatted as a decimal string*
          > - *blockNumber: Block height, which means that this transaction is valid before this height, the format is a hexadecimal string*

          **Return**

          `[]byte`  - If the transaction is successful, return the signed byte array; if the signature fails, return nil

          `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

          **Example**

          ```
          package main
          import (
              "github.com/wormholes-org/wormholes-client/client"
              "fmt"
          )
          
          const (
              priKey   = "fc92219d8c663afde9708be321354b9b1c2e67c8680647f73c5bf64ce13cca66"
          )
          
          func main() {
              worm := client.NewClient(priKey, "")
                 seller2, err := worm1.Wallet.SignSeller2("0x38D7EA4C68000", "0xa", "/ipfs/qqqqqqqqqq", "0", exchangeAddress, "0x7be")
              if err != nil {
                  log.Fatalln("Signing failed")
              }
          
              fmt.Println("sign ", string(seller2))
              //seller2: {"price":"0x38D7EA4C68000","royalty":"0xa","meta_url":"/ipfs/qqqqqqqqqq","exclusive_flag":"0","exchanger":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","block_number":"0x7be","sig":"0x84c0c293298557e38fa5064a6fb3b9e6930fa46b234fcd0a923cd677369f5aad3f014a164b21077f713e25b4e986673f614f6ce824561fbda2b4e67e018fac6f1b"}
          ```

    - ### Sign Exchange

      ```
      SignExchanger(exchangerOwner, to, blockNumber string) ([]byte, error)
      ```

      When conducting NFT transactions, sign the exchange information

      **Params**

      > - *exchangerOwner: Authorize exchange, formatted as a hexadecimal string*
      > - *to:      Authorized exchange, formatted as a hexadecimal string*
      > - *blockNumber: Block height, which means that this transaction is valid before this height, the format is a hexadecimal string*

      **Return**

      `[]byte`  - If the transaction is successful, return the signed byte array; if the signature fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          priKey   = "ea86b620d97c256434fabbc46a3350cadd42019d7a5953914a8fc1e1093f214c"
      )
      
      func main() {
          worm := client.NewClient(priKey, "")
             exchangeAuth, err := worm2.Wallet.SignExchanger(exchangeAddress, exchangeAddress1, "0x26")
          if err != nil {
              log.Fatalln("Signing failed")
          }
      
          fmt.Println("sign ", string(exchangeAuth))
          //exchangerAuth:	{"exchanger_owner":"0x83c43f6F7bB4d8E429b21FF303a16b4c99A59b05","to":"0xB685EB3226d5F0D549607D2cC18672b756fd090c","block_number":"0x26","sig":"0x8c1706b407f50ed5cec8a392eac5f66f0338e9cf4eb71a465dc264ac7e315d2068f6061dfec02ee6b6f7f1150d1594c829436c36bc49c806ee5f5b4ad04e43631c"}
      ```

- ## NFT interface

    - ### Recharge

      ```
      Mint(royalty uint32, metaURL string, exchanger string) (string, error)
      ```

      Implement NFT minting, create an NFT on the wormholes chain

      **Params**

      > - *royalty:         Royalty, formatted as an integer*
      > - metaURL:       NFT metadata address, formatted as string*
      > - *exchanger:  The exchange when the NFT is minted, the format is a string. When this field is filled, the exchange will exclusively own the NFT. If it is not filled in, no exchange will exclusively own the NFT*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.Mint(10,"/ipfs/ddfd90be9408b4","0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4")
          fmt.Println(rs) //0x8fa2d4b70013407012d002fa395939cb0d322553e4848aaae78d4fad638bef55
      }
      ```

    - ### Mint

      ```
      Mint(royalty uint32, metaURL string, exchanger string) (string, error)
      ```

      Implement NFT minting, create an NFT on the wormholes chain

      **Params**

      > - *royalty:         Royalty, formatted as an integer*
      > - metaURL:       NFT metadata address, formatted as string*
      > - *exchanger:  The exchange when the NFT is minted, the format is a string. When this field is filled, the exchange will exclusively own the NFT. If it is not filled in, no exchange will exclusively own the NFT*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.Mint(10,"/ipfs/ddfd90be9408b4","0xe61e5Bbe724B8F449B5C7BB4a09F99A057253eB4")
          fmt.Println(rs) //0x8fa2d4b70013407012d002fa395939cb0d322553e4848aaae78d4fad638bef55
      }
      ```

    - ### Transfer

      ```
      Transfer(nftAddress, to string) (string, error)
      ```

      Change the ownership of the NFT, and the changed owner owns the NFT

      **Params**

      > - *nftAddress:   NFT address, the format is a decimal string, when it is SNFT, the length can be less than 42 (including 0x), representing the synthesized SNFT*
      > - *to:                  new owner address*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.Transfer("0x0000000000000000000000000000000000000001", "0x814920c33b1a037F91a16B126282155c6F92A10F")
          fmt.Println(rs) //0x5e8dd659b0ceb95ab53ce32d37daa8688accab601ce58c75e706f08bb47617f4
      }
      ```

    - ### Author

      ```
      Author(nftAddress, to string) (string, error)
      ```

      NFT is authorized to an exchange, so that the exchange has the right to sell NFT

      **Params**

      > - *nftAddress:   NFT address, the format is a decimal string, when it is SNFT, the length can be less than 42 (including 0x), representing the synthesized SNFT*
      > - to:                 authorized person*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
           "github.com/wormholes-org/wormholes-client/client"
           "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
      
          rs, _ := worm.NFT.Author("0x0000000000000000000000000000000000000001","0x814920c33b1a037F91a16B126282155c6F92A10F")
          fmt.Println(rs) //0x2657d46f0c4ef16cadbc6842896c1b50f41333d6a247ee43e5085da5d7e3feff
      }
      ```

    - ### AuthorRevoke

      ```
      AuthorRevoke(nftAddress, to string) (string, error)
      ```

      Authorize NFT to an exchange, so that the exchange has the right to sell NFT

      **Prams**

      > - *nftAddress:   NFT address, the format is a decimal string, when it is SNFT, the length can be less than 42 (including 0x), representing the synthesized SNFT*
      > - *to:                 The address of the deauthorizer*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.AuthorRevoke("0x0000000000000000000000000000000000000001", "0x814920c33b1a037F91a16B126282155c6F92A10F")
          fmt.Println(rs) //0xe043dc7d8505d01f6cd949b7a7cc4ed685a9e1b640195801c3c6265b7d11efee
      }
      ```

    - ### AccountAuthor

      ```
      AccountAuthor(to string) (string, error)
      ```

      Authorize all NFTs under the account to an exchange, so that the exchange has the right to sell all NFTs

      **Params**

      > - *to:                Licensee's address*

      **Return**

      `string`  - If the transaction is successful, return the signed byte array; if the signature fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **示例**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.AccountAuthor("0x814920c33b1a037F91a16B126282155c6F92A10F")
          fmt.Println(rs) //0x6b42237b9dad13211d89f1e6c66cf947bb371f407a4621ffcf7fd73e385f6fea
      }
      ```

    - ### AccountAuthorRevoke

      ```
      AccountAuthorRevoke(to string) (string, error)
      ```

      Cancel all authorized NFTs

      **Params**

      > - *to:                The address of the deauthorizer*

      **Return**

      `string`  - If the transaction is successful, return the signed byte array; if the signature fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
      priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
    
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.AccountAuthorRevoke("0x814920c33b1a037F91a16B126282155c6F92A10F")
          fmt.Println(rs) //0x1dee05dff7ea39874ed8401c91288ae627b56ae1df6dc4c26a856fafab0447c5
      }
      ```

    - ### SNFTToERB

      ```
      SNFTToERB(nftAddress string) (string, error)
      ```

      Exchange a certain SNFT fragment mined by the account into ERB

      **Params**

      > - *nftAddress:        The converted snft address, in the format of a decimal string, the length can be less than 42 (including 0x), representing the synthesized SNFT*

      **Return**

      `string`  - If the transaction is successful, return the signed byte array; if the signature fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
         "github.com/wormholes-org/wormholes-client/client"
         "fmt"
      )
      
      const (
      endpoint = "http://192.168.4.237:8574"
      priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
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

      When a user wants to become a miner, he needs to do an ERB pledge transaction first to become a miner with the pledge

      **Return**

      `string`  - If the transaction is successful, return the signed byte array; if the signature fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
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

      When the user does not want to be a miner, or no longer wants to pledge so much ERB, he can do ERB to revoke the pledge

      **Return**

      `string`  - If the transaction is successful, return the signed byte array; if the signature fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
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

      When a user wants to open an exchange, an exchange can be created

      **Params**

      > - *feeRate                  The exchange rate, the format is an integer type*
      > - *name                      Exchange name, formatted as a string*
      > - **url                           Exchange server address, the format is a string*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
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

      The transaction is used to close the exchange

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
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

      This transaction is used to inject NFT fragments that can be mined by miners. Only official accounts can do this transaction

      **Params**

      > - *dir                    The path address where the SNFT is located, the format is a string*
      > - *startIndex        The start number of the SNFT fragment, formatted as a hexadecimal string*
      > - *number           The number of injected SNFT fragments, formatted as a decimal string*
      > - *royalty             Royalty, formatted as an integer*
      > - *creator           Creator, format is a hex string*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
    
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.InsertNFTBlock("wormholes2", "0x640001", 6553600, 20, "0xEaE404DCa7c22A15A59f63002Df54BBb8D90c5FB")
          fmt.Println(rs) //0x61cd018d6e70af47c6204fea18db5b33fdecc92162cca66b0089783733809e84
      }
      ```

    - ### TransactionNFT

      ```
      TransactionNFT(buyer []byte, to string) (string, error)
      ```

      This transaction is used to buy and sell NFTs that have been minted. The transaction originator can be an exchange or a seller

      **Params**

      > - *buyer:               Buyer*
      > - *to                      Buyer address, format is a hexadecimal string*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
    
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "88aaf04596c2c9e71c94c1ec5c160d4326346511b28324d6f19efa9716cb66fd"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.TransactionNFT(buyer,"0x5051B76579BC966A9480dd6E72B39A4C89c1154C")
          fmt.Println(rs) //0xc9c4e6652ba411a0435d2e3187f019329b084734f19ae6699ee7f1fa9a92123b
      }
      ```

    - ### BuyerInitiatingTransaction

      ```
      BuyerInitiatingTransaction(seller1 []byte) (string, error)
      ```

      This transaction is used to buy and sell NFTs that have been minted, and the transaction originator is the buyer

      **Params**

      > - *seller1:             Seller*
      > - *to                      Seller's address, formatted as a hexadecimal string*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.BuyerInitiatingTransaction(seller1, "0x814920c33b1a037F91a16B126282155c6F92A10F")
          fmt.Println(rs) //0xfb9cf0100340c9bf965fc0f8ef44bb8a75af58175deab0dcff3979a97a8ebefa
      }
      ```

    - ### FoundryTradeBuyer

      ```
       FoundryTradeBuyer(seller2 []byte) (string, error)
      ```

      For buying and selling unminted NFTs, the transaction originator is the buyer

      **Params**

      > - *seller2:             Seller*
      > - *to:                    Seller address*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **示例**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
    
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.FoundryTradeBuyer(seller2, "0x814920c33b1a037F91a16B126282155c6F92A10F")
          fmt.Println(rs) //0x4634d6bbc36b9444914a259c2acf0410af0b99122baef30d7a8701a496bc3b6c
      }
      ```

    - ### FoundryExchange

      ```
       FoundryExchange(buyer, seller2 []byte, to string) (string, error)
      ```

      For buying and selling unminted NFTs, the transaction originator is the exchange, or the seller

      **Params**

      > - *buyer:             Buyer*
      > - *seller2:           Seller*
      > - *to:                    Buyer address*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "88aaf04596c2c9e71c94c1ec5c160d4326346511b28324d6f19efa9716cb66fd"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.FoundryTradeBuyer(buyer,seller2 "0x5051B76579BC966A9480dd6E72B39A4C89c1154C")
          fmt.Println(rs) //0x70853466fdf5dc4476fab34b79f9be2e66f0448789937094de0b0aa5f3345e8c
      }
      ```

    - ### NftExchangeMatch

      ```
       NftExchangeMatch(buyer, exchangerAuth []byte, to string) (string, error)
      ```

      It is used to buy and sell NFTs that have been minted. The transaction originator is the exchange. This transaction is used when exchange A authorizes another exchange B, and exchange B initiates the transaction.

      **Params**

      > - *buyer                         Buyer*
      > - *exchangerAuth:       Authorized Exchange A*
      > - *to                               buyer, the format is a hexadecimal string*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "abdab4bd87d3ea117b5512d6ce28522f7b9421511d8bb08bb20277dde6fb8320"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.NftExchangeMatch(buyer,exchangAuth,"0x5051B76579BC966A9480dd6E72B39A4C89c1154C")
          fmt.Println(rs) //0xf11e024297b89e6dfd02bc2da4680cea353ea6956c3ea9084afa40d58477932f
      }
      ```

    - ### FoundryExchangeInitiated

      ```
       FoundryExchangeInitiated(buyer, seller2, exchangerAuthor []byte, to string) (string, error)
      ```

      It is used to buy and sell unminted NFTs. The transaction originator is the exchange. The transaction is used when exchange A authorizes another exchange B, and exchange B initiates the transaction.

      **Params**

      > - *buyer                         Buyer*
      > - *Seller2                       Seller*
      > - *exchangerAuth:       Authorized Exchange A*
      > - *to                               Buyer, the format is a hexadecimal string*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "abdab4bd87d3ea117b5512d6ce28522f7b9421511d8bb08bb20277dde6fb8320"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.NftExchangeMatch(buyer,seller2,exchangAuth,  "0x5051B76579BC966A9480dd6E72B39A4C89c1154C")
          fmt.Println(rs) //0xc9cc570057faf1edd83f48833520f9d546e4972083ee705152b5f35630f1588d
      }
      ```

    - ### FtDoesNotAuthorizeExchanges

      ```
       FtDoesNotAuthorizeExchanges(buyer, seller1 []byte, to string) (string, error)
      ```

      Used to buy and sell NFTs that have been minted, the transaction originator is the exchange, and the transaction is used when the NFT is not authorized to the exchange

      **Params**

      > - *buyer               Buyer*
      > - *seller1:             Seller*
      > - *to                     Buyer address, format is a hexadecimal string*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "11e4259f98e6a18772be5e1b2e2c9e5b12b4a9fe8e3cfa0853df59fa0825e861"
      )
      
      func main() {
          worm := client.NewClient(priKey, endpoint)
          rs, _ := worm.NFT.FtDoesNotAuthorizeExchanges(buyer,seller1,"0x5051B76579BC966A9480dd6E72B39A4C89c1154C")
          fmt.Println(rs) //0x95615a6c7a164537257492c112a9fcd99907315893706a1b104456d9e3aa8af6
      }
      ```

    - ### AdditionalPledgeAmount

      ```
       AdditionalPledgeAmount(value int64) (string, error)
      ```

      The amount used by the exchange to increase the pledged ERB

      **Params**

      > - Values:         additional amount

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
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

      The amount used by the exchange to increase the pledged ERB

      **Params**

      > - Values:          reduced amount

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
          "github.com/wormholes-org/wormholes-client/client"
          "fmt"
      )
      
      const (
          endpoint = "http://192.168.4.237:8574"
          priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
      
      func main() {
          worm := client.NewNFT(priKey, endpoint)
          rs, _ := worm.NFT.RevokesPledgeAmount(100)
          fmt.Println(rs) //0xd2c7f943f0f5364b0928c518e7b6de7491c0e8efb6abf912a17e6860f70ebec1
      }
      ```
    - ### VoteOfficialNFT

      ```
      VoteOfficialNFT(dir, startIndex string, number uint64, royalty uint32, creator string) (string, error)
      ```

      This transaction is used to inject NFT fragments that can be mined by miners. Only official accounts can do this transaction

      **Params**

      > - *dir                    The path address where the SNFT is located, the format is a string*
      > - *startIndex        The start number of the SNFT fragment, formatted as a hexadecimal string*
      > - *number           The number of injected SNFT fragments, formatted as a decimal string*
      > - *royalty             Royalty, formatted as an integer*
      > - *creator           Creator, format is a hex string*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
        "github.com/wormholes-org/wormholes-client/client"
        "fmt"
      )
      
      const (
        endpoint = "http://192.168.4.237:8574"
        priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
    
      func main() {
        worm := client.NewClient(priKey, endpoint)
        rs, _ := worm.NFT.VoteOfficialNFT("wormholes2", "0x640001", 6553600, 20, "0xEaE404DCa7c22A15A59f63002Df54BBb8D90c5FB")
        fmt.Println(rs) //0x61cd018d6e70af47c6204fea18db5b33fdecc92162cca66b0089783733809e84
      }
      ```
    - ### VoteOfficialNFTByApprovedExchanger
        
      ```
      VoteOfficialNFTByApprovedExchanger(dir, startIndex string, number uint64, royalty uint32, creator string, exchangerAuth []byte) (string, error)
      ```

      This transaction is used to inject NFT fragments that can be mined by miners. Only official accounts can do this transaction

      **Params**

      > - *dir                    The path address where the SNFT is located, the format is a string*
      > - *startIndex        The start number of the SNFT fragment, formatted as a hexadecimal string*
      > - *number           The number of injected SNFT fragments, formatted as a decimal string*
      > - *royalty             Royalty, formatted as an integer*
      > - *creator           Creator, format is a hex string*
      > - *exchangerAuth           exchangerAuth*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
           "github.com/wormholes-org/wormholes-client/client"
           "fmt"
      )
      
      const (
           endpoint = "http://192.168.4.237:8574"
           priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
    
      func main() {
           worm := client.NewClient(priKey, endpoint)
           rs, _ := worm.NFT.VoteOfficialNFTByApprovedExchanger("wormholes2", "0x640001", 6553600, 20, "0xab7624f47fd7dadb6b8e255d06a2f10af55990fe", exchangeAuth)
           fmt.Println(rs)
      }
      ```

    - ### ChangeRewardsType

      ```
      ChangeRewardsType() (string, error)
      ```

      This transaction is used to change rewards

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
         "github.com/wormholes-org/wormholes-client/client"
         "fmt"
      )
      
      const (
         endpoint = "http://192.168.4.237:8574"
         priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
    
      func main() {
         worm := client.NewClient(priKey, endpoint)
         rs, _ := worm.NFT.ChangeRewardsType()
         fmt.Println(rs)
      }
      ```

    - ### AccountDelegate

      ```
      AccountDelegate(proxyAddress string) (string, error)
      ```

      This transaction is used to change rewards

      **Params**

      > - *proxyAddress                    proxy address for delegation*

      **Return**

      `string`  - If the transaction is successful, return the hash of the transaction; if the transaction fails, return nil

      `error`   - If the transaction is successful, return nil; if the transaction fails, return the corresponding error

      **Example**

      ```
      package main
      import (
       "github.com/wormholes-org/wormholes-client/client"
       "fmt"
      )
      
      const (
       endpoint = "http://192.168.4.237:8574"
       priKey   = "b2ebd0889351eb22dc73c3a02c63e783794a9de3f578d6d07bb370cc112d2ec7"
      )
    
      func main() {
       worm := client.NewClient(priKey, endpoint)
       rs, _ := worm.NFT.AccountDelegate("0x814920c33b1a037F91a16B126282155c6F92A10F")
       fmt.Println(rs)
      }
      ```