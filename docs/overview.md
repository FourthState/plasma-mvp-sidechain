# Plasma Sidechain Overview

This sidechain is built on top of the Cosmos SDK and uses Tendermint for consensus. Therefore it is important to note the flow of a transaction as it would occur on any chain built off the SDK and not specifically on this one. 

[Cosmos SDK documentation](https://cosmos.network/docs/)

## Tx and Msg
There are two main types in the SDK, a Tx and a Msg. 
A message (Msg) can be thought of in the traditional distributed systems sense of being data passed around between nodes.
A transaction (Tx) contains a message as well as any other authentication data necessary to validate a message. 
Both Tx and Msg are interfaces such that there can be many different types of messages and transactions. 

In our implementation, we have a SpendMsg, IncludeDepositMsg and a Transaction. 
The SpendMsg contains all the information necessary to Spend 1 or 2 UTXOs, while Transaction contains a SpendMsg and signatures of the RLP encoded SpendMsg.
We also have IncludeDepositMsg, which signifies to the sidechain that a deposit has occured on the rootchain and should be included into the utxo store. 
 
## Server/App
In the server directory we have two packages, app and plasmad.
Within the app directory, there exists a struct called PlasmaMVPChain in app.go.
This represents the blockchain itself. 
A validator must first initialize an instance of the PlasmaMVP and then begin passing transaction through it to get it to update state. 
The plasmad directory represents this initialization via a command line program. 
Any user who wants to run a full node must install this package.

When passing around transactions and messages within the application, the SDK utilizes Context which provides all the necessary information to access context on the state of the blockchain when the transaction is being processed.
There are two types of processes that can be run with a Context, checkTx and deliverTx. 

CheckTx is executed on a transaction by a validator when it is deciding whether to include a transaction into a block (checkTx does not update state). 
DeliverTx is executed on a transaction after a transaction has been included into a block and therefore updates the state of our blockchain. 

## Processing a transaction 
When the tx bytes are sent to a validator, the validator executes the function ValidateBasic() which belongs to the Msg interface. 
ValidateBasic does a simple check to ensure that the message created is well formed.
For example, SpendMsg will check that the two inputs provided don't equal each other (double spend) and that fields such as Oindex, which require a certain range of numbers, have been filled in appropriately. 

If ValidateBasic does not return any errors, then the transaction is added to the mempool.
CheckTx will then be executed before the tx is included into a block. 
CheckTx uses the ante handler to verify the authentication data in a transaction, verify the message against the state of the blockchain, and check that the fee provided is sufficient. 

Our ante handler will check that the address that created the signatures in Transaction match those that own the inputs of the transaction as well as check that the inputs being spent exist.
The ante handler also checks that the inputs = outputs + fee and that the utxo has not been exitted on the rootchain.
InputConfirmationSigantures will also be checked for each input that is not a deposit nor a fee utxo.  

Once a transaction has been included into a block, DeliverTx will be executed which will do the same functionality as CheckTx as well as route the Msg to a handler. 

The UTXOMapper is our utxo database. 
Our mapper uses keys in the form: < encoded address > + < encoded position > . 
It maps to the encoded utxo and uses go-amino for its encoding. 
The < encoded position > at the beginning of the key is used for prefix iteration which will return all the utxo's owned by a specified address. 

## Store

There are two key value stores in our implementation, utxo store and plasma store. 
The UTXO store is our utxo database and contains all information related to utxo's.
The plasma store contains plasma specific information including plasma block headers and confirmation signatures.

## Client

The client directory contains command line client in plasmacli/ and key value stores for confirmation signatures and private keys in store/.
Plasmacli uses the eth command to do rootchain related actions and the query command to query the current state of the sidechain. 
The keys command manages private key storage. 
Sign can be used to generate and store confirmation signatures.
Confirmation signatures are stored using a goleveldb in .plasmacli/data/signatures.ldb.
We use a mapping from string to address to allow users assign human readable names to their keys. 
This data is stored in .plasmacli/data/accounts.ldb


## Plasma

**Block**

A block contains root hash, number of transactions, and fee associated with a plasma block. 
A plasma block is different from a tendermint block. A tendermint block will contains all the messages of the sidechain while a plasma block contains only valid utxo transactions.

**Deposit**

A deposit contains the owner address, the amount of the deposit, and the ethereum block number it was created in.

**Input**

A input contains the Position, transaction signature, and confirmation signature for the input. 

**Output**

A output contains the owner address and amount of the output. 

**Position**

A Position contains the block number, transaction index, output index, and deposit nonce.

`FromPositionString(string)` takes in a string in the format "(blknum.txindex.oindex.depositnonce)" and returns a Position with the passed in values. 

**Transaction**


## Utils

IsZeroAddress(common.Address) returns true if the address provided is the zero address (0x00000...) and false otherwise

RemoveHexPrefix(string) returns the passed in string without a "0x" prefix if it exists. 

ToEthSignedMessage([] byte) returns the hash of the passed in message with "\x19EthereumSignedMessage\n32" prefixed. 
This is [standard procedure](https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign) for generating ethereum transactions. 
