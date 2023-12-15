## Configuration

To configure a node you need 3 files:
- Node configuration
- Genesis configuration
- Prover configuration

### Node Config
This file is a [TOML](https://en.wikipedia.org/wiki/TOML#) formatted file.
You could find some examples here:
 - `config/environments/local/local.node.config.toml`: running a permisionless node
  - `config/environments/mainnet/node.config.toml`
  - `config/environments/public/node.config.toml`
  - `test/config/test.node.config.toml`: configuration for a trusted node used in CI

  For details about the contents you can read specifications [here](config-file/node-config-doc.md)

This file is used for trusted and for permisionless nodes. In the case of permissionless node you only need to setup next sections:

 You could **override values with environment variables**. The variables needs to have next format:
`ZKEVM_NODE_`[<section in uppercase\>`_`]*<var in uppercase\>`

For example:
`ZKEVM_NODE_STATE_DB_HOST="localhost"` override value of section `[StateDB]` key `Host`

### Network Genesis Config
This file is a [JSON](https://en.wikipedia.org/wiki/JSON) formatted file.
This contain all the info information relating to the relation between L1 and L2 network's (e.g. contracts, etc..) also known as genesis file

You could find an example here:
- `config/environments/local/local.genesis.config.json`:

For details about the contents you can read specifications [here](config-file/custom_network-config-doc.md)


### Prover Config

Please check [prover repository](https://github.com/okx/x1-prover)  for further information

Examples:
 - `config/environments/mainnet/prover.config.json`
 - `config/environments/testnet/prover.config.json`
