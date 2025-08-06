# Finality Gadget Metrics

This document describes the Prometheus metrics exposed by the finality gadget service.

## Endpoint

Metrics are exposed at `/metrics` on the HTTP server.

## Metrics

### finality_gadget_finalized_blocks_total
- **Type**: Counter
- **Description**: Total number of finalized blocks processed by the finality gadget
- **Labels**: None
- **Usage**: Track overall finalization progress and rate

### finality_gadget_latest_finalized_block_height
- **Type**: Gauge  
- **Description**: Height of the latest block confirmed as finalized by QueryIsBlockBabylonFinalizedFromBabylon
- **Labels**: None
- **Usage**: Monitor current finalization status and detect stalls

### finality_gadget_fp_latest_block_voted
- **Type**: Gauge
- **Description**: Latest block height that each finality provider voted on
- **Labels**: 
  - `fp_pubkey`: Finality provider BTC public key (hex)
- **Usage**: Track individual FP participation and detect lagging providers

### finality_gadget_block_voters
- **Type**: Gauge
- **Description**: List of finality providers who voted for each block
- **Labels**:
  - `block_height`: Block height
  - `fp_pubkeys`: Comma-separated list of FP public keys who voted
- **Value**: Number of FPs who voted for this block
- **Usage**: Analyze voting patterns and block-by-block participation

### finality_gadget_fp_voting_power_per_block
- **Type**: Gauge
- **Description**: Voting power of each finality provider at each block height
- **Labels**:
  - `block_height`: Block height
  - `fp_pubkey`: Finality provider BTC public key (hex)
- **Value**: Voting power in satoshis (total BTC delegations)
- **Usage**: Track voting power distribution and changes over time

## Notes

- Voting power is measured in satoshis (1e8 satoshis = 1 BTC)
- FP public keys are BTC public keys in hexadecimal format
- Block heights correspond to L2 chain blocks
- Metrics are updated in real-time as blocks are processed and finalized 