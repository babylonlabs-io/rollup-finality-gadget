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

1. **Fetch activation params from contract**  
   - Query the Rollup BSN contract for `bsn_activation_height` and 
     `finality_signature_interval`

2. **Check block eligibility**  
   - If `currentBlockHeight < bsn_activation_height`, **skip**: set starting
     block to `bsn_activation_height`
   - Else if  
     `(currentBlockHeight - bsn_activation_height) % 
     finality_signature_interval != 0`, **skip**: move to next valid block

3. **Fetch all FP public keys**  
   - Query Babylon for the list of registered FPs for this consumer

4. **Convert block timestamp to BTC height**  
   - Use the Bitcoin client to map the L2 block’s timestamp to the BTC 
     block height

5. **Fetch voting power for each FP**  
   - For each FP:  
     1. Query Babylon for their active delegations at that BTC height
     2. Sum the delegation shares to compute the FP’s voting power

6. **Calculate total voting power**  
   - Sum the voting power of all FPs fetched

7. **Fetch votes for the block**  
   - Query the Rollup BSN contract for votes on the given block height 
     and hash

8. **Calculate voted voting power**  
   - Sum the voting power of only those FPs present in the vote list

9. **Check if quorum is met**  
   - If `votedPower * 3 >= totalPower * 2`, **mark block as finalized**
   - Otherwise, **block remains unfinalized**

> **Note:**  
> The current implementation does **not** enforce the contract’s allow‑list. 
> It only considers FPs registered on Babylon for the given Consumer ID, 
> which may slightly diverge from the on‑chain allow‑list.
