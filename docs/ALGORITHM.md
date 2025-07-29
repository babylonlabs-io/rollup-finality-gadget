# Finality Quorum Calculation Algorithm

## Table of Contents

1. [Overview](#1-overview)
2. [Data Sources](#2-data-sources)  
3. [Inputs](#3-inputs)  
4. [Step‑by‑Step Algorithm](#4-step-by-step-algorithm)

---

## 1. Overview

The Finality Gadget determines if a rollup block is finalized by 
checking if a quorum of voting power has signed the block. This 
involves fetching votes from the contract, calculating each 
provider’s voting power, and comparing the total voted power to 
the required threshold (> 2/3).

---

## 2. Data Sources

- **Babylon Genesis chain**  
  The Rollup BSN contract is deployed on Babylon Genesis. From it we fetch: 
  - signature submissions 
  - public randomness commitments 
  - configuration parameters  
  Babylon Genesis also tracks registered Finality Providers for Rollup BSN.

- **Rollup full node**  
  We query rollup blocks by height and extract each block’s hash and timestamp.

- **Bitcoin full node**  
  We retrieve BTC delegation data at specific Bitcoin block heights, 
  which we use to calculate voting power and assemble the voting table.  


## 3. Inputs

The gadget runs this algorithm for every rollup block:
```go
   func QueryIsBlockBabylonFinalizedFromBabylon(block *types.Block) (bool, error)
```

Function Inputs:  
- single Rollup Block (height, hash, timestamp)  

All other data is fetched within the function:  
- Registered Finality Providers and their voting power from Babylon Genesis 
- Quorum threshold (2/3 of total voting power) 
---

## 4. Step‑by‑Step Algorithm

1. **Fetch activation parameters**  
   Read `bsn_activation_height` and `finality_signature_interval` from the Rollup BSN contract

2. **Determine start block**  
   ```pseudo
   if currentBlockHeight < bsn_activation_height:
       currentBlockHeight = bsn_activation_height
   ```

3. **Align to signature interval**  
   Advance to the next block that satisfies  
   ```pseudo
   (currentBlockHeight - bsn_activation_height) % finality_signature_interval == 0
   ```
4. **Map timestamp → Babylon height**  
   Use the Babylon client to locate the Babylon block whose timestamp  
     matches the L2 block’s timestamp

4. **Fetch FP public keys**  
   Get the list of all registered FP's (validator public keys) from Babylon

5. **Fetch and filter FP public keys**  
   - Query Babylon Genesis at that height for registered FP public keys  
   - Query the Rollup BSN contract’s allow‑list at that height  
   - Compare both lists and use their intersection as the effective FP set 

6. **Map timestamp → BTC height**  
   Use your Bitcoin client to find the BTC block whose timestamp matches the L2 block’s timestamp

7. **Compute each FP’s voting power**  
   For each FP:  
   1. If the FP has no timestamped pub rand at this height, set 
      votingPower[FP] = 0 and move to the next FP
   2. Query Babylon for active delegations at that BTC height 
   3. Sum up shares → `votingPower[FP]`

8. **Calculate total voting power**  
   For all FPs:
   ```pseudo
   totalPower = Σ votingPower[FP]
   ```

9. **Fetch votes for the block**  
   Read from the Rollup BSN contract which FPs voted on this block

10. **Sum voted power**  
   For all FPs in voteList:
      ```pseudo
      votedPower = Σ votingPower[FP]
      ```

11. **Check quorum**  
    ```pseudo
    if votedPower * 3 >= totalPower * 2:
        mark block as FINALIZED
        go-to next valid block
    else:
        block remains UNFINALIZED
        wait for block to be finalized
    ```

> **Note:**
> - The current implementation does **not** enforce the contract’s allow‑list.  
>   It only considers FPs registered on Babylon for the given Consumer ID at  
>   the latest block.
> - It does **not** track FP public keys historically; it only uses the latest  
>   active set.
> - It does **not** verify whether an FP submitted public randomness before  
>   its signature.