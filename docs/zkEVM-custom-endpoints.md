# X1 custom endpoints

The X1 Node JSON RPC server works as is when compared to the official Ethereum JSON RPC, but there are some extra information that also needs to be shared when talking about a L2 Networks, in our case we have information about Batches, Proofs, L1 transactions and much more

In order to allow users to consume this information, a custom set of endpoints were created to provide this information, they are provided under the prefix `zkevm_`

The endpoint documentation follows the [OpenRPC Specification](https://spec.open-rpc.org/) and can be found next to the endpoints implementation as a json file, [here](../jsonrpc/endpoints_zkevm.openrpc.json)

The spec can be easily visualized using the oficial [OpenRPC Playground](https://playground.open-rpc.org/), just copy and paste the json content into the playground area to find a friendly UI showing the methods
