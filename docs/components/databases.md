# Component: Databases

Different services rely on several standard PostgresQL databases, the most important one being the State Database.

Configuring each DB is trivial if done with an orchestration tool such as Docker Compose.

The following are examples on how one would provision each DB to serve along the components (rpc, aggregator, sequencer...)

Note the `environment` values will change per DB.

- **StateDB**:

The StateDB needs to generate some extra databases and tables (`merkletree`) for use with the MerkleTree/Executor service.

This is done via an sql file: [init_prover_db.sql](https://github.com/okx/x1-node/blob/develop/db/scripts/init_prover_db.sql)

```yaml
x1-state-db:
    container_name: x1-state-db
    image: postgres:15
    deploy:
      resources:
        limits:
          memory: 2G
        reservations:
          memory: 1G
    ports:
      - 5432:5432
    volumes:
      - ./init_prover_db.sql:/docker-entrypoint-initdb.d/init.sql
    environment:
      - POSTGRES_USER=state_user
      - POSTGRES_PASSWORD=state_password
      - POSTGRES_DB=state_db
    command: ["postgres", "-N", "500"]
```

- **Other DBs: Pool/RPC**:

```yaml
  x1-pool-db:
    container_name: x1-pool-db
    image: postgres:15
    deploy:
      resources:
        limits:
          memory: 2G
        reservations:
          memory: 1G
    ports:
      - 5433:5432
    environment:
      - POSTGRES_USER=pool_user
      - POSTGRES_PASSWORD=pool_password
      - POSTGRES_DB=pool_db
    command: ["postgres", "-N", "500"]
```
