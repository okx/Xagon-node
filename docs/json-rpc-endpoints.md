# JSON RPC Endpoints

Here you will find the list of all supported JSON RPC endpoints and any differences between them in comparison to the default behavior of an ethereum node.

If the endpoint is not in the list below, it means this specific endpoint is not supported yet, feel free to open an issue requesting it to be added and please explain the reason why you need it. 

> Warning: debug endpoints are considered experimental as they have not been deeply tested yet
<!-- DEBUG -->
- `debug_traceBlockByHash`
- `debug_traceBlockByNumber`
- `debug_traceTransaction`
- `debug_traceBatchByNumber`

<!-- ETH -->
- `eth_blockNumber`
- `eth_call`
  - _doesn't support state override at the moment and pending block. Will be implemented [#1990](https://github.com/0xPolygonHermez/zkevm-node/issues/1990)_ 
  - _doesn't support `from` values that are smart contract addresses. Will be implemented [#2017](https://github.com/0xPolygonHermez/zkevm-node/issues/2017)_  
- `eth_chainId`
- `eth_estimateGas` _* if the block number is set to pending we assume it is the latest_
- `eth_gasPrice`
- `eth_getBalance` _* if the block number is set to pending we assume it is the latest_
- `eth_getBlockByHash` _* allows an extra boolean parameter to query l2 extra information_
- `eth_getBlockByNumber` _* allows an extra boolean parameter to query l2 extra information_
- `eth_getBlockTransactionCountByHash`
- `eth_getBlockTransactionCountByNumber`
- `eth_getCode` _* if the block number is set to pending we assume it is the latest_
- `eth_getCompilers` _* response is always empty_
- `eth_getFilterChanges`
- `eth_getFilterLogs`
- `eth_getLogs`
- `eth_getStorageAt` _* if the block number is set to pending we assume it is the latest_
- `eth_getTransactionByBlockHashAndIndex` _* allows an extra boolean parameter to query l2 extra information_
- `eth_getTransactionByBlockNumberAndIndex` _* if the block number is set to pending we assume it is the latest; * allows an extra boolean parameter to query l2 extra information_
- `eth_getTransactionByHash` _* allows an extra boolean parameter to query l2 extra information_
- `eth_getTransactionCount`
- `eth_getTransactionReceipt` _* doesn't include effectiveGasPrice. Will include once EIP1559 is implemented_
- `eth_getUncleByBlockHashAndIndex` _* response is always empty_
- `eth_getUncleByBlockNumberAndIndex` _* response is always empty_
- `eth_getUncleCountByBlockHash` _* response is always zero_
- `eth_getUncleCountByBlockNumber` _* response is always zero_
- `eth_newBlockFilter`
- `eth_newFilter`
- `eth_protocolVersion` _* response is always zero_
- `eth_sendRawTransaction` _* can relay TXs to another node_
- `eth_subscribe`
- `eth_syncing`
- `eth_uninstallFilter`
- `eth_unsubscribe`

<!-- NET -->
- `net_version`

<!-- TXPOOL -->
- `txpool_content` _* response is always empty_

<!-- WEB3 -->
- `web3_clientVersion`
- `web3_sha3`

<!-- ZKEVM -->
- `zkevm_batchNumber`
- `zkevm_batchNumberByBlockNumber`
- `zkevm_consolidatedBlockNumber`
- `zkevm_estimateFee`
- `zkevm_estimateGasPrice`
- `zkevm_estimateCounters`
- `zkevm_getBatchByNumber`
- `zkevm_getExitRootsByGER`
- `zkevm_getFullBlockByHash`
- `zkevm_getFullBlockByNumber`
- `zkevm_getLatestGlobalExitRoot`
- `zkevm_getNativeBlockHashesInRange`
- `zkevm_getTransactionByL2Hash`
- `zkevm_getTransactionReceiptByL2Hash`
- `zkevm_isBlockConsolidated`
- `zkevm_isBlockVirtualized`
- `zkevm_verifiedBatchNumber`
- `zkevm_virtualBatchNumber`
