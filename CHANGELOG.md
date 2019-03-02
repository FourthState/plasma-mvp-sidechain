# Changelog
All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- [plasmacli] Added keys subcommand with account mapping
- [plasmacli] Added local confirmation signature storage
- Ethereum connection to smart contract
- Implemented Fees
- Unit tests
- Multinode network
- Query sidechain state
- Plasma configuration file
### Changed
- Fixed Length TxBytes (811), compatible with rootchain v1.0.0
- [plasmacli] home flag renamed to directory, flags have suffic "F" for local flags, and "Flag" for persistent flags
- [plasmacli] client keystore/ renamed to store/
- Made UTXO model modular
- Transaction verification to be compatible with rootchain
- Decrease dependency on amino encoding
- Updated client
- Updated documentation
- Upgrade to v0.32.0 of Cosmos SDK

## PreHistory

#### v0.2.0 - July 14th, 2018
- Functional client. Can initialize with genesis UTXOs, start Tendermint daemon, and spend UTXOs using CLI.
- More extensive app tests
- Upgrade to SDK v0.21.0
- Can now retrieve all UTXOs owned by an address.

#### v0.1.1 [HOTFIX] - July 8th, 2018 
- Fix double spend bug when same position is spent twice in single tx
- Added documentation

#### v0.1.0 - June 11th, 2018
- Contains base layer of the blockchain.
- Only capable of validating transactions and updating state.


