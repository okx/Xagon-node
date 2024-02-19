# Production Setup

This guide describes how to run a node that:

- Synchronizes the network
- Expose a JSON RPC interface, acting as an archive node

Note that sequencing and proving functionalities are not covered in this document **yet**.

## Requirements

- A machine to run the zkEVM node with the following requirements:
  - Hardware: 32G RAM, 4 cores, 128G Disk with high IOPS (as the network is super young the current disk requirements are quite low, but they will increase over time. Also note that this requirement is true if the DBs run on the same machine, but it's recommended to run Postgres on dedicated infra). Currently ARM-based CPUs are not supported
  - Software: Ubuntu 22.04, Docker
- A L1 node: we recommend using geth, but what it's actually needed is access to a JSON RPC interface for the L1 network (Sepolia for Cardona zkEVM testnet, Goerli for zkEVM testnet, Ethereum mainnet for zkEVM mainnet)

## Setup

This is the most straightforward path to run a zkEVM node, and it's perfectly fine for most use cases, however if you are interested in providing service to many users it's recommended to do some tweaking over the default configuration. Furthermore, this is quite opinionated, feel free to run this software in a different way, for instance it's not needed to use Docker, you could use the Go and C++ binaries directly.

tl;dr:

```bash
# DOWNLOAD ARTIFACTS
ZKEVM_NET=mainnet
ZKEVM_DIR=./path/to/install # CHANGE THIS
ZKEVM_CONFIG_DIR=./path/to/config  # CHANGE THIS
curl -L https://github.com/0xPolygonHermez/zkevm-node/releases/latest/download/$ZKEVM_NET.zip > $ZKEVM_NET.zip && unzip -o $ZKEVM_NET.zip -d $ZKEVM_DIR && rm $ZKEVM_NET.zip
mkdir -p $ZKEVM_CONFIG_DIR && cp $ZKEVM_DIR/$ZKEVM_NET/example.env $ZKEVM_CONFIG_DIR/.env

# EDIT THIS env file:
nano $ZKEVM_CONFIG_DIR/.env

# RUN:
docker compose --env-file $ZKEVM_CONFIG_DIR/.env -f $ZKEVM_DIR/$ZKEVM_NET/docker-compose.yml up -d
```

### Explained step by step:

1. Define network: `ZKEVM_NET=cardona` or `ZKEVM_NET=testnet` or `ZKEVM_NET=mainnet`
2. Define installation path: `ZKEVM_DIR=./path/to/install`
3. Define a config directory: `ZKEVM_CONFIG_DIR=./path/to/config`
4. It's recommended to source this env vars in your `~/.bashrc`, `~/.zshrc` or whatever you're using
5. Download and extract the artifacts: `curl -L https://github.com/0xPolygonHermez/zkevm-node/releases/latest/download/$ZKEVM_NET.zip > $ZKEVM_NET.zip && unzip -o $ZKEVM_NET.zip -d $ZKEVM_DIR && rm $ZKEVM_NET.zip`. Note you may need to install `unzip` for this command to work. 

> **NOTE:** Take into account this works for the latest release (mainnet), in case you want to deploy a pre-release (testnet) you should get the artifacts directly for that release and not using the "latest" link depicted here. [Here](https://github.com/0xPolygonHermez) you can check the node release deployed for each network.

6. Copy the file with the env parameters into config directory: `mkdir -p $ZKEVM_CONFIG_DIR && cp $ZKEVM_DIR/$ZKEVM_NET/example.env $ZKEVM_CONFIG_DIR/.env`
7. Edit the env file, with your favourite editor. The example will use nano: `nano $ZKEVM_CONFIG_DIR/.env`. This file contains the configuration that anyone should modify. For advanced configuration:
   1. Copy the config files into the config directory `cp $ZKEVM_DIR/$ZKEVM_NET/config/environments/$ZKEVM_NET/* $ZKEVM_CONFIG_DIR/`
   2. Make sure the modify the `ZKEVM_ADVANCED_CONFIG_DIR` from `$ZKEVM_CONFIG_DIR/.env` with the correct path
   3. Edit the different configuration files in the $ZKEVM_CONFIG_DIR directory and make the necessary changes
8. Run the node: `docker compose --env-file $ZKEVM_CONFIG_DIR/.env -f $ZKEVM_DIR/$ZKEVM_NET/docker-compose.yml up -d`. You may need to run this command using `sudo` depending on your Docker setup.
9. Make sure that all components are running: `docker compose --env-file $ZKEVM_CONFIG_DIR/.env -f $ZKEVM_DIR/$ZKEVM_NET/docker-compose.yml ps`. You should see the following containers:
   1. zkevm-rpc
   2. zkevm-sync
   3. zkevm-state-db
   4. zkevm-pool-db
   5. zkevm-prover
10. If everything has gone as expected you should be able to run queries to the JSON RPC at `http://localhost:8545`. For instance you can run the following query that fetches the latest synchronized L2 block, if you call this every few seconds, you should see the number increasing: 

`curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}' http://localhost:8545` 

## Troubleshooting

- It's possible that the machine you're using already uses some of the necessary ports. In this case you can change them directly at `$ZKEVM_DIR/$ZKEVM_NET/docker-compose.yml`
- If one or more containers are crashing please check the logs using `docker compose --env-file $ZKEVM_CONFIG_DIR/.env -f $ZKEVM_DIR/$ZKEVM_NET/docker-compose.yml logs <cointainer_name>`

## Stop

```bash
docker compose --env-file $ZKEVM_CONFIG_DIR/.env -f $ZKEVM_DIR/$ZKEVM_NET/docker-compose.yml down
```

## Updating

In order to update the software, you have to repeat the steps of the setup, but taking care of not overriding the config that you have modified. Basically, instead of running `cp $ZKEVM_DIR/$ZKEVM_NET/example.env $ZKEVM_CONFIG_DIR/.env`, check if the variables of `$ZKEVM_DIR/$ZKEVM_NET/example.env` have been renamed or there are new ones, and update `$ZKEVM_CONFIG_DIR/.env` accordingly.

## Advanced setup

> DISCLAIMER: right now this part of the documentation attempts to give ideas on how to improve the setup for better performance, but is far from being a detailed guide on how to achieve this. Please open issues requesting more details if you don't understand how to achieve something. We will keep improving this doc for sure!

There are some fundamental changes that can be done towards the basic setup, in order to get better performance and scale better:

### DB

In the basic setup, there are Postgres being instantiated as Docker containers. For better performance is recommended to:

- Run dedicated instances for Postgres. To achieve this you will need to:
  - Remove the Postgres services (`zkevm-pool-db` and `zkevm-state-db`) from the `docker-compose.yml`
  - Instantiate Postgres elsewhere (note that you will have to create credentials and run some queries to make this work, following the config files and docker-compose should give a clear idea of what to do)
  - Update the `node.config.toml` to use the correct URI for both DBs
  - Update `prover.config.json` to use the correct URI for the state DB
- Use a setup of Postgres that allows to have separated endpoints for read / write replicas

### JSON RPC

Unlike the synchronizer, that needs to have only one instance running (having more than one synchronizer running at the same time connected to the same DB can be fatal), the JSON RPC can scale horizontally.

There can be as many instances of it as needed, but in order to not introduce other bottlenecks, it's important to consider the following:

- Read replicas of the State DB should be used
- Synchronizer should have an exclusive instance of `zkevm-prover`
- JSON RPCs should scale in correlation with instances of `zkevm-prover`. The most obvious way to do so is by having a dedicated `zkevm-prover` for each `zkevm-rpc`. But depending on the payload of your solution it could be worth to have `1 zkevm-rpc : many zkevm-prover` or `many zkevm-rpc : 1 zkevm-prover`, ... For reference, the `zkevm-prover` implements the EVM, and therefore will be heavily used when calling endpoints such as `eth_call`. On the other hand, there are other endpoints that relay on the `zkevm-state-db`
