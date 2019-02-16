# Plasma MVP Sidechain

[![license](https://img.shields.io/github/license/FourthState/plasma-mvp-rootchain.svg)](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/LICENSE)

Branch    | Tests | Coverage
----------|-------|----------
develop   | [![Build Status](https://travis-ci.org/FourthState/plasma-mvp-sidechain.svg?branch=develop)](https://travis-ci.org/FourthState/plasma-mvp-sidechain) | [![codecov](https://codecov.io/gh/FourthState/plasma-mvp-sidechain/branch/develop/graph/badge.svg)](https://codecov.io/gh/FourthState/plasma-mvp-sidechain)
master	  | [![Build Status](https://travis-ci.org/FourthState/plasma-mvp-sidechain.svg?branch=master)](https://travis-ci.org/FourthState/plasma-mvp-sidechain) | [![codecov](https://codecov.io/gh/FourthState/plasma-mvp-sidechain/branch/master/graph/badge.svg)](https://codecov.io/gh/FourthState/plasma-mvp-sidechain)

This is the latest [Minimum Viable Plasma](https://ethresear.ch/t/minimal-viable-plasma/426) version.  

**Note**: This sidechain is being constructed to be compatible with our [rootchain contract](https://github.com/FourthState/plasma-mvp-rootchain)  

## What is Plasma?
Plasma has two major components: verification and computation. 
Verification is handled by the rootchain smart contract, which resolves any disputes and distributes funds accordingly. 

Computation is handled separately by a sidechain, which leverages the Cosmos SDK to create a modular and flexible blockchain.
This sidechain maintains its security through reporting proofs via merkle roots to the Ethereum mainchain. 
We will be using [Tendermint](https://github.com/tendermint/tendermint) as a consensus algorithm.

Plasma MVP utilizes a UTXO model, which allows for secure and compact proofs when interacting with the rootchain smart contract. 

Learn more about plasma on [learnplasma.org](https://www.learnplasma.org/en/)!

## Quick Start

**Requirements**: 
- [golang](https://golang.org/)
- [dep](https://github.com/golang/dep)

Pull the latest version of the develop branch.
Run `dep ensure -vendor-only`

Plasma Node:

- Navigate to `server/plasmad/` directory
- Run `go install` via command line

The plasma node (plasmad) is now installed and can be called from any directory with `plasmad`

Run `plasmad init` via command line to start an instance of a plasma node with a connection to a tendermint validator.
Use the `--directory <dirpath>` to specify a location where you want your plasma node to exist.  

Navigate to `<dirpath>/config/` (default is `$HOME/.plasmad/config`), set configuration parameters in config.toml and plasma.toml.
Run `plasmad start` via command line to begin running the plasma node. 

Plasma Light Client:

- Navigate to `client/plasmacli/` directory
- Run `go install` via command line

Use `plasmacli` to run any of the commands for this light client

### dep ensure 
When building the sidechain, go dep is used to manage dependencies. 
Running `dep ensure` followed by `go build` will result in the following output:

```
# github.com/FourthState/plasma-mvp-sidechain/vendor/github.com/ethereum/go-ethereum/crypto/secp256k1
../vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/curve.go:42:44: fatal error: libsecp256k1/include/secp256k1.h: No such file or directory
```
This is caused by a go dep issue outlined [here](https://github.com/tools/godep/issues/422).
To fix this locally, add the following in Gopkg.lock under `crypto/secp256k1` and above `crypto/sha3`:

```
"crypto/secp256k1/libsecp256k1",
"crypto/secp256k1/libsecp256k1/include",
"crypto/secp256k1/libsecp256k1/src",
"crypto/secp256k1/libsecp256k1/src/modules/recovery",
```

Run `dep ensure -vendor-only`

Your vendor folder should now contain all the necessary dependencies, there is no need to run `dep ensure`. 
  
### Plasma Architecture 
See our [research repository](https://github.com/FourthState/plasma-research) for architectural explanations of our Plasma implementation. 

### Documentation
See our [documentation](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/docs/overview.md)

### Contributing
See our [contributing guidelines](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/CONTRIBUTING.md). Join our [Discord Server](https://discord.gg/YTB5A4P).
