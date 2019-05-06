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

Exiting a fee can be done in the same format as exiting a deposit. 
Specifiying the position and committed fee is the only information required to do a successful fee withdrawal. 


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
plasmacli prove acc1 "(22.0.1.0)"
Roothash: 0xC6BA74C556C3114598214AC828766DC485E688F217B1506C33F2095045B0300E
Total: 1
LeafHash: 0xc6ba74c556c3114598214ac828766dc485e688f217b1506c33f2095045b0300e
TxBytes: 0xf90328f9029da0000000000000000000000000000000000000000000000000000000000000000ba00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000b882b0452100c6e01ab4e04e44bc5fd767dbcaa8abf930585be257d9e33dfab9b8230d39e48c3a350d3727237efe9c70d24e2e769516b930b56eea5188bec35065a7010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000005b88200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000945475b99e01ac3bb08b24fd754e2868dbb829bc3aa0000000000000000000000000000000000000000000000000000000000000232894ec36ead9c897b609a4ffa5820e1b2b137d454343a000000000000000000000000000000000000000000000000000000000000007d0a00000000000000000000000000000000000000000000000000000000000000000f886b8415a1ba592dc188288fedd7dfd86cfd953e9993a0e50948fa230c4a6b33b71ecb87a26c15322034c3db1f58be930a25ef1565461cee75bda4103624f9bd30f6a1f00b8415a1ba592dc188288fedd7dfd86cfd953e9993a0e50948fa230c4a6b33b71ecb87a26c15322034c3db1f58be930a25ef1565461cee75bda4103624f9bd30f6a1f00

plasmacli eth exit acc1 "(22.0.1.0)" -b 0xf90328f9029da0000000000000000000000000000000000000000000000000000000000000000ba00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000b882b0452100c6e01ab4e04e44bc5fd767dbcaa8abf930585be257d9e33dfab9b8230d39e48c3a350d3727237efe9c70d24e2e769516b930b56eea5188bec35065a7010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000005b88200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000945475b99e01ac3bb08b24fd754e2868dbb829bc3aa0000000000000000000000000000000000000000000000000000000000000232894ec36ead9c897b609a4ffa5820e1b2b137d454343a000000000000000000000000000000000000000000000000000000000000007d0a00000000000000000000000000000000000000000000000000000000000000000f886b8415a1ba592dc188288fedd7dfd86cfd953e9993a0e50948fa230c4a6b33b71ecb87a26c15322034c3db1f58be930a25ef1565461cee75bda4103624f9bd30f6a1f00b8415a1ba592dc188288fedd7dfd86cfd953e9993a0e50948fa230c4a6b33b71ecb87a26c15322034c3db1f58be930a25ef1565461cee75bda4103624f9bd30f6a1f00
Sent exit transaction
Transaction Hash: 0xeea2e9e6ff8f93ba189d938cf531052c55f724cce87c38895d3eeec20299e615
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
plasmacli spend 0x5475b99e01ac3bb08b24fd754e2868dbb829bc3a 9000 acc1 --position "(0.0.0.6)" --fee 1000

plasmacli eth exit acc1 "(0.0.0.6)" --fee 0
Enter passphrase:
Sent deposit exit transaction
Transaction Hash: 0x6b23f7fc7cfdd5b4948e78740f38220133392e77b99b535984faafcc380567a8

plasmacli eth query exit --position "(0.0.0.6)"
Owner: 0xec36ead9c897b609a4ffa5820e1b2b137d454343
Amount: 10000
State: Pending
Committed Fee: 0
Created: 2019-03-26 20:17:36 +0000 UTC

Exit will be finalized in about: 167.99566772949944 hours

plasmacli eth challenge "(0.0.0.6)" "(25.0.0.0)" --tx-bytes 0xf90328f9029da00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000006b88200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000b88200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000945475b99e01ac3bb08b24fd754e2868dbb829bc3aa00000000000000000000000000000000000000000000000000000000000002328940000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000003e8f886b8412efef4db94eaf2186e9cd51aacba30876596a2d467ed954c06e348925bc60611028ec63150039305d32ebe8081511f1b8542935758c841be5aa424d4d7f7f61200b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
Enter passphrase:
Sent challenge transaction
Transaction Hash: 0xbf7ea33c09a187eeb6bb9172e0459c9a137058b93a56ea869d206e615648cc5f

plasmacli eth query exit --position "(0.0.0.6)"
Owner: 0xec36ead9c897b609a4ffa5820e1b2b137d454343
Amount: 10000
State: Nonexistent
Committed Fee: 0
Created: 2019-03-26 20:17:36 +0000 UTC
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

Balance before finalizing:
```
plasmacli eth query balance acc1
Rootchain Balance: 200000
```

Finalize exits:
```
plasmacli eth finalize acc1
Enter passphrase:
Successfully sent finalize exits transaction
Transaction Hash: 0x95d028a4ae3e90833ab222880b366e32818353195acbb7e9c7516972d224c1f6
```

Balance after finalizing:
```
plasmacli eth query balance acc1
Rootchain Balance: 618000
```

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

