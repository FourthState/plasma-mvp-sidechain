# Using the Sidechain Example #

## Setup a testnet ##

cd into client/plasmacli

run `go install`

cd into server/plasmad

run `go install`

run `plasmad init` to initalize a validator. cd into `~/.plasmad/config`. Open genesis.json and add an ethereum address to `fee_address`. See our example [genesis.json](https://github.com/FourthState/plasma-mvp-sidechain/blob/develop/docs/testnet-setup/example_genesis.json)

Open config.toml and add any configurations you would like to add for your validator, such as a moniker.

Open plasma.toml, set `is_operator` to true if you are running a validator. Set `ethereum_operator_privatekey` to be the unencrypted private key that will be used to submit blocks to the rootchain. It must contain sufficient eth to pay gas costs for every submitted plasma block. Set `ethereum_plasma_contract_address` to be the contract address of the deployed rootchain. Set `plasma_block_commitment_rate` to be the rate at which you want plasma blocks to be submitted to the rootchain. Set `ethereum_nodeurl` to be the url which contains your ethereum full node. Set `ethereum_finality` to be the number of ethereum blocks until a submitted header has presumed finality.

run `plasmad unsafe-reset-all` followed by `plasmad start`

You should be successfully producing empty blocks

Things to keep in mind: 
- You can change `timeout_commit` in config.toml to slow down block time. 
- go install `plasmacli` and `plasmad` when updating to newer versions
- Use `plasmad unsafe-reset-all` if you encounter an unexpected error. If the error persists open an issue.

## Generating Keys ##

In order to spend utxos on the sidechain, we will need keys corresponding to the addresses that own those utxos. We can use plasmacli to generate these keys.

For example:

```
plasmacli add
Enter a passphrase for your key:
Repeat the passphrase:

**Important** do not lose your passphrase.
It is the only way to recover your account
You should export this account and store it in a secure location
Your new account address is: 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2
```

You can also import unencrypted private keys:

Ganache-cli display

```
 Ganache CLI v6.1.8 (ganache-core: 2.2.1)

Available Accounts
==================
(0) 0xea6ed4bb7cba09c391c11a15d5472e806caa3986 (~100 ETH)
(1) 0x6a06b2fd816021568b5b8abad00eab4679ab3450 (~100 ETH)
(2) 0x864fe3ffbde7cd14e4feb309df2d2b99674af22d (~100 ETH)
(3) 0x2ad263e223db255547fef5796bc451db3fa6714a (~100 ETH)
(4) 0x3f4acad1741381594bceec558edce934865ea61f (~100 ETH)
(5) 0xc1f178dfc9c05f82614563304455f8ac61aea4e9 (~100 ETH)
(6) 0xd28d45a70a750f90e8fb045cacbe11470742fdf8 (~100 ETH)
(7) 0xc7b1e3daaeed0c4a6032d330df388aa3412057d6 (~100 ETH)
(8) 0xd28b24bd1ebb3de7cd519c7e5aca955864a31ecd (~100 ETH)
(9) 0x856277d5c6cae04203a740cf7ee0758a46572033 (~100 ETH)

Private Keys
==================
(0) 0x46b7f5573b110d21ca809d1e17d97c582e6e402b1d7daa7d09264c13015e73ae
(1) 0x2a21b4bcbc0d812adf41ff098f336a86325e3c7df31d5c576ac7b760ed2d56f0
(2) 0x1d44ea210bf3bd1d074e46dcf4c23f57e3e1dba044a943dcc74d65de1599454f
(3) 0x40cbbddda351a5bba39f2951534f8ad5dd3830e81fd3e85375457cb2de02b132
(4) 0xbee7332633f9fdf6b6c1bd902425d56085e3d94f608facb407b2e934854b7714
(5) 0x2aeaf0e4256104a5dcb214e6ac84ba1225bcbdd038d60bfb3b92f50e0d0b1815
(6) 0xc7ad5cd86bfdeae63365ad54283044b900b9394ff9708e1a48b0ff5bfd0ed947
(7) 0xb4594bdb60bedc21cbcd507a2842bd6a93e8fff719d1e987e84360644b80b747
(8) 0x5e64a2a9dd2745b34bc8a6ed12d6bfd92ec75c0c83853329696e3cc17d77c266
(9) 0x77984146b9c6628b8f54661ab97462fda52cf1d5cadcba91475ecc4118f5a994

```

Put the private key you wish to import in a file

keyfile:
```
46b7f5573b110d21ca809d1e17d97c582e6e402b1d7daa7d09264c13015e73ae
```

Use plasmacli
```
plasmacli import ~/mykey
Please set a passphrase for your imported account.
Passphrase:
Repeat the passphrase:
Address: {ea6ed4bb7cba09c391c11a15d5472e806caa3986}
```


#### List Keys ####

```
plasmacli list keys
Index:		Address:
0		0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2
1		0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986
```

## Spending Deposits/Fees ## 

Deposits and Fees do not need a confirmation signature to be spent. 

```
plasmacli balance 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2
Position: {0 0 0 1} 
Amount: 1000 
```

In the following example, the deposit at position {0 0 0 1} is spent, sending 500 to `0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986`, 400 to back to the spending address, and 100 to the validator as a fee for the transaction.
```
plasmacli send \
--address 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 \
--position 0.0.0.1 \
--amounts 500,400,100 \
--to 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986,0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2
Password to sign with '0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2':
Committed at block 4. Hash B9418E36D2D0A629A3D18C54CDC04E063B9D5285C443116B7156D1C0DB61E08D
``` 

--address:  used to identify the addresses that own the inputs. In this example we only need the owner of the deposit being spent.

--position: the position of the inputs being spent

--amounts: If you are sending to a single address you would have amount1,0,fee. If you are sending to two addresses amount1,amount2,fee. Amount1 + Amount2 + Fee must equal (Amount of Input0 + Amount of Input1)

--to: The addresses you are sending to

Each transaction can have a maximum of 2 inputs and 2 outputs


#### Spending the fee ####

```
plasmacli balance 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986
Position: {4 0 0 0} 
Amount: 500 
Position: {4 65535 0 0} 
Amount: 100 
Position: {0 0 0 2} 
Amount: 1000 
Position: {0 0 0 3} 
Amount: 1000 
```

In this example, the inputs are a deposit {0 0 0 2} and a fee from the previous transaction {4 65535 0 0}. Neither input needs a confirmation signature.
```
plasmacli send \
--address 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986,0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986 \
--position 0.0.0.2::4.65535.0.0 \
--to 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2,0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986 \
--amounts 1000,100,0
Password to sign with '0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986':
Password to sign with '0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986':
Committed at block 9. Hash BBEDC0E5F6412E1D643CF26BE0C2A6312772824E3F943E3F8DCF953C3ABA90B1
```

## Spending UTXOS ##

The info command is useful in displaying all the information you need to spend utxos that need confirmation signatures

```
plasmacli info 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986

Position: {4 0 0 0} 
Amount: 500 
Denomination: Ether 
First Input Address: 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 

Position: {9 0 1 0} 
Amount: 100 
Denomination: Ether 
First Input Address: 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986 
Second Input Address: 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986 

Position: {0 0 0 3} 
Amount: 1000 
Denomination: Ether 
```

```
plasmacli info 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2

Position: {4 0 1 0} 
Amount: 400 
Denomination: Ether 
First Input Address: 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 

Position: {9 0 0 0} 
Amount: 1000 
Denomination: Ether 
First Input Address: 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986 
Second Input Address: 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986 
```

Whenever a transaction occurs, confirmation signatures must be generated as a way the sender can verify that the transaction was included in a block.

In order for `0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986` to spend {4 0 0 0}, a confirmation signature must be recieved from `0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2` 

In order for `0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2` to spend {9 0 0 0}, two confirmation signatures must be recieved. They both will be signed by `0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986`


In the following example, the inputs are {4 0 0 0} and {9 0 0 0}. We will need a total of 3 confirmation signatures to spend these inputs

We can use the sign command to generate confirmation signatures

```
plasmacli sign 4.0.0.0 \
--from 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 \
--owner 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986
Password to sign with '0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2':

Confirmation Signature for utxo with
position: {4 0 0 0} 
amount: 500
signature: 2ee2d01080111bffa7ba2cf6f997883b8c915a2b02d3ab3e7b1bd07959a46de51ac4172fe6d28ef8c76db79abfc4fb5469dba582429908a498f02aa874109cf400
UTXO had 1 inputs

```

--from: the address that of the input in the transaction that created this utxo. The same address as "First Input Address"

--owner: the address that owns the utxo


```
plasmacli sign 9.0.0.0 \
--owner 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 \
--from 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986
Password to sign with '0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986':

Confirmation Signature for utxo with
position: {9 0 0 0} 
amount: 1000
signature: 5aa6b5daeb9e9f848ec0fe7d2a9d220dc008d1a5ac85ca296b846800e42f9f4141e8360d6ce14ce59c7728f09b5c97bed3af5fba352b0ee7aec041ff7ae5804501
UTXO had 2 inputs
```

We use the first signature {4 0 0 0} for Input0ConfirmSigs and the second signature {9 0 0 0} for Input1ConfirmSigs

```
plasmacli send \
--address 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986,0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 \
--position 4.0.0.0::9.0.0.0 \
--amounts 700,800,0 \
--to 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986,0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 \
--Input0ConfirmSigs 2ee2d01080111bffa7ba2cf6f997883b8c915a2b02d3ab3e7b1bd07959a46de51ac4172fe6d28ef8c76db79abfc4fb5469dba582429908a498f02aa874109cf400 \
--Input1ConfirmSigs 5aa6b5daeb9e9f848ec0fe7d2a9d220dc008d1a5ac85ca296b846800e42f9f4141e8360d6ce14ce59c7728f09b5c97bed3af5fba352b0ee7aec041ff7ae5804501,5aa6b5daeb9e9f848ec0fe7d2a9d220dc008d1a5ac85ca296b846800e42f9f4141e8360d6ce14ce59c7728f09b5c97bed3af5fba352b0ee7aec041ff7ae5804501
Password to sign with '0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986':
Password to sign with '0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2':
Committed at block 15. Hash 04AAD098969EDFFC24C73FB21F84FFEFC072D2D4F8C15CADE6DED5D8832996EE
```

Since {9 0 0 0} had two input addresses, two signatures were provided (in this example the signatures were the same because the input addresses were the same)


#### Confirmation Signatures from different addresses ####

```
plasmacli info 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986

Position: {9 0 1 0} 
Amount: 100
Denomination: Ether 
First Input Address: 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986 
Second Input Address: 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986 

Position: {15 0 0 0} 
Amount: 700 
Denomination: Ether 
First Input Address: 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986 
Second Input Address: 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 

Position: {0 0 0 3} 
Amount: 1000 
Denomination: Ether 
```

```
plasmacli info 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2

Position: {4 0 1 0} 
Amount: 400 
Denomination: Ether 
First Input Address: 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 

Position: {15 0 1 0} 
Amount: 800 
Denomination: Ether 
First Input Address: 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986 
Second Input Address: 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 
```

In order to spend {15 0 0 0}, a confirmation signature from `0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986` and `0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2` is needed.

In order to spend {15 0 1 0}, a confirmation signature from `0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986` and `0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2` is needed.


```
plasmacli sign 15.0.0.0
--from 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986
--owner 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986
Password to sign with '0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986':

Confirmation Signature for utxo with
position: {15 0 0 0} 
amount: 700
signature: 55ad2e6b0c4f0a4477c269fee4090b96c35085aa1e1025ff6c9d2e286f0b4c89656f6029eb15ddfd99c0d1f700f27dacd9a758f51fad6af3c0ba7659f550c6ba00
UTXO had 2 inputs
```

```
plasmacli sign 15.0.0.0 --from 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 --owner 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986
Password to sign with '0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2':

Confirmation Signature for utxo with
position: {15 0 0 0} 
amount: 700
signature: 97d194c038a8d4afcad3eebd93a328b1dcb685df4f5c29827fbf0ea0b574166b280cafa937bb6e7aef4a047e4c152d4738f5064990d5add0552ff82f474291d800
UTXO had 2 inputs
```

```
plasmacli sign 15.0.1.0 --from 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986 --owner 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2
Password to sign with '0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986':

Confirmation Signature for utxo with
position: {15 0 1 0} 
amount: 800
signature: 55ad2e6b0c4f0a4477c269fee4090b96c35085aa1e1025ff6c9d2e286f0b4c89656f6029eb15ddfd99c0d1f700f27dacd9a758f51fad6af3c0ba7659f550c6ba00
UTXO had 2 inputs
```

```
plasmacli sign 15.0.1.0 --from 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 --owner 0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2
Password to sign with '0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2':

Confirmation Signature for utxo with
position: {15 0 1 0} 
amount: 800
signature: 97d194c038a8d4afcad3eebd93a328b1dcb685df4f5c29827fbf0ea0b574166b280cafa937bb6e7aef4a047e4c152d4738f5064990d5add0552ff82f474291d800
UTXO had 2 inputs
```

Notice how the signatures with the same --from address have the same signautres. This is because a confirmation signature is signed by the --from address over `Hash(Hash(txbytes) + root_hash)`. Therefore inputs in a transaction with the same from address will have the same confirmation signature even if they have different oindex values

```
plasmacli send \
--address 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986,0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 \
--amounts 800,700,0 \
--Input0ConfirmSigs 55ad2e6b0c4f0a4477c269fee4090b96c35085aa1e1025ff6c9d2e286f0b4c89656f6029eb15ddfd99c0d1f700f27dacd9a758f51fad6af3c0ba7659f550c6ba00,97d194c038a8d4afcad3eebd93a328b1dcb685df4f5c29827fbf0ea0b574166b280cafa937bb6e7aef4a047e4c152d4738f5064990d5add0552ff82f474291d800 \
--Input1ConfirmSigs 55ad2e6b0c4f0a4477c269fee4090b96c35085aa1e1025ff6c9d2e286f0b4c89656f6029eb15ddfd99c0d1f700f27dacd9a758f51fad6af3c0ba7659f550c6ba00,97d194c038a8d4afcad3eebd93a328b1dcb685df4f5c29827fbf0ea0b574166b280cafa937bb6e7aef4a047e4c152d4738f5064990d5add0552ff82f474291d800 \
--to 0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986,0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2 \
--position 15.0.0.0::15.0.1.0
Password to sign with '0xeA6eD4bB7CbA09c391C11a15D5472e806Caa3986':
Password to sign with '0xb79A48171a4BAF707CFcEF8A919F25aDbE7108a2':
Committed at block 78. Hash 405F5A168AA26DC63B3C243D02DF48441D32A9FD760EEAEAD9F8E7763ABB178B
```
