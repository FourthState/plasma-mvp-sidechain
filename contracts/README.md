# PLASMA MVP

[![travis build](https://travis-ci.org/FourthState/plasma-mvp-rootchain.svg?branch=master)](https://travis-ci.org/FourthState/plasma-mvp-rootchain)
[![license](https://img.shields.io/github/license/FourthState/plasma-mvp-rootchain.svg)](https://github.com/FourthState/plasma-mvp-rootchain/blob/master/LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/FourthState/plasma-mvp-rootchain/badge.svg?branch=master)](https://coveralls.io/github/FourthState/plasma-mvp-rootchain?branch=master)

Implementation of [Minimum Viable Plasma](https://ethresear.ch/t/minimal-viable-plasma/426)

## Overview
Plasma is a layer 2 scaling solution which conducts transaction processing off chain and allows for only merkle roots of each block to be reported to a root chain. This allows for users to benefit from off chain scaling while still relying on decentralized security.

The root contract of a Plasma child chain represents an intermediary who can resolve any disputes. The root contract is responsible for maintaining a mapping from block number to merkle root, processing deposits, and processing withdrawals.

## Root Contract Details
A transaction is encoded in the following form:

```
[
  [Blknum1, TxIndex1, Oindex1, DepositNonce1, Owner1, Input1ConfirmSig,

   Blknum2, TxIndex2, Oindex2, DepositNonce2, Owner2, Input2ConfirmSig,

   NewOwner, Denom1, NewOwner, Denom2, Fee],

  [Signature1, Signature2]
]
```
The signatures are over the hash of the transaction list signed by the owner of each respective utxo input.

### Documentation

See our [documentation](https://github.com/FourthState/plasma-mvp-rootchain/blob/master/docs/rootchainFunctions.md) for a more detailed description of the smart contract functions.

### Testing
1. ``git clone https://github.com/fourthstate/plasma-mvp-rootchain``
2. ``cd plasma-mvp-rootchain``
3. ``npm install``
4. ``npm install -g truffle ganache-cli``  // if not installed already
5. ``ganache-cli`` // run as a background process
6. ``npm test``

### Running
The first migration file `1_initial_migration` deploys the `PriorityQueue` library and links it to the `RootChain` contract, while the second one `2_deploy_rootchain` finally makes the deployment. Ethereum requires libraries to already be deployed prior to be used by other contracts.

If you encounter problems, make sure your local test rpc (e.g. [ganache](https://github.com/trufflesuite/ganache-core)) has the same network id as the contract's json from the `build` folder.

### Contributing

See our [contribution guidelines](https://github.com/FourthState/plasma-mvp-rootchain/blob/master/CONTRIBUTING.md). Join our [Discord Server](https://discord.gg/YTB5A4P).
