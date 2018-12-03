# Rootchain Documentation
```solidity
function submitBlock(bytes32 root)
```
The validator submits the block header, `root` of each child chain block. More than one block can be submitted per call by appending the roots to one another.  

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
function startTransactionExit(uint256[3] txPos, bytes txBytes, bytes proof, bytes sigs, bytes confirmSignatures)
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
function challengeTransactionExit(uint256[3] txPos, uint256[2] newTxPos, bytes txBytes, bytes proof, bytes sigs, bytes confirmSignature)
```
`txPos` and `newTxPos` follow the convention - `[blockNumber, transcationIndex, outputIndex]`

A uxto that has starting an exit phase but was already spent on the child chain can be challenged using this function call. A successfull challenge awards the caller with the exit bond.
The `txPos` locates the malicious utxo and is used to calculate the priority. `newTxPos` locates the transaction that is the parent (offending transaction is an input into this tx).
The `proof`, `txBytes` and `sigs` is sufficient for a proof of inclusion in the child chain of the parent transaction. The `confirmSignature`, signed by the owner of the malicious transaction,
acknowledges the inclusion of it's parent in the plasma chain and allows anyone with this confirm signature to challenge a malicious exit of the child.

<br />

```solidity
function challengeDepositExit(uint256 nonce, uint256[3] newTxPos, bytes txBytes, bytes sigs, bytes proof, bytes confirmSignature)
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
function getChildBlock(uint256 blockNumber) returns (bytes32 header, uint256 created_at)
```
Getter for the block header and when the block was submitted

<br />

```solidity
function getTransactionExit(uint256 priority) returns (address owner, uint256 amount, uint256[3] utxoPos, uint256 created_at, uint8 state)
```
Getter for all information about an exit

<br />

```solidity
function getDeposit(uint256 nonce) returns (address owner, uint256 amount, uint256 created_at)
```
Getter for all information about a deposit

<br />

```solidity
function childChainBalance() returns (uint256 funds)
```
Query the total funds of the plasma chain
