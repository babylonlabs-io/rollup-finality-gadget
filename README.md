# ⚠️ NOT FOR PRODUCTION USE ⚠️

> **This repository is a Proof-of-Concept (POC) and example implementation 
> to showcase how a Finality Gadget can work. It is not a production-ready 
> system and should not be used in production environments. We do not advise 
> companies or projects to use this codebase as-is.**

# Rollup BSN Finality Gadget

The **Rollup BSN Finality Gadget** is an off-chain program that can be 
run by the Rollup BSN network or by users of Rollup BSN. Its primary 
purpose is to track and determine whether a given L2 block is BTC-finalized. 
A block is considered BTC‑finalized if it has:
- received signatures from a quorum of registered Finality Providers (FPs) for that block, 
according to the rules and registration on the Babylon Genesis chain
- its prior block at height h − X is also BTC‑staking‑finalized, where X is the finality signature interval

This tool enables monitoring and querying of the BTC-finalized status of L2 blocks, 
providing an extra layer of security and finality for rollup chains integrated with 
the Rollup BSN ecosystem.


## Modules

- `cmd` : entry point for `opfgd` finality gadget daemon
- `finalitygadget` : top-level umbrella module that exposes query methods and coordinates calls to other clients
- `client` : grpc client to query the finality gadget
- `server` : grpc server for the finality gadget
- `proto` : protobuf definitions for the grpc server
- `config` : configs for the finality gadget
- `btcclient` : wrapper around Bitcoin RPC client
- `bbnclient` : wrapper around Babylon RPC client
- `ethl2client` : wrapper around OP stack L2 ETH RPC client
- `cwclient` : client to query CosmWasm smart contract deployed on BabylonChain
- `db` : handler for local database to store finalized block state
- `types` : common types
- `log` : custom logger
- `testutil` : test utilities and helpers

## Instructions

### Download and configuration

To get started, clone the repository.

```bash
git clone https://github.com/babylonlabs-io/finality-gadget.git
```

Copy the `config.toml.example` file to `config.toml`:

```bash
cp config.toml.example config.toml
```

Configure the `config.toml` file with the following parameters:

```toml
# L2 Chain Configuration
L2RPCHost = "http://localhost:8545"        # RPC URL of OP stack L2 chain

# Bitcoin Node Configuration  
BitcoinRPCHost = "localhost:18443"         # Bitcoin RPC URL
BitcoinRPCUser = "rpcuser"                 # Bitcoin RPC username (optional)
BitcoinRPCPass = "rpcpass"                 # Bitcoin RPC password (optional)  
BitcoinDisableTLS = true                   # Disable TLS for HTTP connections (required for http://)

# Babylon Chain Configuration
FGContractAddress = ""                     # Rollup BSN contract address
BBNChainID = ""                            # Babylon chain ID
BBNRPCAddress = "http://localhost:26657"   # Babylon RPC host URL

# Database Configuration
DBFilePath = "./finalitygadget.db"         # Path to local bbolt DB file

# Server Configuration
GRPCListener = "0.0.0.0:50051"             # Host:port to listen for gRPC connections
HTTPListener = "0.0.0.0:8085"              # Host:port to listen for HTTP connections

# Processing Configuration
PollInterval = "1s"                        # Interval to poll for new L2 blocks
BatchSize = 1                              # Number of blocks to process in a batch
StartBlockHeight = 0                       # Block height to start processing from (0 = use latest)
LogLevel = "info"                          # Log level (debug, info, warn, error)
```

### Building and installing the binary

At the top-level directory of the project

```bash
make install
```

The above command will build and install the `opfgd` binary to
`$GOPATH/bin`.

If your shell cannot find the installed binaries, make sure `$GOPATH/bin` is in
the `$PATH` of your shell. Usually these commands will do the job

```bash
export PATH=$HOME/go/bin:$PATH
echo 'export PATH=$HOME/go/bin:$PATH' >> ~/.profile
```

### Running the daemon

To start the daemon, run:

```bash
opfgd start --cfg config.toml
```

## Querying Block Finality Status

The finality gadget provides gRPC endpoints to query the finalization status of blocks. You can use `grpcurl` to test these endpoints.

### Prerequisites

Install `grpcurl` if you haven't already:

```bash
# macOS
brew install grpcurl

# Ubuntu/Debian
apt install grpcurl

# Or go install
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### Available gRPC Endpoints

#### 1. Check if a specific block height is finalized

```bash
grpcurl -plaintext -proto proto/finalitygadget.proto \
  -d '{"block_height": 27259}' \
  localhost:50051 proto.FinalityGadget/QueryIsBlockFinalizedByHeight
```

**Response:**
```json
{
  "isFinalized": true
}
```

#### 2. Get the latest finalized block

```bash
grpcurl -plaintext -proto proto/finalitygadget.proto \
  -d '{}' \
  localhost:50051 proto.FinalityGadget/QueryLatestFinalizedBlock
```

**Response:**
```json
{
  "block": {
    "blockHash": "0x511b9c896c93ab4b82ce3be5061ec3323db7dfdaa02d3e22579c14658e0f1ca3",
    "blockHeight": "26788",
    "blockTimestamp": "1753706277"
  }
}
```

## Build Docker image

### Prerequisites

1. **Docker Desktop**: Install from [Docker's official website](https://docs.docker.com/desktop/).

2. **Make**: Required for building service binaries. Installation guide available [here](https://sp21.datastructur.es/materials/guides/make-install.html).

3. **GitHub SSH Key**:
   - Create a non-passphrase-protected SSH key.
   - Add it to GitHub ([instructions](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account)).
   - Export the key path:
     ```shell
     export BBN_PRIV_DEPLOY_KEY=FULL_PATH_TO_PRIVATE_KEY/.ssh/id_ed25519
     ```

4. **Repository Setup**:
   ```shell
   git clone https://github.com/babylonlabs-io/finality-gadget.git
   ```
To build the docker image:

```bash
make build-docker
```
