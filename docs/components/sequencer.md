# Component: Sequencer

## X1 Sequencer:

The X1 Sequencer is an optional but ancillary module that proposes new batches using transactions stored in the Pool Database.

## Running:

The preferred way to run the X1 Sequencer component is via Docker and Docker Compose.

```bash
docker pull okx/x1-node
```

To orchestrate multiple deployments of the different X1 Node components, a `docker-compose.yaml` file for Docker Compose can be used:

```yaml
  x1-sequencer:
    container_name: x1-sequencer
    image: x1-node
    command:
        - "/bin/sh"
        - "-c"
        - "/app/x1-node run --genesis /app/genesis.json --cfg /app/config.toml --components sequencer"
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
