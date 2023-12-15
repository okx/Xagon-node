# How to use snapshot/restore feature

This feature is for fast replication of nodes. It creates a backup of database and allows restoration in another database to save synchronization time.
- It uses the tools `pg_dump` and `pg_restore` and requires them to match the same version as the server.

## Snapshot

This feature creates a dump of entire database

### Usage

```
NAME:
   x1-node snapshot - Snapshot the state db

USAGE:
   x1-node snapshot [command options] [arguments...]

OPTIONS:
   --cfg FILE, -c FILE  Configuration FILE
   --help, -h           show help
```

**Make sure that the config file contains the data required to connect to `HashDB` database**, for example:
```
[HashDB]
User = "prover_user"
Password = "prover_pass"
Name = "prover_db"
Host = "x1-state-db"
Port = "5432"
EnableLog = false
MaxConns = 200
```

This generates two files in the current working path:
* For stateDB: <database_name>`_`\<timestamp>`_`\<version>`_`\<gitrev>`.sql.tar.gz`
* For hashDB: <database_name>`_`\<timestamp>`_`\<version>`_`\<gitrev>`.sql.tar.gz`

#### Example of invocation:
```
# cd /tmp/ && /app/x1-node snap -c /app/config.toml
(...)
# ls -1
prover_db_1689925019_v0.2.0-RC9-15-gd39e7f1e_d39e7f1e.sql.tar.gz
state_db_1689925019_v0.2.0-RC9-15-gd39e7f1e_d39e7f1e.sql.tar.gz
```


## Restore
It populates state, and hash databases with the previous backup

**Be sure that none node service is running!**

### Usage

```
NAME:
   x1-node restore - Restore snapshot of the state db

USAGE:
   x1-node restore [command options] [arguments...]

OPTIONS:
   --inputfilestate value, --is value  Input file stateDB
   --inputfileHash value, --ih value   Input file hashDB
   --cfg FILE, -c FILE                 Configuration FILE
   --help, -h                          show help
```

#### Example of invocation:
```
/app/x1-node restore -c /app/config.toml  --is /tmp/state_db_1689925019_v0.2.0-RC9-15-gd39e7f1e_d39e7f1e.sql.tar.gz  --ih /tmp/prover_db_1689925019_v0.2.0-RC9-15-gd39e7f1e_d39e7f1e.sql.tar
.gz
```

# How to test
You could use `test/docker-compose.yml` to interact with `x1-node`:
* Run the containers: `make run`
* Launch a interactive container:
```
docker-compose up -d x1-sh
docker-compose exec x1-sh /bin/sh
```
* Inside this shell you can execute the examples of invocation
