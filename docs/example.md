# Using the Sidechain Example #

The following assumes you have already deployed the rootchain contract to either ganache or a testnet.
See our rootchain deployment [example]()

Plasmacli: the command-line interface for interacting with the sidechain and rootchain. 

Plasmad: runs a sidechain full-node

## Setting up a full-node ##

Install the latest version of plasmad: 

```
cd server/plasmad/
go install
```

Run `plasmad init` to initalize a validator. cd into `~/.plasmad/config`. 
Open genesis.json and add an ethereum address to `fee_address`. 
See our example [genesis.json](https://github.com/FourthState/plasma-mvp-sidechain/blob/develop/docs/testnet-setup/example_genesis.json)

Open config.toml and add any configurations you would like to add for your validator, such as a moniker. TODO: add section on seeds

Open plasma.toml, set `is_operator` to true if you are running a validator. 
Set `ethereum_operator_privatekey` to be the unencrypted private key that will be used to submit blocks to the rootchain.
It must contain sufficient eth to pay gas costs for every submitted plasma block.
Set `ethereum_plasma_contract_address` to be the contract address of the deployed rootchain. 
Set `plasma_block_commitment_rate` to be the rate at which you want plasma blocks to be submitted to the rootchain. 
Set `ethereum_nodeurl` to be the url which contains your ethereum full node. 
Set `ethereum_finality` to be the number of ethereum blocks until a submitted header is presumed to be final.

See our example [plasma.toml](https://github.com/FourthState/plasma-mvp-sidechain/blob/develop/docs/testnet-setup/example_plasma.toml)

Run `plasmad unsafe-reset-all` followed by `plasmad start`

You should be successfully producing empty blocks

Things to keep in mind: 
- You can change `timeout_commit` in config.toml to slow down block time. 
- go install `plasmacli` and `plasmad` when updating to newer versions
- Using `plasmad unsafe-reset-all` will erase all chain history. You will need to redeploy the rootchain contract. 

## Setting up the client ##

You will need to run a full eth node to interact with the rootchain contract.
See the install [script](https://github.com/FourthState/plasma-mvp-sidechain/blob/develop/scripts/plasma_install.sh) for an example of setting up a full eth node.

Install the latest version of plasmacli:

```
cd client/plasmacli/
go install
```

cd into `~/.plasmacli/`. Open plasma.toml.
Set `ethereum_plasma_contract_address` to be the contract address of the deployed rootchain. 
Set `ethereum_nodeurl` to be the url which contains your ethereum full node. 
Set `ethereum_finality` to be the number of ethereum blocks until a submitted header is presumed to be final.

Things to keep in mind:
- plasmacli can be used without a full node, but certain features will be disabled such as interacting with the rootchain
- Using the `-h` will provide short documentation and example usage for each command 

See [keys documentation](https://github.com/FourthState/plasma-mvp-sidechain/blob/develop/docs/keys.md) for examples on how to use the keys subcommand.

See [eth documentation]() for examples on how to use the eth subcommand.

## Spending Deposits ## 

In order to spend a deposit on the sidechain, first a user must deposit on the rootchain and then send an include-deposit transaction (after presumed finality).
A user can deposit using the eth subcommand. See this [example]()

Sending an include-deposit transaction: 
```
plasmacli include-deposit 1 acc1
```

You can also use the --sync flag
```
plasmacli include-deposit 1 acc1 --sync
Error: broadcast_tx_commit: Response error: RPC error -32603 - Internal error: Error on broadcastTxCommit: Tx already exists in cache
```

The above error simply means that the above transaction has been sent but not yet included in a block.

If you query your account balance you should see your deposit:
```
plasmacli query balance acc1
Position: (0.0.0.1) , Amount: 1000
Total: 1000
```

**Note:** The include-deposit transaction will fail if presumed finality specified by the validator has not yet been reached.

Spending the deposit:
First argument is address being sent to, followed by amounts to send (first output, second output), followed by the account to send from. 
Position flag of inputs to be spent must be specified. 

```
 plasmacli spend 0x5475b99e01ac3bb08b24fd754e2868dbb829bc3a 1000,0 acc1 --position "(0.0.0.1)"
```

The above address being sent to corresponds to acc2

```
plasmacli query balance acc2
Position: (2.0.0.0) , Amount: 1000
Total: 1000
```

Deposits and Fees do not need a confirmation signature to be spent. 


## Spending UTXOS ##

