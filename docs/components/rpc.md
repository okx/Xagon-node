# Component: RPC

## X1 RPC:

The X1 RPC relays transactions to the Trusted sequencer.

## Hard dependencies:

- [Synchronizer](./synchronizer.md)
- [StateDB Database](./databases.md)
- [RPCDB Database](./databases.md)
- [Merkle Tree and Executor](./prover.md)

## Running:

The preferred way to run the X1 RPC component is via Docker and Docker Compose.

```bash
docker pull okx/x1-node
```

To orchestrate multiple deployments of the different X1 Node components, a `docker-compose.yaml` file for Docker Compose can be used:

```yaml
  x1-rpc:
    container_name: x1-rpc
    image: x1-node
    command:
        - "/bin/sh"
        - "-c"
        - "/app/x1-node run --genesis /app/genesis.json --cfg /app/config.toml --components rpc"
```

The container alone needs some parameters configured, access to certain configuration files and the appropriate ports exposed.

- ports:
    - `8545:8545`: RPC Port
    - `9091:9091`: Needed if Prometheus metrics are enabled
- environment: Env variables that supersede the config file
    - `X1_NODE_STATE_DB_HOST`: Name of StateDB Database Host
    - `X1_NODE_POOL_HOST`: Name of PoolDB Database Host
    - `X1_NODE_RPC_DB_HOST`: Name of RPCDB Database Host
- volumes:
    - `your config.toml file`: /app/config.toml
    - `your genesis file`: /app/genesis.json
