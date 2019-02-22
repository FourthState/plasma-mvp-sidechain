#Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Ethereum connection to smart contract
- Implemented Fees
- Unit tests
- Multinode network
- Query sidechain state
- Plasma configuration file
### Changed
- Made UTXO model modular
- Transaction verification to be compatible with rootchain
- Decrease dependency on amino encoding
- Updated client
- Updated documentation

## PreHistory

##### [0.2.0]
- Functional client. Can initialize with genesis UTXOs, start Tendermint daemon, and spend UTXOs using CLI.
- More extensive app tests
- Upgrade to SDK v0.21.0
- Can now retrieve all UTXOs owned by an address.

##### [0.1.1] HOTFIX - July 8th, 2018 
- Fix double spend bug when same position is spent twice in single tx
- Added documentation

##### [0.1.0] - June 11th, 2018
- Contains base layer of the blockchain.
- Only capable of validating transactions and updating state.


