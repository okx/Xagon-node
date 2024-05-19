# upgrade testnet rpc fork.9
XLayer Testnet has been upgraded to Fork9. You can upgrade permission less rpc nodes in two ways. 

For more detailed information, please refer to:https://www.okx.com/zh-hans/xlayer/docs/developer/setup-zknode/setup-production-zknode
## 1.Re-download the script and snapshot. It will take about 3 hours. (Recommended)
``` bash
apt install axel  # Parallel download tools
apt install pigz  # Parallel compression tools
wget https://static.okex.org/cdn/chain/xlayer/snapshot/run_xlayer_testnet.sh && chmod +x run_xlayer_testnet.sh && ./run_xlayer_testnet.sh init && cp ./testnet/example.env ./testnet/.env
vim ./testnet/.env # Modify XLAYER_NODE_ETHERMAN_URL = "http://your.L1node.url"
./run_xlayer_testnet.sh restore 
./run_xlayer_testnet.sh start
```

## 2.Manual upgrade
### 2.1 Stop the xlayer-sync and xlayer-rpc services.
### 2.2 Modify the configuration file.
Modify the docker-compose.yml file and replace the image version as follows:
```
xlayer-rpc->image: okexchain/xlayer-node:origin_sync-v0.3.6_20240412160714_76857189
xlayer-sync->image: okexchain/xlayer-node:origin_sync-v0.3.6_20240412160714_76857189
xlayer-prover->image: okexchain/xlayer-prover:origin_release_v0.3.1_20240327040854_458b8d26
```

Modify the config/genesis.config.json file and add the fields rollupCreationBlockNumber and rollupManagerCreationBlockNumber as follows:
```
{
  ...
  "genesisBlockNumber":  4648290,

  "rollupCreationBlockNumber": 4648290,
  "rollupManagerCreationBlockNumber": 4648290,

  "root": "0xb2fbff62137228e52809081a425bfcd30c0fdc8c1213085278c739676a7669b8",
  ...
}
```

Modify the vim config/node.config.toml file and add the Fork9UpgradeBatch field as follows:
```
...
ForkUpgradeBatchNumber=0
ForkUpgradeNewForkId=0

Fork9UpgradeBatch=476000
...

```

### 2.3 Start the xlayer-sync and xlayer-rpc services.

