# Component: Prover

NOTE: The Prover is not considered part of the X1 Node and all issues and suggestions should be sent to the [Prover repo](https://github.com/okx/x1-prover/).

## X1 Prover:

The X1 Prover image hosts different components, *Merkle Tree*, *Executor* and finally the actual *Prover*.

## Hard dependencies:

- [Aggregator](./aggregator.md)

## Running:

The preferred way to run the X1 Prover component is via Docker and Docker Compose.

```bash
docker pull hermeznetwork/x1-prover
```

To orchestrate multiple deployments of the different X1 Node components, a `docker-compose.yaml` file for Docker Compose can be used:

```yaml
  x1-prover:
    container_name: x1-prover
    image: x1-prover
    volumes:
      - ./prover-config.json:/usr/src/app/config.json
    command: >
      zkProver -c /usr/src/app/config.json
```

The `prover-config.json` file contents will vary depending on your use case, the main document explains different values to be changed to achieve different behaviors.

### Ports:

- `50051`: Prover
- `50061`: Merkle Tree
- `50071`: Executor