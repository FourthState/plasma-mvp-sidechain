The eth subcommand acts an interface enabling interaction with the rootchain contract. 
It requires a connection to a full eth node which can be specified in `~/.plasmacli/plasma.toml`
See [example_plasmacli_plasma.toml]() for an example setup.

You must have eth in your account to use any non querying commands.

## Depositing ##
Example usage: 

```
plasmacli eth deposit 1000 acc1
Enter passphrase:
Successfully sent deposit transaction
Transaction Hash: 0x04d2a92c52e4417382c8a1a59ada3aa3a8619bea3ea61d70870ecb4e7bbded30
```

You can use etherscan to check the status of your transaction.

You can also use the query command to check to see if your deposit occured on the rootchain

```
plasmacli eth query deposit --all
Owner: 0xec36ead9c897b609a4ffa5820e1b2b137d454343
Amount: 1000
Nonce: 1
Rootchain Block: 4071013
```

If you know what your deposit nonce should be you can also query in the following manner:

```
plasmacli eth query deposit 1
Owner: 0xec36ead9c897b609a4ffa5820e1b2b137d454343
Amount: 1000
Nonce: 1
Rootchain Block: 4071013
```

## Checking Submitted headers ##


## Exiting ##


## Challenging ##


## Withdraw ##


## Query Rootchain ##

Query rootchain specific information:

```
plasmacli eth query rootchain
Last Committed Block: 0
Contract Balance: 1000
Withdraw Balance: 0
Minimum Exit Bond: 200000
Operator: 0xec36ead9c897b609a4ffa5820e1b2b137d454343
```

