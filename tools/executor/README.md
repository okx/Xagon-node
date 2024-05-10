# Executor tool

Tool intended to test and debug the executor. The main point is to provide JSON files as test vectors that resembles an executor request.
This way we can easily reproduce scenarios that are problematic.

## Run

`cd` into this directory and `go run .`. You can decide which vectors will be run by seting `skip` to true/false, on the json files found at `./vectors`.

## Generate vector

Export the logs of the sequencer into a separated file and look for logs that look like this:

```log
2022-08-25T13:16:08Z	[35mDEBUG[0m	state/state.go:344	*******************************************
2022-08-25T13:16:08Z	[35mDEBUG[0m	state/state.go:345	ProcessSequencerBatch start
2022-08-25T13:16:08Z	[35mDEBUG[0m	state/helper.go:26	0 1000000000 519947 <nil> 0 2661 1000
2022-08-25T13:16:08Z	[35mDEBUG[0m	state/state.go:405	processBatch[processBatchRequest.BatchNum]: 1
2022-08-25T13:16:08Z	[35mDEBUG[0m	state/state.go:406	processBatch[processBatchRequest.BatchL2Data]: 0xf90a7980843b9aca0...ecadc11c
2022-08-25T13:16:08Z	[35mDEBUG[0m	state/state.go:407	processBatch[processBatchRequest.From]: 
2022-08-25T13:16:08Z	[35mDEBUG[0m	state/state.go:408	processBatch[processBatchRequest.OldStateRoot]: 0xc8f751215e3f69b83361af8cdf8c657773743af845e313659f6a8189b098e038
2022-08-25T13:16:08Z	[35mDEBUG[0m	state/state.go:409	processBatch[processBatchRequest.GlobalExitRoot]: 0x0000000000000000000000000000000000000000000000000000000000000000
2022-08-25T13:16:08Z	[35mDEBUG[0m	state/state.go:410	processBatch[processBatchRequest.OldLocalExitRoot]: 0x0000000000000000000000000000000000000000000000000000000000000000
2022-08-25T13:16:08Z	[35mDEBUG[0m	state/state.go:411	processBatch[processBatchRequest.EthTimestamp]: 1661433355
2022-08-25T13:16:08Z	[35mDEBUG[0m	state/state.go:412	processBatch[processBatchRequest.Coinbase]: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
2022-08-25T13:16:08Z	[35mDEBUG[0m	state/state.go:413	processBatch[processBatchRequest.UpdateMerkleTree]: 1
2022-08-25T13:16:09Z	[35mDEBUG[0m	state/converters.go:69	ProcessTransactionResponse[TxHash]: 0xb559e861a8dbcdc2d9c62963505782920b0473ebc0a7a7cde65da521721d521d
2022-08-25T13:16:09Z	[35mDEBUG[0m	state/converters.go:70	ProcessTransactionResponse[StateRoot]: 0xd69b140edd763769df5cc4cfa314d7c4055b5a96d15001954a3841965057e7d0
2022-08-25T13:16:09Z	[35mDEBUG[0m	state/converters.go:71	ProcessTransactionResponse[Error]: <nil>
2022-08-25T13:16:09Z	[35mDEBUG[0m	state/converters.go:72	ProcessTransactionResponse[GasUsed]: 519947
2022-08-25T13:16:09Z	[35mDEBUG[0m	state/converters.go:73	ProcessTransactionResponse[GasLeft]: 0
2022-08-25T13:16:09Z	[35mDEBUG[0m	state/converters.go:74	ProcessTransactionResponse[GasRefunded]: 0
2022-08-25T13:16:09Z	[35mDEBUG[0m	state/converters.go:75	ProcessTransactionResponse[IsProcessed]: true
2022-08-25T13:16:09Z	[35mDEBUG[0m	state/state.go:358	ProcessSequencerBatch end
2022-08-25T13:16:09Z	[35mDEBUG[0m	state/state.go:359	*******************************************
```

Then, create a new JSON file under `test/executor/vextors/`, with the following content:

```json
{
    "title": "Fail deploying uniswap repeated nonce",
    "description": "When running the scripts to deploy Uniswap, a repeated nonce got in into a batch. The executor response then was unexpected",
    "genesisFile": "default-genesis.json",
    "batches": [
        {
            "batchL2Data": "processBatch[processBatchRequest.BatchL2Data]",
            "numBatch": "processBatch[processBatchRequest.BatchNum]",
            "oldLocalExitRoot": "processBatch[processBatchRequest.OldLocalExitRoot]",
            "oldStateRoot": "processBatch[processBatchRequest.OldLocalStateRoot]",
            "sequencerAddr": "processBatch[processBatchRequest.Coinbase]",
            "timestamp": "processBatch[processBatchRequest.EthTimestamp]"
        },
        {
            "batchNum 2 ..."
        },
        {
            "batchNum 3..."
        }
    ]
}
```

## Generate genesis

In case some vector doesn't use the default genesis:

1. Start the node without txs:

```bash
make run-db
make run-zkprover
docker compose up -d zkevm-sync
```

2. Get the entries of the merkletree in JSON format: `PGPASSWORD=prover_pass psql -h 127.0.0.1 -p 5432 -U prover_user -d prover_db -c "select row_to_json(t) from (select encode(hash, 'hex') as hash, encode(data, 'hex') as data from state.merkletree) t" > newGenesis.json`
3. Tweak the file until it's a valid json
