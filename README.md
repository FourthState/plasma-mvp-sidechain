# Plasma MVP Sidechain
[![Go Report](https://goreportcard.com/badge/github.com/FourthState/plasma-mvp-sidechain)](https://goreportcard.com/report/github.com/FourthState/plasma-mvp-sidechain)
[![Build Status](https://travis-ci.org/FourthState/plasma-mvp-sidechain.svg?branch=develop)](https://travis-ci.org/FourthState/plasma-mvp-sidechain)
[![codecov](https://codecov.io/gh/FourthState/plasma-mvp-sidechain/branch/develop/graph/badge.svg)](https://codecov.io/gh/FourthState/plasma-mvp-sidechain)
[![Discord](https://img.shields.io/badge/discord-join%20chat-blue.svg)](https://discord.gg/YTB5A4P)
[![license](https://img.shields.io/github/license/FourthState/plasma-mvp-rootchain.svg)](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/LICENSE)

Implementation of [Minimum Viable Plasma](https://ethresear.ch/t/minimal-viable-plasma/426) compatible with our [rootchain contract](https://github.com/FourthState/plasma-mvp-rootchain)  

## What is Plasma?
Plasma has two major components: verification and computation. 
Verification is handled by the rootchain contract, which resolves any disputes and distributes funds accordingly.
Computation is handled separately by a sidechain, which maintains its security through reporting proofs via merkle roots to the rootchain contract. 

Plasma MVP utilizes a UTXO model, which allows for secure and compact proofs. Learn more about plasma on [learnplasma.org](https://www.learnplasma.org/en/)!

We are using [Tendermint](https://github.com/tendermint/tendermint) for our consensus protocol.
This sidechain currently supports a single validator, but will be updated in the future to support multiple validators.

## Quick Start

### Install using a script

This script can be used on a fresh server that has no dependencies installed.

```
curl https://raw.githubusercontent.com/FourthState/plasma-mvp-sidechain/develop/scripts/plasma_install.sh > install.sh
chmod +x install.sh
./install.sh
```

### Manual Install

**Requirements**: 
- [golang](https://golang.org/)
- [dep](https://github.com/golang/dep)

Pull the latest version of the develop branch.
Run `dep ensure -vendor-only`

***Plasma Node:***

```
cd server/plasmad/
go install
```

Run `plasmad init` to start an instance of a plasma node.
Use the `--home <dirpath>` to specify a location where you want your plasma node to exist.

Navigate to `<dirpath>/config/` (default is `$HOME/.plasmad/config`), set configuration parameters in config.toml and plasma.toml.
Run `plasmad start` to begin running the plasma node. 

***Plasma Client:***

```
cd client/plasmacli/ 
go install
```

Navigate to `$HOME/.plasmacli`, set ethereum client configuration parameters in plasma.toml.
Use `plasmacli` to run any of the commands for this light client
  
### Plasma Architecture 
See our [research repository](https://github.com/FourthState/plasma-research) for architectural explanations of our Plasma implementation. 

### Documentation
See our [documentation](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/docs/overview.md)

### Contributing
See our [contributing guidelines](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/.github/CONTRIBUTING.md). Join our [Discord Server](https://discord.gg/YTB5A4P).
