# Plasma MVP Sidechain

[![Build Status](https://travis-ci.org/FourthState/plasma-mvp-sidechain.svg?branch=master)](https://travis-ci.org/FourthState/plasma-mvp-sidechain)
[![license](https://img.shields.io/github/license/FourthState/plasma-mvp-rootchain.svg)](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/LICENSE)

We're implementing [Minimum Viable Plasma](https://ethresear.ch/t/minimal-viable-plasma/426) 

**Note**: This sidechain is being constructed to be compatible with our [rootchain contract](https://github.com/FourthState/plasma-mvp-rootchain/master)  

## Overview
As a layer 2 scaling solution, Plasma has two major components: verification and computation. Verification is handled by the rootchain contract which resolves any disputes and distributes funds accordingly. 

Computation is handled off chain by a sidechain. This sidechain levarges the Cosmos SDK to create a scalable and flexible blockchain, that can maintain it's security through reporting merkle roots to the root chain. We will be using [Tendermint](https://github.com/tendermint/tendermint) for consensus on this blockchain. 

We are using a UTXO model for this blockchain. This allows us to do secure and compact proofs when interacting with the rootchain contract. 

### Plasma Architecture 
See our [research repository](https://github.com/FourthState/plasma-research) for architectural explanations of our Plasma implementation. 

### Documentation
See our [documentation](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/documentation.md)

### Contributing
See our [contributing guidelines](https://github.com/FourthState/plasma-mvp-sidechain/blob/master/CONTRIBUTING.md)
