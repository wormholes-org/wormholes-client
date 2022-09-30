package types

const WormHolesVersion = "v0.0.1"

const (
	Mint = iota
	Transfer
	Author
	AuthorRevoke
	AccountAuthor
	AccountAuthorRevoke
	SNFTToERB
	SNFTPledge
	SNFTRevokesPledge
	TokenPledge
	TokenRevokesPledge
	Open
	Close
	_
	TransactionNFT
	BuyerInitiatingTransaction
	FoundryTradeBuyer
	FoundryExchange
	NftExchangeMatch
	FoundryExchangeInitiated
	FtDoesNotAuthorizeExchanges
	AdditionalPledgeAmount
	RevokesPledgeAmount
	VoteOfficialNFT
	VoteOfficialNFTByApprovedExchanger
	UnforzenAccount
	_
	_
	_
	_
	_
	AccountDelegate
)

// Transaction struct for handling NFT transactions
type Transaction struct {
	Type       uint8  `json:"type"`
	Dir        string `json:"dir,omitempty"`
	StartIndex string `json:"start_index,omitempty"`
	Number     uint64 `json:"number,omitempty"`
	NFTAddress string `json:"nft_address,omitempty"`
	Royalty    uint32 `json:"royalty,omitempty"`
	MetaURL    string `json:"meta_url,omitempty"`
	Exchanger  string `json:"exchanger,omitempty"`
	//ApproveAddress string		`json:"approve_address"`
	FeeRate       uint32         `json:"fee_rate,omitempty"`
	Name          string         `json:"name,omitempty"`
	Url           string         `json:"url,omitempty"`
	Buyer         *Buyer         `json:"buyer,omitempty"`
	Seller1       *Seller1       `json:"seller1,omitempty"`
	Seller2       *Seller2       `json:"seller2,omitempty"`
	ExchangerAuth *ExchangerAuth `json:"exchanger_auth,omitempty"`
	Creator       string         `json:"creator,omitempty"`
	RewardFlag    int            `json:"reward_flag,omitempty"`
	ProxyAddress  string         `json:"proxy_address,omitempty"`
	ProxySign     string         `json:"proxy_sign,omitempty"`
	Version       string         `json:"version"`
}

type Buyer struct {
	Amount      string `json:"price,omitempty"`
	NFTAddress  string `json:"nft_address,omitempty"`
	Exchanger   string `json:"exchanger,omitempty"`
	BlockNumber string `json:"block_number,omitempty"`
	Seller      string `json:"seller,omitempty"`
	Sig         string `json:"sig,omitempty"`
}

type Seller1 struct {
	Amount      string `json:"price,omitempty"`
	NFTAddress  string `json:"nft_address,omitempty"`
	Exchanger   string `json:"exchanger,omitempty"`
	BlockNumber string `json:"block_number,omitempty"`
	Sig         string `json:"sig,omitempty"`
}

type Seller2 struct {
	Amount        string `json:"price,omitempty"`
	Royalty       string `json:"royalty,omitempty"`
	MetaURL       string `json:"meta_url,omitempty"`
	ExclusiveFlag string `json:"exclusive_flag,omitempty"`
	Exchanger     string `json:"exchanger,omitempty"`
	BlockNumber   string `json:"block_number,omitempty"`
	Sig           string `json:"sig,omitempty"`
}

type ExchangerAuth struct {
	ExchangerOwner string `json:"exchanger_owner,omitempty"`
	To             string `json:"to,omitempty"`
	BlockNumber    string `json:"block_number,omitempty"`
	Sig            string `json:"sig,omitempty"`
}
