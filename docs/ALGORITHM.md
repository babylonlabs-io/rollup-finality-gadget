# Finality Quorum Calculation Algorithm

## Table of Contents

1. [Overview](#1-overview)  
2. [Inputs and Data Sources](#2-inputs-and-data-sources)  
3. [Step‑by‑Step Algorithm](#3-step-by-step-algorithm)

---

## 1. Overview

The Finality Gadget determines if a rollup block is finalized by 
checking if a quorum of voting power has signed the block. This 
involves fetching votes from the contract, calculating each 
provider’s voting power, and comparing the total voted power to 
the required threshold (> 2/3).

---

## 2. Inputs and Data Sources

1. Rollup blocks  
   - Fetch height, hash, and timestamp from the rollup chain  
2. Finality Providers  
   - Fetch FP public keys from the Babylon chain  
3. Voting Power  
   - Calculate FP voting power via Bitcoin delegation data at the 
     specific BTC block height  
4. Quorum Threshold  
   - Hardcoded as 2/3 of total voting power  

---

## 3. Step‑by‑Step Algorithm

1. **Fetch activation parameters**  
   Read `bsn_activation_height` and `finality_signature_interval` from the Rollup BSN contract.

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

4. **Fetch validator public keys**  
   Get the list of all registered FPs (validator public keys) from Babylon.

5. **Map timestamp → BTC height**  
   Use your Bitcoin client to find the BTC block whose timestamp matches the L2 block’s timestamp.

6. **Compute each validator’s voting power**  
   For each FP:  
   1. Query Babylon for active delegations at that BTC height.  
   2. Sum up shares → `votingPower[FP]`

7. **Calculate total voting power**  
   For all FPs:
   ```pseudo
   totalPower = Σ votingPower[FP]
   ```

8. **Fetch votes for the block**  
   Read from the Rollup BSN contract which FPs voted on this block.

9. **Sum voted power**  
   For all FPs in voteList:
   ```pseudo
   votedPower = Σ votingPower[FP]
   ```

10. **Check quorum**  
    ```pseudo
    if votedPower * 3 >= totalPower * 2:
        mark block as FINALIZED
        go-to next valid block
    else:
        block remains UNFINALIZED
        wait for block to be finalized
    ```

> **Note:**  
> The current implementation does **not** enforce the contract’s allow‑list. 
> It only considers FPs registered on Babylon for the given Consumer ID, 
> which may slightly diverge from the on‑chain allow‑list.
