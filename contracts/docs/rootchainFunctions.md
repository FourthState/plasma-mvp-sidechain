# Rootchain Documentation

The transcation bytes, `txBytes`, in the contract follow the convention:  
```
RLP_ENCODE ([ 
  [Blknum1, TxIndex1, Oindex1, DepositNonce1, Amount1, ConfirmSig1,

  Blknum2, TxIndex2, Oindex2, DepositNonce2, Amount2, ConfirmSig2,

  NewOwner, Denom1, NewOwner, Denom2, Fee],

  [Signature1, Signature2]
])
```
```solidity
function submitBlock(bytes32 blocks, uint256[] txnsPerBlock, uint256 blockNum)
```
The validator submits appended block headers in ascending order. Each block can be of variable block size(capped at 2^16 txns per block). The total number of transactions per block must be passed in through `txnsPerBlock`.
`blockNum` must be the intended block number of the first header in this call. Ordering is enforced on each call. `blockNum == lastCommittedBlock + 1`.

<br >

```solidity
function deposit(address owner)
```
Entry point into the child chain. The user has the option to create a spendable utxo owned by the address, `owner`. Once created,
the private keys of the `owner` address has complete control of the new utxo.

Deposits are not recorded in the child chain blocks and are entirely represented on the rootchain. Each deposit is identified with an incremental nonce.
Validators catch deposits through event handlers and maintain a collection of spendable deposits.
```solidity
mapping(uint256 => depositStruct) deposits; // The key is the incrementing nonce
struct depositStruct {
    address owner;
    uint256 amount;
    uint256 created_at;
}
```

<br />

```solidity
function startTransactionExit(uint256[3] txPos, bytes txBytes, bytes proof, bytes confirmSignatures)
```
`txPos` follows the convention - `[blockNumber, transactionIndex, outputIndex]`

Exit procedure for exiting a utxo on the child chain(not deposits). The `txPos` locates the transaction on the child chain. The leaf, hash(hash(`txBytes`), `sigs`) is checked against the block header using the `proof`.
The `confirmSignatures` represent the acknowledgement of the inclusion by both inputs. If only one input was used to create this transactions, only one confirm signature should be passed in for the corresponding
input. However, if there are two distinct inputs in the exiting transactions, both confirm signatures should be appended together in order for a total of 130 bytes.

A valid exit satisfies the following properties:
  - Exit has not previously been finalized or challenged
  - The creator of this exit posted a sufficient bond. Excess funds are refunded the the senders rootchain balance and are immediately withdrawable.
  - If present, the confirm signatures are correct and signed by the same address which signed the corresponding input signatures.

<br />

```solidity
function startDepositExit(uint256 nonce)
```
Exit procedure for deposits that have not been spent. Deposits are purely identified by their `nonce`. The caller's address must match the owner of the deposit.
A valid exit must satisfy the same constraints listed above for normal utxo exits except confirm signatures. Deposits exits are also collected into their own seperate queue from normal transcations.
This is because of the differing priority calculation. The priority of a deposit is purely it's nonce while the priority of a utxo is calculated from it's location in the child chain.

<br />

```solidity
function challengeTransactionExit(uint256[3] exitingTxPos, uint256[2] challengingTxPos, bytes txBytes, bytes proof, bytes confirmSignature)
```
`exitingTxPos` and `challengingTxPos` follow the convention - `[blockNumber, transcationIndex, outputIndex]`

A uxto that has starting an exit phase but was already spent on the child chain can be challenged using this function call. A successfull challenge awards the caller with the exit bond.
The `exitingTxPos` locates the malicious utxo and is used to calculate the priority. `challengingTxPos` locates the transaction that is the parent (offending transaction is an input into this tx).
The `proof`, `txBytes` and `sigs` is sufficient for a proof of inclusion in the child chain of the parent transaction. The `confirmSignature`, signed by the owner of the malicious transaction,
acknowledges the inclusion of it's parent in the plasma chain and allows anyone with this confirm signature to challenge a malicious exit of the child.

<br />

```solidity
function challengeDepositExit(uint256 nonce, uint256[3] newTxPos, bytes txBytes, bytes proof, bytes confirmSignature)
```
A deposit that has been spent in the child chain is challenged here. The `txBytes` of the the parent transaction must include the nonce as one if it's input. The `txBytes`, `sigs` and `proof` is
sufficient for a proof of inclusion. Similar to a normal challenge, the owner of the deposit must have also broadcasted a `confirmSignature` acknowledging the spend. A successfull challenge awards the
caller with the exit bond.

<br />

```solidity
function finalizeTransactionExits()
```
Process all "finalized" exits in the priority queue. "Finalized" exits are those that have been in the priority queue for at least one week and have not been proven to be malicious through a challenge.

<br />

```solidity
function finalizeDepositExits()
```
Process all "finalized" deposit exits in the priority queue. "Finalized" exits are those that have been in the priority queue for at least one week and have not been proven to be malicious through a challenge.

<br />

```solidity
function withdraw()
```
Sender withdraws all funds associated with their balance from the contract.

<br />

```solidity
function balanceOf(address _address) returns (uint256 amount)
```
Getter for the withdrawable balance of `_address`

<br />

```solidity
function childChainBalance() returns (uint256 funds)
```
Query the total funds of the plasma chain
