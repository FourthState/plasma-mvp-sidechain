# Changelog
All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- [\#141](https://github.com/FourthState/plasma-mvp-sidechain/pull/141) Added Makefile
- [\#126](https://github.com/FourthState/plasma-mvp-sidechain/pull/126) Added installation script
- **plasmacli:** [\#110](https://github.com/FourthState/plasma-mvp-sidechain/pull/110) Added eth subcommand for rootchain interaction
- **plasmacli:** [\#110](https://github.com/FourthState/plasma-mvp-sidechain/pull/110) Added plasma.toml in .plasmacli/ for rootchain connection configuration
- **plasmacli:** [\#108](https://github.com/FourthState/plasma-mvp-sidechain/pull/108) Added keys subcommand with account mapping
- **plasmacli:** [\#116](https://github.com/FourthState/plasma-mvp-sidechain/pull/116) Added local confirmation signature storage
- **plasmacli:** [\#120](https://github.com/FourthState/plasma-mvp-sidechain/pull/120) Added `include-deposit` command to add deposit nonce into sidechain
- Ethereum connection to smart contract
- Implemented Fees
- Unit tests
- Multinode network
- Query sidechain state
- Plasma configuration file
- Added IncludeDepositMsg with handling to allow explicit deposit inclusion into sidechain
### Changed
- [\#141](https://github.com/FourthState/plasma-mvp-sidechain/pull/141) Dependency management is now handled by go modules instead of Dep
- [\#129](https://github.com/FourthState/plasma-mvp-sidechain/pull/129) Updated sign command to iterate over an account to finalize transactions
- [\#129](https://github.com/FourthState/plasma-mvp-sidechain/pull/129) Updated spend to auto generate transaction for users based on the utxos they own
- [\#120](https://github.com/FourthState/plasma-mvp-sidechain/pull/118) Fixed Length TxBytes (811), compatible with rootchain v1.0.0
- **plasmacli:** [\#108](https://github.com/FourthState/plasma-mvp-sidechain/pull/108) home flag renamed to directory, flags have suffix "F" for local flags, and "Flag" for persistent flags
- **plasmacli:** [\#116](https://github.com/FourthState/plasma-mvp-sidechain/pull/116) client keystore/ renamed to store/
- Made UTXO model modular
- Transaction verification to be compatible with rootchain
- Decrease dependency on amino encoding
- Updated client
- Updated documentation
- Upgrade to v0.32.0 of Cosmos SDK, v0.28.0 of TM
### Fixed
- [\#147](https://github.com/FourthState/plasma-mvp-sidechain/pull/147) Fix Syncing bug where syncing nodes would panic after processing exitted inputs/deposits. Bug is explained in detail here: [\#143](https://github.com/FourthState/plasma-mvp-sidechain/issues/143)
- [\#154](https://github.com/FourthState/plasma-mvp-sidechain/pull/154) Fixes issue where include-Deposit msg.Owner == deposit.Owner not enforced. This is necessary to prevent malicious users from rewriting an already included UTXO in store.
### Deprecated 
- Dep is no longer be supported

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


