The eth subcommand acts an interface enabling interaction with the rootchain contract. 
It requires a connection to a full eth node which can be specified in `~/.plasmacli/plasma.toml`
See an example [plasma.toml](https://github.com/FourthState/plasma-mvp-sidechain/blob/develop/docs/testnet-setup/example_plasmacli_plasma.toml) for an example setup.

You must have eth in your account to use the commands deposit, exit, challenge, withdraw, and finalize.

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

```
plasmacli eth query block 1
Block: 1
Header: 0x5a4f97a64e82a4aa090b0946ed299c2d75f4c8353be0a4b97df8607120713183
Txs: 1
Fee: 0
Created: 2019-03-21 19:03:55 +0000 UTC
```

## Exiting ##

UTXO's and unspent deposits can be exitted. 
When exiting, the user can use the "trust-node" flag if they trust the full node specified by the "node" flag.
When "trust-node" flag is used, information necessary for exiting will be retireved from the connected full node. 
Exiting a deposit, only requires its position and committed fee so no flags are necessary. 
A proof is not required for transactions included in a block of size 1.

Exiting an unspent deposit:

```
plasmacli eth exit acc1 "(0.0.0.3)"
Enter passphrase:
Sent deposit exit transaction
Transaction Hash: 0xa16f3909e1dd749f5093199e7b597974f5f643d9b833b3a9546ff674dac0af28
```

Exiting a utxo with trust-node:
```
plasmacli eth exit acc1 "(10.0.0.0)" -t
Enter passphrase:
Warning: No proof was found or provided. If the exiting transaction was not the only transaction included in the block then this transaction will fail.
Sent exit transaction
Transaction Hash: 0x4df9c79d17036b7468b69f83a313a1fb87ccc987e9e28be1e78f81d2d19afab8
```

Exiting a utxo without trust-node:

```

```

Querying for deposit exits:

```
plasmacli eth query exit --deposits --all
Owner: 0xec36ead9c897b609a4ffa5820e1b2b137d454343
Amount: 1000
State: Pending
Committed Fee: 0
Created: 2019-03-21 22:05:40 +0000 UTC

Exit will be finalized in about: 167.40055942302556 hours

Owner: 0x5475b99e01ac3bb08b24fd754e2868dbb829bc3a
Amount: 10000
State: Challenged
Committed Fee: 0
Created: 2019-03-21 22:21:40 +0000 UTC
```

Querying for transaction exits:

```
plasmacli eth query exit --all
Owner: 0xec36ead9c897b609a4ffa5820e1b2b137d454343
Amount: 10000
State: Pending
Committed Fee: 0
Created: 2019-03-21 22:40:10 +0000 UTC

Exit will be finalized in about: 167.96671164467028 hours
```

Querying for a specific exit:

```
plasmacli eth query exit --position "(10.0.0.0)"
Owner: 0xec36ead9c897b609a4ffa5820e1b2b137d454343
Amount: 10000
State: Pending
Committed Fee: 0
Created: 2019-03-21 22:40:10 +0000 UTC

Exit will be finalized in about: 167.9601120748311 hours
```

## Challenging ##

A pending exit may be challenged if the exit committed to an incorrect fee amount or if the utxo was spent on the sidechain in a finalized transaction.
Every exit commits to the fee of an unfinalized spend of that deposit/utxo. 
If the deposit/utxo was never spent, the committed fee is 0. 
If the deposit/utxo was involved in an unfinalized spend which included a non zero fee, then the exit must commit to that non zero fee or risk being challenged. 
A deposit/utxo may also be challenged with a finalized spend, if it exists.  

Challenging a deposit with an incorrect committed fee:


```

```

Challenging a deposit with a finalized spend:

Deposit '(0.0.0.4)' was exitted and spent on the sidechain.

```
plasmacli eth challenge "(0.0.0.4)" "(10.0.0.0)" acc1 -t --signatures 0x19984e40ce233d31db3a5bbf724f079e306644046c0d30bdbd93027bd3e3c04f314c41f14a8b3c771e1faad6c41a07ed84d418123a8c2a0e82237eb8a22ca62501
Enter passphrase:
Warning: No proof was found or provided. If the exiting transaction was not the only transaction included in the block then this transaction will fail.
Sent challenge transaction
Transaction Hash: 0x6347965d4cb1af2160ff56f53cce18fdd50e0be8bbbc7f03ff475139f8422d8b
```

**Note:** 

If your local signature storage contains the confirmation signature, then the "signature" flag is unnecessary. 
If "trust-node" is not used, "proof" and "tx-bytes" flags are required. 

Challenging a transaction exit with an incorrect committed fee:

The following exiting utxo was involved in a spend that committed a fee of 1000, but exitted with a committed fee of 0.

```
plasmacli eth challenge "(22.0.0.0)" "(24.0.0.0)" acc1 -t 
Enter passphrase:
Warning: No proof was found or provided. If the exiting transaction was not the only transaction included in the block then this transaction will fail.
Sent challenge transaction
Transaction Hash: 0xdcf512c38b12670e36914fc7f2129e246e0ba61aa0abd1c284077e364c36ae1b

plasmacli eth query exit --position "(22.0.0.0)"
Owner: 0x5475b99e01ac3bb08b24fd754e2868dbb829bc3a
Amount: 9000
State: Nonexistent
Committed Fee: 0
Created: 2019-03-22 20:09:55 +0000 UTC
```

Since the exiting utxo can still be exitted with the correct fee amount, its state is set to Nonexistent.

Challenging a transaction exit with a finalized spend:

Now the utxo exits with the correct fee amount, but the spend in block 24 is finalized so we can challenge it again. 

```
plasmacli eth exit acc2 "(22.0.0.0)" --fee 1000
Enter passphrase:
Warning: No proof was found or provided. If the exiting transaction was not the only transaction included in the block then this transaction will fail.
Sent exit transaction
Transaction Hash: 0xef7a0798bcafc4b575d4005e0dedacc69c4bfe3c37e13fbbb71397e32c99bc05

plasmacli sign acc2

UTXO
Position: (24.0.0.0)
Owner: 0xec36ead9c897b609a4ffa5820e1b2b137d454343
Value: 15000
> Would you like to finalize this transaction? [Y/n]
Y
Enter passphrase:
Confirmation Signature for output with position: (24.0.0.0)
0x151367f8b7ab02d4a616e175f3a6d5955306933fe54dce6d5247f45a543db071798fbd61acb254cfb02cc53d72c19d7ce0903c5fe3bd0570f98a11ed26dc83e701

UTXO
Position: (24.0.0.0)
Owner: 0xec36ead9c897b609a4ffa5820e1b2b137d454343
Value: 15000
> Would you like to finalize this transaction? [Y/n]
Y 
Enter passphrase:
Confirmation Signature for output with position: (24.0.0.0)
0x151367f8b7ab02d4a616e175f3a6d5955306933fe54dce6d5247f45a543db071798fbd61acb254cfb02cc53d72c19d7ce0903c5fe3bd0570f98a11ed26dc83e701

plasmacli eth challenge "(22.0.0.0)" "(24.0.0.0)" acc1
Enter passphrase:
Warning: No proof was found or provided. If the exiting transaction was not the only transaction included in the block then this transaction will fail.
Sent challenge transaction
Transaction Hash: 0x5a4633533895bcd65514975b755509d452cc283bf5f299cbe1224999a66b00e7


```

## Finalize ##

Exits may be finalized, after the challenge period has ended. 
The default challenge period is 5 days for deposits and 7 days for transaction exits

## Withdraw ##

Checking the balance on the rootchain that can be withdrawn:

```
plasmacli eth query balance acc1
Rootchain Balance: 200000
```

Withdrawing an entire balance from the rootchain:

```
plasmacli eth withdraw acc1
Enter passphrase:
Successfully sent withdraw transaction
Transaction Hash: 0xd790a512f15051ee866bf8751b9913de3d2fedae0e2228b5c0ff97ed6451f2e7
```

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

