package bbnclient

import (
	"math"

	"github.com/babylonlabs-io/babylon/v3/client/query"
	bbntypes "github.com/babylonlabs-io/babylon/v3/x/btcstaking/types"
	sdkquerytypes "github.com/cosmos/cosmos-sdk/types/query"
)

type BabylonClient struct {
	*query.QueryClient
}

//////////////////////////////
// CONSTRUCTOR
//////////////////////////////

func NewBabylonClient(queryClient *query.QueryClient) *BabylonClient {
	return &BabylonClient{
		QueryClient: queryClient,
	}
}

//////////////////////////////
// METHODS
//////////////////////////////

func (bbnClient *BabylonClient) QueryAllFpBtcPubKeys(consumerId string) ([]string, error) {
	pagination := &sdkquerytypes.PageRequest{}
	resp, err := bbnClient.QueryClient.FinalityProviders(consumerId, pagination)
	if err != nil {
		return nil, err
	}

	var pkArr []string

	for _, fp := range resp.FinalityProviders {
		pkArr = append(pkArr, fp.BtcPk.MarshalHex())
	}
	return pkArr, nil
}

func (bbnClient *BabylonClient) QueryMultiFpPower(
	fpPubkeyHexList []string,
	btcHeight uint32,
) (map[string]uint64, error) {
	// Pre-fetch parameters once for all FPs (they're the same for all delegations at this height)
	btccheckpointParams, err := bbnClient.QueryClient.BTCCheckpointParams()
	if err != nil {
		return nil, err
	}
	btcstakingParams, err := bbnClient.QueryClient.BTCStakingParams()
	if err != nil {
		return nil, err
	}

	// Process FPs in parallel for better performance
	type fpResult struct {
		fpPubkeyHex string
		power       uint64
		err         error
	}

	resultCh := make(chan fpResult, len(fpPubkeyHexList))

	// Extract values once
	kValue := btccheckpointParams.GetParams().BtcConfirmationDepth
	wValue := btccheckpointParams.GetParams().CheckpointFinalizationTimeout
	covQuorum := btcstakingParams.GetParams().CovenantQuorum

	// Launch goroutines for parallel processing
	for _, fpPubkeyHex := range fpPubkeyHexList {
		go func(fpPk string) {
			power, err := bbnClient.queryFpPower(fpPk, btcHeight, kValue, wValue, covQuorum)
			resultCh <- fpResult{fpPubkeyHex: fpPk, power: power, err: err}
		}(fpPubkeyHex)
	}

	// Collect results
	fpPowerMap := make(map[string]uint64)
	for i := 0; i < len(fpPubkeyHexList); i++ {
		result := <-resultCh
		if result.err != nil {
			return nil, result.err
		}
		fpPowerMap[result.fpPubkeyHex] = result.power
	}

	return fpPowerMap, nil
}

// QueryEarliestActiveDelBtcHeight returns the earliest active BTC staking height
func (bbnClient *BabylonClient) QueryEarliestActiveDelBtcHeight(fpPkHexList []string) (uint32, error) {
	allFpEarliestDelBtcHeight := uint32(math.MaxUint32)

	for _, fpPkHex := range fpPkHexList {
		fpEarliestDelBtcHeight, err := bbnClient.QueryFpEarliestActiveDelBtcHeight(fpPkHex)

		if err != nil {
			return math.MaxUint32, err
		}
		if fpEarliestDelBtcHeight < allFpEarliestDelBtcHeight {
			allFpEarliestDelBtcHeight = fpEarliestDelBtcHeight
		}
	}

	return allFpEarliestDelBtcHeight, nil
}

func (bbnClient *BabylonClient) QueryFpEarliestActiveDelBtcHeight(fpPubkeyHex string) (uint32, error) {
	pagination := &sdkquerytypes.PageRequest{
		Limit: 100,
	}

	// queries the BTCStaking module for all delegations of a finality provider
	resp, err := bbnClient.QueryClient.FinalityProviderDelegations(fpPubkeyHex, pagination)
	if err != nil {
		return math.MaxUint32, err
	}

	// queries BtcConfirmationDepth, CovenantQuorum, and the latest BTC header
	btccheckpointParams, err := bbnClient.QueryClient.BTCCheckpointParams()
	if err != nil {
		return math.MaxUint32, err
	}

	// get the BTC staking params
	btcstakingParams, err := bbnClient.QueryClient.BTCStakingParams()
	if err != nil {
		return math.MaxUint32, err
	}

	// get the latest BTC header
	btcHeader, err := bbnClient.QueryClient.BTCHeaderChainTip()
	if err != nil {
		return math.MaxUint32, err
	}

	kValue := btccheckpointParams.GetParams().BtcConfirmationDepth
	covQuorum := btcstakingParams.GetParams().CovenantQuorum
	latestBtcHeight := btcHeader.GetHeader().Height

	earliestBtcHeight := uint32(math.MaxUint32)
	pageCount := 0
	maxPages := 500 // Safety limit to prevent infinite loops

	for {
		// btcDels contains all the queried BTC delegations
		for _, btcDels := range resp.BtcDelegatorDelegations {
			for _, btcDel := range btcDels.Dels {
				activationHeight := getDelFirstActiveHeight(btcDel, latestBtcHeight, kValue, covQuorum)
				if activationHeight < earliestBtcHeight {
					earliestBtcHeight = activationHeight
				}
			}
		}

		if resp.Pagination == nil || resp.Pagination.NextKey == nil {
			break
		}

		pageCount++
		if pageCount >= maxPages {
			break
		}

		// Set up pagination for next query
		pagination.Key = resp.Pagination.NextKey

		resp, err = bbnClient.QueryClient.FinalityProviderDelegations(fpPubkeyHex, pagination)
		if err != nil {
			return math.MaxUint32, err
		}
	}
	return earliestBtcHeight, nil
}

