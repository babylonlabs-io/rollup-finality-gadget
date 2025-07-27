# ⚠️ NOT FOR PRODUCTION USE ⚠️

> **This repository is a Proof-of-Concept (POC) and example implementation 
> to showcase how a Finality Gadget can work. It is not a production-ready 
> system and should not be used in production environments. We do not advise 
> companies or projects to use this codebase as-is.**

# Rollup BSN Finality Gadget

The **Rollup BSN Finality Gadget** is an off-chain program that can be 
run by the Rollup BSN network or by users of Rollup BSN. Its primary 
purpose is to track and determine whether a given L2 block is BTC-finalized. 
A block is considered BTC-finalized if it has received signatures from a quorum 
of registered Finality Providers (FPs) for that block, according to the rules and 
registration on the Babylon chain.

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
L2RPCHost = # RPC URL of OP stack L2 chain
BitcoinRPCHost = # Bitcoin RPC URL
DBFilePath = # Path to local bbolt DB file
FGContractAddress = # Babylon finality gadget contract address
BBNChainID = # Babylon chain id
BBNRPCAddress = # Babylon RPC host URL
GRPCListener = # Host:port to listen for gRPC connections
PollInterval = # Interval to poll for new L2 blocks
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

### Running tests

To run tests:

```bash
make test
```

## Querying Finality Status

You can query the finality status of blocks and transactions using the gRPC API provided by the Rollup BSN Finality Gadget. The main queries are:

1. **Is a block finalized?**
   - Use `QueryIsBlockBabylonFinalized` (checks by fetching on chain data)

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
