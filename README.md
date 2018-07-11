# Plasma MVP Sidechain

[![license](https://img.shields.io/github/license/FourthState/plasma-mvp-rootchain.svg)](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/LICENSE)

Branch    | Tests | Coverage
----------|-------|----------
develop   | [![Build Status](https://travis-ci.org/FourthState/plasma-mvp-sidechain.svg?branch=develop)](https://travis-ci.org/FourthState/plasma-mvp-sidechain) | [![codecov](https://codecov.io/gh/FourthState/plasma-mvp-sidechain/branch/develop/graph/badge.svg)](https://codecov.io/gh/FourthState/plasma-mvp-sidechain)
master	  | [![Build Status](https://travis-ci.org/FourthState/plasma-mvp-sidechain.svg?branch=master)](https://travis-ci.org/FourthState/plasma-mvp-sidechain) | [![codecov](https://codecov.io/gh/FourthState/plasma-mvp-sidechain/branch/master/graph/badge.svg)](https://codecov.io/gh/FourthState/plasma-mvp-sidechain)

This is the latest [Minimum Viable Plasma](https://ethresear.ch/t/minimal-viable-plasma/426) version.  

We have decided to modify our implementation by removing confirmation signatures as described in the [More Viable Plasma](https://ethresear.ch/t/more-viable-plasma/2160) construction.

**Note**: This sidechain is being constructed to be compatible with our [rootchain contract](https://github.com/FourthState/plasma-mvp-rootchain/master)  

## Overview
As a layer 2 scaling solution, Plasma has two major components: verification and computation. Verification is handled by the rootchain contract which resolves any disputes and distributes funds accordingly. 

Computation is handled off chain by a sidechain. This sidechain leverages the Cosmos SDK to create a scalable and flexible blockchain, that can maintain it's security through reporting merkle roots to the root chain. We will be using [Tendermint](https://github.com/tendermint/tendermint) for consensus on this blockchain. 

We are using a UTXO model for this blockchain. This allows us to do secure and compact proofs when interacting with the rootchain contract. 

## Starting a sidechain

In order to run a sidechain with tendermint consensus and a client to form transaction, a plasma node and light client will need to be initialized. 

**Note**: The following assumes you have [golang](https://golang.org/) properly setup and all dependecies have already been installed. See [Contribution Guidelines](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/CONTRIBUTING.md) for more information.

Plasma Node:

- Navigate to `client/plasmad/` directory
- Run `go install` via command line

The plasma node (plasmad) is now installed and can be called from any directory with `plasmad`

Run `plasmad init` via command line to start an instance of a plasma node with a connection to a tendermint validator.

Run `plasmad start` via command line to begin running the plasma node. You should see empty blocks being proposed and committed.

Plasma Light Client:

- Navigate to `client/plasmacli/` directory
- Run `go install` via command line

Use `plasmacli` to run any of the commands for this light client

The light client uses the Ethereum keystore to create and store passphrase encrypted keys in `$HOME/.plasmacli/keys/`

### Plasma Architecture 
See our [research repository](https://github.com/FourthState/plasma-research) for architectural explanations of our Plasma implementation. 

### Documentation
See our [documentation](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/documentation.md)

### Contributing
See our [contributing guidelines](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/CONTRIBUTING.md)
