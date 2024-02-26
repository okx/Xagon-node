# Component: Sequencer

## XLayer Sequencer:

The XLayer Sequencer is an optional but ancillary module that proposes new batches using transactions stored in the Pool Database.

## Running:

The preferred way to run the XLayer Sequencer component is via Docker and Docker Compose.

```bash
docker pull okx/xlayer-node
```

To orchestrate multiple deployments of the different XLayer Node components, a `docker-compose.yaml` file for Docker Compose can be used:

```yaml
  xlayer-sequencer:
    container_name: xlayer-sequencer
    image: xlayer-node
    command:
        - "/bin/sh"
        - "-c"
        - "/app/xlayer-node run --genesis /app/genesis.json --cfg /app/config.toml --components sequencer"
```

The container alone needs some parameters configured, access to certain configuration files and the appropriate ports exposed.

- environment: Env variables that supersede the config file
    - `X1_NODE_POOLDB_HOST`: Name of PoolDB Database Host
    - `X1_NODE_STATE_DB_HOST`: Name of StateDB Database Host
- volumes:
    - `your Account Keystore file`: /pk/keystore (note, this `/pk/keystore` value is the default path that's written in the Public Configuration files on this repo, meant to expedite deployments, it can be superseded via an env flag `X1_NODE_ETHERMAN_PRIVATEKEYPATH`.)
    - `your config.toml file`: /app/config.toml
    - `your genesis.json file`: /app/genesis.json

[How to generate an account keystore](./account_keystore.md)
