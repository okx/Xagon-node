# zkEVM Node

zkEVM Node is a Go implementation of a node that operates the Polygon zkEVM Network.

## About the Polygon zkEVM network

Since this is an implementation of a protocol it's fundamental to understand it, [here](https://zkevm.polygon.technology/docs/zknode/zknode-overview) you can find the specification of the protocol.

Glossary:

- L1: Base blockchain where the rollup smart contracts are deployed. It's Ethereum or a testnet of Ethereum, but it could be any EVM compatible blockchain.
- L2: the rollup network aka the Polygon zkEVM network.
- Batch: a group of transactions that are executed/proved, using the [zkEVM prover](https://github.com/0xPolygonHermez/zkevm-prover) and sent to / synchronized from L1
- Sequencer: the actor that is responsible for selecting transactions, putting them in a specific order, and sending them in batches to L1
- Trusted sequencer: sequencer that has special privileges, there can only be one trusted sequencer. The privileges granted to the trusted sequencer allow it to forecast the batches that will be applied to L1. This way it can commit to a specific sequence before interacting with L1. This is done to achieve fast finality and reduce costs associated with using the network (lower gas fees)
- Permissionless sequencer: sequencer role that can be performed by anyone. It has competitive disadvantages compared to the trusted sequencer (slow finality, MEV attacks). Its main purpose is to provide censorship resistance and unstoppability features to the network.
- Sequence: Group of batches and other metadata that the trusted sequencer sends to L1 to update the state
- Forced batch: batch that is sent by permissionless sequencers to L1 to update the state
- L2 Block: Same as an L1 block, but for L2. This is mostly used by the JSON RPC interface. Currently, all the L2 Blocks are set to only include one transaction, this is done to achieve instant finality: it's not necessary to close a batch to allow the JSON RPC to expose results of already processed transactions
- Trusted state: state reached through processing transactions that have been shared by the trusted sequencer. This state is considered trusted as the trusted sequencer could commit to a certain sequence, and then send a different one to L1
- Virtual state: state reached through processing transactions that have already been submitted to L1. These transactions are sent in batches by either trusted or permissionless sequencers. Those batches are also called virtual batches. Note that this state is trustless as it relies on L1 security assumptions
- Consolidated state: state that is proven on-chain by submitting a ZKP (Zero Knowledge Proof) that proves the execution of a sequence of the last virtual batch.
- Invalid transaction: a transaction that can't be processed and doesn't affect the state. Note that such a transaction could be included in a virtual batch. The reason for a transaction to be invalid could be related to the Ethereum protocol (invalid nonce, not enough balance, ...) or due to limitations introduced by the zkEVM (each batch can make use of a limited amount of resources such as the total amount of keccak hashes that can be computed)
- Reverted transaction: a transaction that is executed, but is reverted (because of smart contract logic). The main difference with *invalid transaction* is that this transaction modifies the state, at least to increment nonce of the sender.

## Architecture

<p align="center">
  <img src="./docs/architecture.drawio.png"/>
</p>

The diagram represents the main components of the software and how they interact between them. Note that this reflects a single entity running a node, in particular a node that acts as the trusted sequencer. But there are many entities running nodes in the network, and each of these entities can perform different roles. More on this later.

- (JSON) RPC: an HTTP interface that allows users (dApps, metamask, etherscan, ...) to interact with the node. Fully compatible with Ethereum RPC + some extra [custom endpoints](./docs/zkEVM-custom-endpoints.md) specifics of the network. It interacts with the `state` (to get data and process transactions) as well as the `pool` (to store transactions).
- L2GasPricer: it fetches the L1 gas price and applies some formula to calculate the gas price that will be suggested for the users to use for paying fees on L2. The suggestions are stored on the `pool`, and will be consumed by the `rpc`
- Pool: DB that stores transactions by the `RPC` to be selected/discarded by the `sequencer` later on
- Sequencer: responsible for building the trusted state. To do so, it gets transactions from the pool and puts them in a specific order. It needs to take care of opening and closing batches while trying to make them as full as possible. To achieve this it needs to use the executor to actually process the transaction not only to execute the state transition (and update the hashDB) but also to check the consumed resources by the transactions and the remaining resources of the batch. After executing a transaction that fits into a batch, it gets stored on the `state`. Once transactions are added into the state, they are immediately available through the `rpc`.
- SequenceSender: gets closed batches from the `state`, tries to aggregate as many of them as possible, and at some point, decides that it's time to send those batches to L1, turning the state from trusted to virtualized. In order to send the L1 tx, it uses the `ethtxmanager`
- EthTxManager: handles requests to send L1 transactions from `sequencesender` and `aggregator`. It takes care of dealing with the nonce of the accounts, increasing the gas price, and other actions that may be needed to ensure that L1 transactions get mined
- Etherman: abstraction that implements the needed methods to interact with the Ethereum network and the relevant smart contracts.
- Synchronizer: Updates the `state` (virtual batches, verified batches, forced batches, ...) by fetching data from L1 through the `etherman`. If the node is not a `trusted sequencer` it also updates the state with the data fetched from the `rpc` of the `trusted sequencer`. It also detects and handles reorgs that can happen if the `trusted sequencer` sends different data in the rpc vs the sequences sent to L1 (trusted reorg aka L2 reorg). Also handles L1 reorgs (reorgs that happen on the L1 network)
- State: Responsible for managing the state data (batches, blocks, transactions, ...) that is stored on the `state SB`. It also handles the integration with the `executor` and the `Merkletree` service
- State DB: persistence layer for the state data (except the Merkletree that is handled by the `HashDB` service), it stores informationrelated to L1 (blocks, global exit root updates, ...) and L2 (batches, L2 blocks, transactions, ...)
- Aggregator: consolidates batches by generating ZKPs (Zero Knowledge proofs). To do so it gathers the necessary data that the `prover` needs as input through the `state` and sends a request to it. Once the proof is generated it sends a request to send an L1 tx to verify the proof and move the state from virtual to verified to the `ethtxmanager`. Note that provers connect to the aggregator and not the other way around. The aggregator can handle multiple connected provers at once and make them work concurrently in the generation of different proofs
- Prover/Executor/hashDB: service that generates ZK proofs. Note that this component is not implemented in this repository, and it's treated as a "black box" from the perspective of the node. The prover/executor has two implementations: [JS reference implementation](https://github.com/0xPolygonHermez/zkevm-proverjs) and [C production-ready implementation](https://github.com/0xPolygonHermez/zkevm-prover). Although it's the same software/binary, it implements three services:
  - Executor: Provides an EVM implementation that allows processing batches as well as getting metadata (state root, transaction receipts, logs, ...) of all the needed results.
  - Prover: Generates ZKPs for batches, batches aggregation, and final proofs.
  - HashDB: service that stores the Merkletree, containing all the account information (balances, nonces, smart contract code, and smart contract storage)

## Roles of the network

The node software is designed to support the execution of multiple roles. Each role requires different services to work. Most of the services can run in different instances, and the JSON RPC can run in many instances (all the other services must have a single instance)

### RPC

This role can be performed by anyone.

Required services and components:

- JSON RPC: can run in a separated instance, and can have multiple instances
- Synchronizer: single instance that can run on a separate instance
- Executor & Merkletree: service that can run on a separate instance
- State DB: Postgres SQL that can be run in a separate instance

There must be only one synchronizer, and it's recommended that it has exclusive access to an executor instance, although it's not necessary. This role can perfectly be run in a single instance, however, the JSON RPC and executor services can benefit from running in multiple instances, if the performance decreases due to the number of requests received

- [`zkEVM RPC endpoints`](./docs/json-rpc-endpoints.md)
- [`zkEVM RPC Custom endpoints documentation`](./docs/zkEVM-custom-endpoints.md)

### Trusted sequencer

This role can only be performed by a single entity. This is enforced in the smart contract, as the related methods of the trusted sequencer can only be performed by the owner of a particular private key.

Required services and components:

- JSON RPC: can run in a separated instance, and can have multiple instances
- Sequencer & Synchronizer: single instance that needs to run together
- Executor & Merkletree: service that can run on a separate instance
- Pool DB: Postgres SQL that can be run in a separate instance
- State DB: Postgres SQL that can be run in a separate instance

Note that the JSON RPC is required to receive transactions. It's recommended that the JSON RPC runs on separated instances, and potentially more than one (depending on the load of the network). It's also recommended that the JSON RPC and the Sequencer don't share the same executor instance, to make sure that the sequencer has exclusive access to an executor

### Aggregator

This role can be performed by anyone.

Required services and components:

- Synchronizer: single instance that can run on a separated instance
- Executor & Merkletree: service that can run on a separate instance
- State DB: Postgres SQL that can be run in a separate instance
- Aggregator: single instance that can run on a separated instance
- Prover: single instance that can run on a separated instance
- Executor: single instance that can run on a separated instance

It's recommended that the prover is run on a separate instance, as it has important hardware requirements. On the other hand, all the other components can run on a single instance,

## Development

It's recommended to use `make` for building, and testing the code, ... Run `make help` to get a list of the available commands.

## Running the node

- [Running locally](docs/running_local.md)
- [Running on production](docs/production-setup.md)

### Requirements

- Go 1.21
- Docker
- Docker Compose
- Make
- GCC

## Contribute

Before opening a pull request, please read this [guide](CONTRIBUTING.md).


