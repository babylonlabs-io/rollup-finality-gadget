package cwclient

import (
	"context"
	"encoding/json"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/babylonlabs-io/finality-gadget/types"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
)

type CosmWasmClient struct {
	rpcclient.Client
	contractAddr string
}

const (
	// hardcode the timeout to 20 seconds. We can expose it to the params once needed
	DefaultTimeout = 20 * time.Second
)

//////////////////////////////
// CONSTRUCTOR
//////////////////////////////

func NewCosmWasmClient(rpcClient rpcclient.Client, contractAddr string) *CosmWasmClient {
	return &CosmWasmClient{
		Client:       rpcClient,
		contractAddr: contractAddr,
	}
}

//////////////////////////////
// METHODS
//////////////////////////////

func (cwClient *CosmWasmClient) QueryListOfVotedFinalityProviders(
	queryParams *types.Block,
) ([]string, error) {
	queryData, err := createBlockVotersQueryData(queryParams)
	if err != nil {
		return nil, err
	}

	resp, err := cwClient.querySmartContractState(queryData)
	if err != nil {
		return nil, err
	}
	// BlockVoters's return type is Option<HashSet<String>> in contract
	// Check empty response before unmarshaling
	if len(resp.Data) == 0 {
		return nil, nil
	}

	// The contract now returns detailed finality signature information
	// We need to extract just the fp_btc_pk_hex field from each object
	var votedFpDetails []struct {
		FpBtcPkHex string `json:"fp_btc_pk_hex"`
		// We can ignore the other fields (pub_rand, finality_signature) for now
	}

	if err := json.Unmarshal(resp.Data, &votedFpDetails); err != nil {
		return nil, err
	}

	// Extract just the public key hex values
	votedFpPkHexList := make([]string, len(votedFpDetails))
	for i, detail := range votedFpDetails {
		votedFpPkHexList[i] = detail.FpBtcPkHex
	}

	return votedFpPkHexList, nil
}

func (cwClient *CosmWasmClient) QueryConsumerId() (string, error) {
	queryData, err := createConfigQueryData()
	if err != nil {
		return "", err
	}

	resp, err := cwClient.querySmartContractState(queryData)
	if err != nil {
		return "", err
	}

	var data contractConfigResponse
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}

	return data.ConsumerId, nil
}

func (cwClient *CosmWasmClient) QueryConfig() (*types.ContractConfig, error) {
	queryData, err := createConfigQueryData()
	if err != nil {
		return nil, err
	}

	resp, err := cwClient.querySmartContractState(queryData)
	if err != nil {
		return nil, err
	}

	var data contractConfigResponse
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}

	return &types.ContractConfig{
		ConsumerId:                data.ConsumerId,
		BsnActivationHeight:       data.BsnActivationHeight,
		FinalitySignatureInterval: data.FinalitySignatureInterval,
		MinPubRand:                data.MinPubRand,
	}, nil
}

//////////////////////////////
// INTERNAL
//////////////////////////////

func createBlockVotersQueryData(queryParams *types.Block) ([]byte, error) {
	queryData := ContractQueryMsgs{
		BlockVoters: &blockVotersQuery{
			Height: queryParams.BlockHeight,
			Hash:   queryParams.BlockHash,
		},
	}
	data, err := json.Marshal(queryData)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type contractConfigResponse struct {
	ConsumerId                string `json:"bsn_id"`
	BsnActivationHeight       uint64 `json:"bsn_activation_height"`
	FinalitySignatureInterval uint64 `json:"finality_signature_interval"`
	MinPubRand                uint64 `json:"min_pub_rand"`
}

type ContractQueryMsgs struct {
	Config      *contractConfig   `json:"config,omitempty"`
	BlockVoters *blockVotersQuery `json:"block_voters,omitempty"`
}

type blockVotersQuery struct {
	Hash   string `json:"hash_hex"`
	Height uint64 `json:"height"`
}

type contractConfig struct{}

func createConfigQueryData() ([]byte, error) {
	queryData := ContractQueryMsgs{
		Config: &contractConfig{},
	}
	data, err := json.Marshal(queryData)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// querySmartContractState queries the smart contract state given the contract address and query data
func (cwClient *CosmWasmClient) querySmartContractState(
	queryData []byte,
) (*wasmtypes.QuerySmartContractStateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	sdkClientCtx := cosmosclient.Context{Client: cwClient.Client}
	wasmQueryClient := wasmtypes.NewQueryClient(sdkClientCtx)

	req := &wasmtypes.QuerySmartContractStateRequest{
		Address:   cwClient.contractAddr,
		QueryData: queryData,
	}
	return wasmQueryClient.SmartContractState(ctx, req)
}
