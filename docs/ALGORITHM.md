# Finality Quorum Calculation Algorithm

## Table of Contents

1. [Overview](#1-overview)
2. [Inputs and Data Sources](#2-inputs-and-data-sources)
3. [Step-by-Step Algorithm](#3-step-by-step-algorithm)
4. [Code Snippets](#4-code-snippets)
5. [Error Handling and Edge Cases](#5-error-handling-and-edge-cases)
6. [References](#6-references)

---

## 1. Overview

The Finality Gadget determines if a rollup block is finalized by checking 
if a quorum of voting power has signed the block. This involves fetching votes from the contract, 
calculating each provider's voting power, and comparing the total voted power 
to the required threshold (>2/3).

## 2. Inputs and Data Sources

- **Block Info:**
  - Block height, hash, and timestamp from the L2 chain.
- **Finality Votes:**
  - List of FP public keys who voted for the block, fetched from the Rollup BSN CosmWasm contract.
- **FP Voting Power:**
  - Voting power for each FP, fetched from the Babylon chain at the corresponding BTC block height.
- **Quorum Threshold:**
  - Hardcoded as 2/3 of total voting power.

## 3. Step-by-Step Algorithm

1. **Check if the finality gadget is enabled:**
   - Query the contract for the enabled flag.
2. **Fetch all FP public keys for the consumer chain:**
   - Query Babylon for the list of FPs.
3. **Convert the L2 block timestamp to BTC block height:**
   - Use the Bitcoin client to map timestamp to BTC height.
4. **Fetch voting power for each FP at this BTC height:**
   - Query Babylon for each FP's active delegations and sum their power.
5. **Calculate total voting power:**
   - Sum the voting power of all FPs.
6. **Fetch the list of FPs who voted for the block:**
   - Query the contract for votes for the given block height and hash.
7. **Calculate voted voting power:**
   - Sum the voting power of FPs who voted.
8. **Check if voted power meets quorum:**
   - If `votedPower * 3 >= totalPower * 2`, the block is finalized.

> **Note:**
> The current implementation does **not** track or take into account the allow-list in 
> the rollup-bsn contract. It only considers Finality Providers (FPs) registered in the 
> Babylon chain for the given Consumer ID. This means that FPs who are not on the allow-list 
> but are registered in Babylon may be counted, and FPs who are on the allow-list but not 
> registered in Babylon will not be counted.

## 4. Code Snippets

### Main Finality Check (Go)
```go
func (fg *FinalityGadget) QueryIsBlockBabylonFinalizedFromBabylon(block *types.Block) (bool, error) {
    // 1. Check if enabled
    isEnabled, err := fg.cwClient.QueryIsEnabled()
    if err != nil || !isEnabled {
        return isEnabled, err
    }

    // 2. Get all FPs
    allFpPks, err := fg.queryAllFpBtcPubKeys()
    if err != nil {
        return false, err
    }

    // 3. Convert timestamp to BTC height
    btcblockHeight, err := fg.btcClient.GetBlockHeightByTimestamp(block.BlockTimestamp)
    if err != nil {
        return false, err
    }

    // 4. Get FP voting power
    allFpPower, err := fg.bbnClient.QueryMultiFpPower(allFpPks, btcblockHeight)
    if err != nil {
        return false, err
    }

    // 5. Calculate total power
    var totalPower uint64
    for _, power := range allFpPower {
        totalPower += power
    }
    if totalPower == 0 {
        return false, types.ErrNoFpHasVotingPower
    }

    // 6. Get FPs who voted
    votedFpPks, err := fg.cwClient.QueryListOfVotedFinalityProviders(block)
    if err != nil || votedFpPks == nil {
        return false, err
    }

    // 7. Calculate voted power
    var votedPower uint64
    for _, key := range votedFpPks {
        if power, exists := allFpPower[key]; exists {
            votedPower += power
        }
    }

    // 8. Check quorum (2/3)
    if votedPower*3 < totalPower*2 {
        return false, nil
    }
    return true, nil
}
```

### Fetching Votes from Contract (Go)
```go
func (cwClient *CosmWasmClient) QueryListOfVotedFinalityProviders(block *types.Block) ([]string, error) {
    // ...construct query...
    resp, err := cwClient.querySmartContractState(queryData)
    // ...unmarshal response...
    return votedFpPkHexList, nil
}
```

### Fetching FP Voting Power (Go)
```go
func (bbnClient *BabylonClient) QueryMultiFpPower(fpPubkeys []string, btcHeight uint32) (map[string]uint64, error) {
    // ...iterate over FPs, sum active delegations...
    return fpPowerMap, nil
}
```

## 5. Error Handling and Edge Cases
- If no FPs have voting power, the block cannot be finalized.
- If the contract is not enabled, all blocks are considered finalized (for testing or pass-through mode).
- If the BTC staking is not activated, finality is not possible.

## 6. References
- [finalitygadget/finalitygadget.go](../finalitygadget/finalitygadget.go)
- [cwclient/cwclient.go](../cwclient/cwclient.go)
- [bbnclient/bbnclient.go](../bbnclient/bbnclient.go)
- [proto/finalitygadget.proto](../proto/finalitygadget.proto) 