//////////////////////////////
// INTERNAL
//////////////////////////////

// queryFpPower is an optimized version that reuses cached parameters
func (bbnClient *BabylonClient) queryFpPower(
	fpPubkeyHex string,
	btcHeight uint32,
	kValue uint32,
	wValue uint32,
	covQuorum uint32,
) (uint64, error) {
	totalPower := uint64(0)
	pagination := &sdkquerytypes.PageRequest{}

	// Query delegations for this FP
	resp, err := bbnClient.QueryClient.FinalityProviderDelegations(fpPubkeyHex, pagination)
	if err != nil {
		return 0, err
	}

	for {
		// Process all delegations in this page
		for _, btcDels := range resp.BtcDelegatorDelegations {
			for _, btcDel := range btcDels.Dels {
				// Check if delegation is active using cached parameters (inline for performance)
				if bbnClient.isDelegationActive(btcDel, btcHeight, kValue, wValue, covQuorum) {
					totalPower += btcDel.TotalSat
				}
			}
		}

		if resp.Pagination == nil || resp.Pagination.NextKey == nil {
			break
		}
		pagination.Key = resp.Pagination.NextKey

		// Query next page
		resp, err = bbnClient.QueryClient.FinalityProviderDelegations(fpPubkeyHex, pagination)
		if err != nil {
			return 0, err
		}
	}

	return totalPower, nil
}

// isDelegationActive checks if delegation is active using pre-calculated parameters
func (bbnClient *BabylonClient) isDelegationActive(
	btcDel *bbntypes.BTCDelegationResponse,
	btcHeight uint32,
	kValue uint32,
	wValue uint32,
	covQuorum uint32,
) bool {
	ud := btcDel.UndelegationResponse

	if ud.DelegatorUnbondingInfoResponse != nil {
		return false
	}

	// k is not involved in the `GetStatus` logic as Babylon will accept a BTC delegation request
	// only when staking tx is k-deep on BTC.
	//
	// But the msg handler performs both checks 1) ensure staking tx is k-deep, and 2) ensure the
	// staking tx's timelock has at least w BTC blocks left.
	// (https://github.com/babylonlabs-io/babylon-private/blob/3d8f190c9b0c0795f6546806e3b8582de716cd60/x/btcstaking/keeper/msg_server.go#L283-L292)
	//
	// So after the msg handler accepts BTC delegation the 1st check is no longer needed
	// the k-value check is added per
	//
	// So in our case, we need to check both to ensure the delegation is active
	if btcHeight < btcDel.StartHeight+kValue || btcHeight+wValue > btcDel.EndHeight {
		return false
	}

	if len(btcDel.CovenantSigs) < int(covQuorum) {
		return false
	}
	if len(ud.CovenantUnbondingSigList) < int(covQuorum) {
		return false
	}
	if len(ud.CovenantSlashingSigs) < int(covQuorum) {
		return false
	}

	return true
}

// The active delegation needs to satisfy:
// 1) the staking tx is k-deep in Bitcoin, i.e., start_height + k
// 2) it receives a quorum number of covenant committee signatures
//
// return math.MaxUint32 if the delegation is not active
//
// Note: the delegation can be unbounded and that's totally fine and shouldn't affect when the chain was activated
func getDelFirstActiveHeight(btcDel *bbntypes.BTCDelegationResponse, latestBtcHeight, kValue, covQuorum uint32) uint32 {
	activationHeight := btcDel.StartHeight + kValue
	// not activated yet
	if latestBtcHeight < activationHeight || len(btcDel.CovenantSigs) < int(covQuorum) {
		return math.MaxUint32
	}
	return activationHeight
}
