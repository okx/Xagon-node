# Component: Synchronizer

## X1 Synchronizer:

The X1 Synchronizer is the **base** component for which all others will depend on. You can *mix and match* different components to achieve a different outcome, be it sending transactions or computing proofs, but the Sync module will need to be up and running.

This module syncs data between the Layer 1 Ethereum network and X1 L2 network.

## Running:

The preferred way to run the X1 Synchronizer component is via Docker and Docker Compose.

```bash
docker pull hermeznetwork/x1-node
```

To orchestrate multiple deployments of the different X1 Node components, a `docker-compose.yaml` file for Docker Compose can be used:

**THIS STEP IS MANDATORY ON ALL DEPLOYMENT MODES**

```yaml
  x1-sync:
    container_name: x1-sync
    image: x1-node
    command:
        - "/bin/sh"
        - "-c"
        - "/app/x1-node run --genesis /app/genesis.json --cfg /app/config.toml --components synchronizer"
```

The container alone needs some parameters configured, access to certain configuration files and the appropriate ports exposed.

- environment: Env variables that supersede the config file
    - `X1_NODE_STATE_DB_HOST`: Name of StateDB Database Host
- volumes:
    - `your config.toml file`: /app/config.toml
    - `your genesis.json file`: /app/genesis.json
