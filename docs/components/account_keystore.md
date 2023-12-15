# Generating an Account Keystore file:

This file contains your Ethereum L1 private key, but it will be encrypted at rest using a password of your choice. The ZKEVM Node - depending on which operating mode it's set up - will use this file in conjunction with the password to authorize L1 transactions.

```bash
docker run --rm okx/x1-node:latest sh -c "/app/x1-node encryptKey --pk=[your private key] --pw=[password to encrypt file] --output=./keystore; cat ./keystore/*" > account.keystore
```

**NOTE**:

- Replace `[your private key]` with your Ethereum L1 account private key
- Replace `[password to encrypt file]` with a password used for file encryption. This password must be passed to the Node later on via env variable (`X1_NODE_ETHERMAN_PRIVATEKEYPASSWORD`)
