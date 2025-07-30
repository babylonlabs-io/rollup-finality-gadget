package types

type Block struct {
	BlockHash      string `json:"block_hash" description:"block hash"`
	BlockHeight    uint64 `json:"block_height" description:"block height"`
	BlockTimestamp uint64 `json:"block_timestamp" description:"block timestamp"`
}

type ChainSyncStatus struct {
	LatestBlockHeight               uint64 `json:"latest_block"`
	LatestBtcFinalizedBlockHeight   uint64 `json:"latest_btc_finalized_block"`
	EarliestBtcFinalizedBlockHeight uint64 `json:"earliest_btc_finalized_block"`
	LatestEthFinalizedBlockHeight   uint64 `json:"latest_eth_finalized_block"`
}

type ContractConfig struct {
	ConsumerId                string `json:"bsn_id"`
	BsnActivationHeight       uint64 `json:"bsn_activation_height"`
	FinalitySignatureInterval uint64 `json:"finality_signature_interval"`
	MinPubRand                uint64 `json:"min_pub_rand"`
}
