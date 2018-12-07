pragma solidity ^0.4.24;

// external modules
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/math/Math.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/ECRecovery.sol";
import "solidity-rlp/contracts/RLPReader.sol";

import "../libraries/Validator.sol";
import "../libraries/PriorityQueue.sol";

contract RootChain is Ownable {
    using PriorityQueue for uint256[];
    using RLPReader for bytes;
    using RLPReader for RLPReader.RLPItem;
    using SafeMath for uint256;
    using Validator for bytes32;

    /*
     * Events
     */

    event AddedToBalances(address owner, uint256 amount);
    event BlockSubmitted(bytes32 root, uint256 blockNumber, uint256 numTxns, uint256 feeAmount);
    event Deposit(address depositor, uint256 amount, uint256 depositNonce);

    event StartedTransactionExit(uint[3] position, address owner, uint256 amount, bytes confirmSignatures);
    event StartedDepositExit(uint nonce, address owner, uint256 amount);

    event ChallengedExit(uint[4] position, address owner, uint256 amount);
    event FinalizedExit(uint[4] position, address owner, uint256 amount);

    /*
     *  Storage
     */

    // child chain
    uint256 public lastCommittedBlock;
    uint256 public depositNonce;
    mapping(uint256 => childBlock) public childChain;
    mapping(uint256 => depositStruct) public deposits;
    struct childBlock {
        bytes32 root;
        uint256 numTxns;
        uint256 feeAmount;
        uint256 createdAt;
    }
    struct depositStruct {
        address owner;
        uint256 amount;
        uint256 createdAt;
    }

    // exits
    uint256 minExitBond;
    uint256[] txExitQueue;
    uint256[] depositExitQueue;
    mapping(uint256 => exit) public txExits;
    mapping(uint256 => exit) public depositExits;
    enum ExitState { NonExistent, Pending, Challenged, Finalized }
    struct exit {
        uint256 amount;
        uint256 createdAt;
        address owner;
        uint256[4] position; // (blkNum, txIndex, outputIndex, depositNonce)
        ExitState state; // default value is `NonExistent`
    }

    // funds
    mapping(address => uint256) public balances;
    uint256 public totalWithdrawBalance;

    // constants
    uint256 public constant txIndexFactor = 10;
    uint256 public constant blockIndexFactor = 1000000;

    constructor() public
    {
        lastCommittedBlock = 0;
        depositNonce = 1;
        minExitBond = 10000;
    }

    // @param blocks       32 byte merkle roots appended in ascending order
    // @param txnsPerBlock number of transactions per block
    // @param feesPerBlock amount of fees the validator has collected per block
    // @param blockNum     the block number of the first header
    function submitBlock(bytes blocks, uint256[] txnsPerBlock, uint256[] feesPerBlock, uint256 blockNum)
        public
        onlyOwner
    {
        require(blockNum == lastCommittedBlock + 1, "inconsistent block number ordering");
        require(blocks.length > 0 && blocks.length % 32 == 0, "block roots must be of size 32 bytes");
        require(blocks.length / 32 == txnsPerBlock.length, "blocks and txnsPerBlock lengths are inconsistent");
        require(blocks.length / 32 == feesPerBlock.length, "blocks and feesPerBlock lengths are inconsistent");

        uint memPtr;
        assembly  {
            memPtr := add(blocks, 0x20)
        }

        uint maxNumTxPerBlock = 2**16 - 1;
        bytes32 root;
        for (uint i = 0; i < txnsPerBlock.length; i++) {
            require(txnsPerBlock[i] <= maxNumTxPerBlock, "number of transactions per block exceeds limit");

            assembly {
                root := mload(add(memPtr, mul(i, 32)))
            }

            childChain[blockNum + i] = childBlock(root, txnsPerBlock[i], feesPerBlock[i], block.timestamp);
            emit BlockSubmitted(root, blockNum + i, txnsPerBlock[i], feesPerBlock[i]);
        }

        lastCommittedBlock = lastCommittedBlock.add(txnsPerBlock.length);
   }

    // @param owner owner of this deposit
    function deposit(address owner)
        public
        payable
    {
        deposits[depositNonce] = depositStruct(owner, msg.value, block.timestamp);
        emit Deposit(owner, msg.value, depositNonce);

        depositNonce = depositNonce.add(1);
    }

    // @param depositNonce the nonce of the specific deposit
    function startDepositExit(uint256 nonce)
        public
        payable
    {
        require(deposits[nonce].owner == msg.sender, "mismatch in owner");
        require(depositExits[nonce].state == ExitState.NonExistent, "exit for this deposit already exists");
        require(msg.value >= minExitBond, "insufficient exit bond");
        if (msg.value > minExitBond) {
            uint256 excess = msg.value - minExitBond;
            balances[msg.sender] = balances[msg.sender].add(excess);
            totalWithdrawBalance = totalWithdrawBalance.add(excess);
        }

        uint amount = deposits[nonce].amount;
        address owner = deposits[nonce].owner;
        depositExitQueue.insert(nonce);
        depositExits[nonce] = exit({
            owner: owner,
            amount: amount,
            createdAt: block.timestamp,
            position: [0,0,0,nonce],
            state: ExitState.Pending
        });

        emit StartedDepositExit(nonce, owner, amount);
    }

    // Transaction encoding:
    // [[Blknum1, TxIndex1, Oindex1, DepositNonce1, Owner1, Input1ConfirmSig,
    //   Blknum2, TxIndex2, Oindex2, DepositNonce2, Owner2, Input2ConfirmSig,
    //   NewOwner, Denom1, NewOwner, Denom2, Fee],
    //  [Signature1, Signature2]]
    //
    // @param txBytes rlp encoded transaction
    // @notice this function will revert if the txBytes are malformed
    function decodeTransaction(bytes txBytes)
        internal
        pure
        returns (RLPReader.RLPItem[] memory txList, RLPReader.RLPItem[] memory sigList, bytes32 txHash)
    {
        RLPReader.RLPItem[] memory spendMsg = txBytes.toRlpItem().toList();
        require(spendMsg.length == 2, "incorrect encoding of the transcation");

        txList = spendMsg[0].toList();
        require(txList.length == 17, "incorrect number of items in the transaction list");

        sigList = spendMsg[1].toList();
        require(sigList.length == 2, "two signatures must be present");

        // bytes the signatures are over
        txHash = keccak256(spendMsg[0].toRlpBytes());
    }

    // @param txPos             location of the transaction [blkNum, txIndex, outputIndex]
    // @param txBytes           raw transaction bytes
    // @param proof             merkle proof of inclusion in the child chain
    // @param confirmSignatures confirm signatures sent by the owners of the inputs acknowledging the spend.
    // @notice `confirmSignatures` and `ConfirmSig0`/`ConfirmSig1` are unrelated to each other.
    // @notice `confirmSignatures` is either 65 or 130 bytes in length dependent on if input2 is used.
    function startTransactionExit(uint256[3] txPos, bytes txBytes, bytes proof, bytes confirmSignatures)
        public
        payable
    {
        bytes32 txHash;
        RLPReader.RLPItem[] memory txList;
        RLPReader.RLPItem[] memory sigList;
        (txList, sigList, txHash) = decodeTransaction(txBytes);

        require(msg.sender == txList[12 + 2 * txPos[2]].toAddress(), "mismatch in utxo owner");
        require(msg.value >= minExitBond, "insufficient exit bond");
        if (msg.value > minExitBond) {
            uint256 excess = msg.value.sub(minExitBond);
            balances[msg.sender] = balances[msg.sender].add(excess);
            totalWithdrawBalance = totalWithdrawBalance.add(excess);
        }

        childBlock storage blk = childChain[txPos[0]];

        // check signatures
        bytes32 merkleHash = sha256(txBytes);
        require(txHash.checkSigs(sha256(abi.encodePacked(merkleHash, blk.root)), // confirmation hash -- sha256(merkleHash, root)
                         // we always assume the first input is always present in a transaction. The second input is optional
                         txList[6].toUint() > 0 || txList[9].toUint() > 0, // existence of input1. Either a deposit or utxo
                         sigList[0].toBytes(), sigList[1].toBytes(), confirmSignatures), "signature mismatch");

        // check proof
        require(merkleHash.checkMembership(txPos[1], blk.root, proof, blk.numTxns), "invalid merkle proof");

        // check that the UTXO's two direct inputs have not been previously exited
        require(validateTransactionExitInputs(txList), "an input is pending an exit or has been finalized");

        uint256 position = blockIndexFactor*txPos[0] + txIndexFactor*txPos[1] + txPos[2];
        require(txExits[position].state == ExitState.NonExistent, "this exit has already been started, challenged, or finalized");

        // calculate the priority of the transaction taking into account the withdrawal delay attack
        // withdrawal delay attack: https://github.com/FourthState/plasma-mvp-rootchain/issues/42
        txExitQueue.insert(Math.max256(blk.createdAt + 1 weeks, block.timestamp) << 128 | position);
        txExits[position] = exit({
            owner: txList[12 + 2 * txPos[2]].toAddress(),
            amount: txList[13 + 2 * txPos[2]].toUint(),
            createdAt: block.timestamp,
            position: [txPos[0], txPos[1], txPos[2], 0],
            state: ExitState.Pending
        });

        emit StartedTransactionExit(txPos, msg.sender, txList[13 + 2 * txPos[2]].toUint(), confirmSignatures);
    }

    // For any attempted exit of an UTXO, validate that the UTXO's two inputs have not
    // been previously exited or are currently pending an exit.
    function validateTransactionExitInputs(RLPReader.RLPItem[] memory txList)
        private
        view
        returns (bool)
    {
        for (uint256 i = 0; i < 2; i++) {
            ExitState state;
            uint depositNonce_ = txList[6*i + 3].toUint();
            if (depositNonce_ == 0) {
                uint256 blkNum = txList[6*i + 0].toUint();
                uint256 inputIndex = txList[6*i + 1].toUint();
                uint256 outputIndex = txList[6*i + 2].toUint();
                uint256 position = blockIndexFactor*blkNum + txIndexFactor*inputIndex + outputIndex;
                state = txExits[position].state;
            } else
                state = depositExits[depositNonce_].state;

            if (state != ExitState.NonExistent && state != ExitState.Challenged)
                return false;
        }

        return true;
    }

    // Validator of any block can call this function to exit the fees collected
    // for that particular block. The fee exit is added to exit queue with the lowest priority for that block.
    // In case of the fee UTXO already spent, anyone can challenge the fee exit by providing
    // the spend of the fee UTXO.
    // @param blockNumber the block for which the validator wants to exit fees
    function startFeeExit(uint256 blockNumber)
        public
        payable
        onlyOwner
    {
        // specified blockNumber must exist in child chain
        require(childChain[blockNumber].root != bytes32(0), "specified block does not exist in child chain.");

        require(msg.value >= minExitBond, "insufficient exit bond");
        if (msg.value > minExitBond) {
            uint256 excess = msg.value.sub(minExitBond);
            balances[msg.sender] = balances[msg.sender].add(excess);
            totalWithdrawBalance = totalWithdrawBalance.add(excess);
        }

        // a fee UTXO has explicitly defined position [blockNumber, 2**16 - 1, 0]
        uint256 txIndex = 2**16 - 1;
        uint256 position = blockIndexFactor*blockNumber + txIndexFactor*txIndex + 0;
        require(txExits[position].state == ExitState.NonExistent, "this exit has already been started, challenged, or finalized");

        txExitQueue.insert(Math.max256(childChain[blockNumber].createdAt + 1 weeks, block.timestamp) << 128 | position);
        uint256 feeAmount = childChain[blockNumber].feeAmount;
        txExits[position] = exit({
            owner: msg.sender,
            amount: feeAmount,
            createdAt: block.timestamp,
            position: [blockNumber, txIndex, 0, 0],
            state: ExitState.Pending
        });

        // pass in empty bytes for confirmSignatures for StartedTransactionExit event.
        emit StartedTransactionExit([blockNumber, txIndex, 0], msg.sender, feeAmount, "");
    }

    // @param depositNonce     the nonce of the deposit trying to exit
    // @param newTxPos         position of the transaction with this deposit as an input [blkNum, txIndex, outputIndex]
    // @param txBytes          bytes of this transcation
    // @param proof            merkle proof of inclusion
    // @param confirmSignature signature used to invalidate the invalid exit. Signature is over (merkleHash, block header)
    function challengeDepositExit(uint256 nonce, uint256[3] newTxPos, bytes txBytes, bytes proof, bytes confirmSignature)
        public
    {
        RLPReader.RLPItem[] memory txList;
        RLPReader.RLPItem[] memory sigList;
        (txList, sigList, ) = decodeTransaction(txBytes);

        // ensure that the txBytes is a direct spend of the deposit
        require(nonce == txList[3].toUint() || nonce == txList[9].toUint(), "challenging transaction is not a direct spend");

        exit memory exit_ = depositExits[nonce];
        require(exit_.state == ExitState.Pending, "no pending exit to challenge");

        // check for inclusion in the side chain
        childBlock storage blk = childChain[newTxPos[0]];

        bytes32 merkleHash = sha256(txBytes);
        bytes32 confirmationHash = sha256(abi.encodePacked(merkleHash, blk.root));
        require(exit_.owner == confirmationHash.recover(confirmSignature), "mismatch in exit owner and confirm signature");
        require(merkleHash.checkMembership(newTxPos[1], blk.root, proof, blk.numTxns), "incorrect merkle proof");

        // exit successfully challenged
        balances[msg.sender] = balances[msg.sender].add(minExitBond);
        totalWithdrawBalance = totalWithdrawBalance.add(minExitBond);

        depositExits[nonce].state = ExitState.Challenged;
        emit ChallengedExit(exit_.position, exit_.owner, exit_.amount);
    }

    // @param exitingTxPos     position of the invalid exiting transaction [blkNum, txIndex, outputIndex]
    // @param challengingTxPos position of the challenging transaction [blkNum, txIndex, outputIndex]
    // @param txBytes          raw transaction bytes of the challenging transaction
    // @param proof            proof of inclusion for this merkle hash
    // @param confirmSignature signature used to invalidate the invalid exit. Signature is over (merkleHash, block header)
    function challengeTransactionExit(uint256[3] exitingTxPos, uint256[3] challengingTxPos, bytes txBytes, bytes proof, bytes confirmSignature)
        public
    {
        RLPReader.RLPItem[] memory txList;
        RLPReader.RLPItem[] memory sigList;
        (txList, sigList, ) = decodeTransaction(txBytes);

        // must be a direct spend
        require(ensureMatchingInputs(exitingTxPos, txList), "challenging transaction is not a direct spend");

        // transaction to be challenged should have a pending exit
        uint256 position = blockIndexFactor*exitingTxPos[0] + txIndexFactor*exitingTxPos[1] + exitingTxPos[2];
        exit memory exit_ = txExits[position];
        require(exit_.state == ExitState.Pending, "no pending exit to challenge");

        // confirm challenging transcation's inclusion and confirm signature
        childBlock storage blk = childChain[challengingTxPos[0]];

        bytes32 merkleHash = sha256(txBytes);
        bytes32 confirmationHash = sha256(abi.encodePacked(merkleHash, blk.root));
        require(exit_.owner == confirmationHash.recover(confirmSignature), "mismatch in exit owner and confirm signature");
        require(merkleHash.checkMembership(challengingTxPos[1], blk.root, proof, blk.numTxns), "incorrect merkle proof");

        // exit successfully challenged. Award the sender with the bond
        balances[msg.sender] = balances[msg.sender].add(minExitBond);
        totalWithdrawBalance = totalWithdrawBalance.add(minExitBond);
        emit AddedToBalances(msg.sender, minExitBond);

        // reflect challenged state
        txExits[position].state = ExitState.Challenged;
        emit ChallengedExit(exit_.position, exit_.owner, exit_.amount);
    }

    // When challenging an exiting transcation located at `exitingTxPos`, we must make sure that the challenging
    // transcation posted is either a direct spend of the exit if the confirm signature was not included in the txBytes of
    // the exiting transaction
    function ensureMatchingInputs(uint256[3] exitingTxPos, RLPReader.RLPItem[] memory challengingTxList)
        private
        pure
        returns (bool)
    {
        // indicator for which input to check int the challenging transaction
        uint i = exitingTxPos[0] == challengingTxList[0].toUint() ? 0 : 1;

        if (exitingTxPos[0] == challengingTxList[0 + 6*i].toUint()
            && exitingTxPos[1] == challengingTxList[1 + 6*i].toUint()
            && exitingTxPos[2] == challengingTxList[2 + 6*i].toUint())
            return true;

        return false;
    }

    function finalizeDepositExits() public { finalize(depositExitQueue, true); }
    function finalizeTransactionExits() public { finalize(txExitQueue, false); }

    // Finalizes exits by iterating through either the depositExitQueue or txExitQueue.
    // Users can determine the number of exits they're willing to process by varying
    // the amount of gas allow finalize*Exits() to process.
    // Each transaction takes < 80000 gas to process.
    function finalize(uint256[] storage queue, bool isDeposits)
        private
    {
        // getMin will fail if nothing is in the queue
        if (queue.currentSize() == 0) {
            return;
        }

        // retrieve the lowest priority and the appropriate exit struct
        uint256 priority = queue.getMin();
        exit memory currentExit;
        uint256 position;
        if (isDeposits) {
            currentExit = depositExits[priority];
        } else {
            // retrieve the right 128 bits from the priority to obtain the position
            assembly {
   			    position := and(priority, div(not(0x0), exp(256, 16)))
		    }
            currentExit = txExits[position];
        }

        /*
        * Conditions:
        *   1. Exits exist
        *   2. Exits must be a week old
        *   3. Funds must exist for the exit to withdraw
        */
        uint256 amountToAdd;
        while (queue.currentSize() > 0 &&
               (block.timestamp - currentExit.createdAt) > 1 weeks &&
               currentExit.amount.add(minExitBond) <= address(this).balance - totalWithdrawBalance &&
               gasleft() > 80000) {

            // skip currentExit if it is not in 'started/pending' state.
            if (currentExit.state != ExitState.Pending) {
                queue.delMin();
            } else {
                amountToAdd = currentExit.amount.add(minExitBond);
                balances[currentExit.owner] = balances[currentExit.owner].add(amountToAdd);
                totalWithdrawBalance = totalWithdrawBalance.add(amountToAdd);

                if (isDeposits)
                    depositExits[priority].state = ExitState.Finalized;
                else
                    txExits[position].state = ExitState.Finalized;

                emit FinalizedExit(currentExit.position, currentExit.owner, amountToAdd);
                emit AddedToBalances(currentExit.owner, amountToAdd);

                // move onto the next oldest exit
                queue.delMin();
            }

            if (queue.currentSize() == 0) {
                return;
            }

            // move onto the next oldest exit
            priority = queue.getMin();
            if (isDeposits) {
                currentExit = depositExits[priority];
            } else {
                // retrieve the right 128 bits from the priority to obtain the position
                assembly {
   			        position := and(priority, div(not(0x0), exp(256, 16)))
		        }
                currentExit = txExits[position];
            }
        }
    }

    function withdraw()
        public
        returns (uint256)
    {
        if (balances[msg.sender] == 0) {
            return 0;
        }

        uint256 transferAmount = balances[msg.sender];
        delete balances[msg.sender];
        totalWithdrawBalance = totalWithdrawBalance.sub(transferAmount);

        // will revert the above deletion if fails
        msg.sender.transfer(transferAmount);
        return transferAmount;
    }

    /*
    * Getters
    */

    function childChainBalance()
        public
        view
        returns (uint)
    {
        // takes into accounts the failed withdrawals
        return address(this).balance - totalWithdrawBalance;
    }

    function balanceOf(address _address)
        public
        view
        returns (uint256)
    {
        return balances[_address];
    }
}
