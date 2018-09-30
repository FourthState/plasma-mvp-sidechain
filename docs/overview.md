# Plasma Sidechain Overview

This sidechain is built on top of the Cosmos SDK and uses Tendermint for consensus. Therefore it is important to note the flow of a transaction as it would occur on any chain built off the SDK and not specifically on this one. 

The cosmos sdk docs can be found [here](https://github.com/cosmos/cosmos-sdk/tree/master/docs)

## Tx and Msg
There are two main types in the SDK, a Tx and a Msg. A message (Msg) can be thought of in the traditional distributed systems sense of being data passed around between nodes. A transaction (Tx) contains a message as well as any other authentication data necessary to validate a message. Both Tx and Msg are interfaces such that there can be many different types of messages and transactions. 

In our implementation, we have a SpendMsg and a BaseTx. The SpendMsg contains all the information necessary to Spend 1 or 2 UTXOs, while BaseTx contains a SpendMsg and signatures of the RLP encoded SpendMsg.
 
## App
Within the app package, there is a struct called ChildChain. This represents the blockchain itself. A validator must first initialize an instance of the ChildChain and then begin passing transaction through it to get it to update state. 

When passing around transactions and messages within the application, the SDK utilizes Context which provides all the necessary information to access context on the state of the blockchain when the transaction is being processed. There are two types of processes that can be run with a Context, checkTx and deliverTx. 

CheckTx is executed on a transaction by a validator when it is deciding whether to include a transaction into a block (checkTx does not update state). DeliverTx is executed on a transaction after a transaction has been included into a block and therefore updates the state of our blockchain. 

## Processing a transaction 
When the tx bytes are sent to a validator, the validator executes the function ValidateBasic() which belongs to the Msg interface. ValidateBasic does a simple check to ensure that the message created is well formed. For example, SpendMsg will check that the two inputs provided don't equal each other (double spend) and that fields such as Oindex, which require a certain range of numbers, have been filled in appropriately. 

If ValidateBasic does not return any errors, then the transaction is added to the mempool. CheckTx will then be executed before the tx is included into a block. CheckTx uses the ante handler to verify the authentication data in a transaction, verify the message against the state of the blockchain, and check that the fee provided is sufficient. 

Our ante handler will check that the address that created the signatures in BaseTx match those that own the inputs of the transaction as well as check that the inputs being spent exist. The ante handler also checks that the inputs = outputs + fee. 

Once a transaction has been included into a block, DeliverTx will be executed which will do the same functionality as CheckTx as well as route the Msg to a handler. 

In our implementation, our handler can only handle a SpendMsg. Handling of a SpendMsg includes deleting the inputs of the transaction from our utxo database and creating new utxos corresponding to the addresses and denominations provided in SpendMsg. handleSpendMsg() will utilize the UTXOKeeper which is used to restrict access to our UTXOMapper (utxo database). The handler will also increment the transaction index number so the newly created utxo's will map from the correct position.

The UTXOMapper is our utxo database. Our mapper uses keys in the form: < encoded address > + < encoded position > . It maps to the encoded utxo and uses go-amino for its encoding. The < encoded position > at the beginning of the key is used for prefix iteration which will return all the utxo's owned by a specified address. 

## Types

**SpendMsg**

A SpendMsg contains the position of the input utxos (block number, transaction index, output index, deposit number), confirm signatures signed by the owners of the parent inputs (the inputs to the SpendMsg inputs), the addresses of the input owners, the addresses of the outputs, the denomination of each output, and the fee amount. 

GetSignBytes() returns the rlp encoded bytes of the SpendMsg

GetSigners() returns the input owner addresses as sdk.Address's 

**Position**

A Position contains the block number, transaction index, output index, and deposit number.

GetSignBytes() returns the rlp encoded bytes of the Position

**UTXO**

A UTXO contains the address of the owner of the utxo, the denomination of the utxo, the position of the utxo, and the input addresses that were used to create the utxo. 

## Utils

ZeroAddress(common.Address) returns true if the address provided is the zero address (0x00000...) and false otherwise

ValidAddress(common.Address) returns true if the address provided is properly formatted and false otherwise.

PrivKeyToAddress(*ecdsa.PrivateKey) returns the common.Address corresponding to the private key provided.

GenerateAddress() generates a random common.Address

