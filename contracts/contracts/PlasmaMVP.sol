pragma solidity ^0.5.0;

// external modules
import "solidity-rlp/contracts/RLPReader.sol";

// libraries
import "./libraries/BytesUtil.sol";
import "./libraries/SafeMath.sol";
import "./libraries/ECDSA.sol";
import "./libraries/TMSimpleMerkleTree.sol";
import "./libraries/MinPriorityQueue.sol";

contract PlasmaMVP {
    using MinPriorityQueue for uint256[];
    using BytesUtil for bytes;
    using RLPReader for bytes;
    using RLPReader for RLPReader.RLPItem;
    using SafeMath for uint256;
    using TMSimpleMerkleTree for bytes32;
    using ECDSA for bytes32;

    /*
     * Events
     */

    event ChangedOperator(address oldOperator, address newOperator);

    event AddedToBalances(address owner, uint256 amount);
    event BlockSubmitted(bytes32 header, uint256 blockNumber, uint256 numTxns, uint256 feeAmount);
    event Deposit(address depositor, uint256 amount, uint256 depositNonce, uint256 ethBlockNum);

    event StartedTransactionExit(uint256[3] position, address owner, uint256 amount, bytes confirmSignatures, uint256 committedFee);
    event StartedDepositExit(uint256 nonce, address owner, uint256 amount, uint256 committedFee);

    event ChallengedExit(uint256[4] position, address owner, uint256 amount);
    event FinalizedExit(uint256[4] position, address owner, uint256 amount);

    /*
     *  Storage
     */

    address public operator;

    uint256 public lastCommittedBlock;
    uint256 public depositNonce;
    mapping(uint256 => plasmaBlock) public plasmaChain;
    mapping(uint256 => depositStruct) public deposits;
    struct plasmaBlock{
        bytes32 header;
        uint256 numTxns;
        uint256 feeAmount;
        uint256 createdAt;
        uint256 ethBlockNum;
    }
    struct depositStruct {
        address owner;
        uint256 amount;
        uint256 createdAt;
        uint256 ethBlockNum;
    }

    // exits
    uint256 public minExitBond;
    uint256[] public txExitQueue;
    uint256[] public depositExitQueue;
    mapping(uint256 => exit) public txExits;
    mapping(uint256 => exit) public depositExits;
    enum ExitState { NonExistent, Pending, Challenged, Finalized }
    struct exit {
        uint256 amount;
        uint256 committedFee;
        uint256 createdAt;
        address owner;
        uint256[4] position; // (blkNum, txIndex, outputIndex, depositNonce)
        ExitState state; // default value is `NonExistent`
    }

    // funds
    mapping(address => uint256) public balances;
    uint256 public totalWithdrawBalance;

    // constants
    uint256 constant txIndexFactor = 10;
    uint256 constant blockIndexFactor = 1000000;
    uint256 constant lastBlockNum = 2**109;
    uint256 constant feeIndex = 2**16-1;

    /** Modifiers **/
    modifier isBonded()
    {
        require(msg.value >= minExitBond);
        if (msg.value > minExitBond) {
            uint256 excess = msg.value.sub(minExitBond);
            balances[msg.sender] = balances[msg.sender].add(excess);
            totalWithdrawBalance = totalWithdrawBalance.add(excess);
        }

        _;
    }

    modifier onlyOperator()
    {
        require(msg.sender == operator);
        _;
    }

    function changeOperator(address newOperator)
        public
        onlyOperator
    {
        require(newOperator != address(0));

        emit ChangedOperator(operator, newOperator);
        operator = newOperator;
    }

    constructor() public
    {
        operator = msg.sender;

        lastCommittedBlock = 0;
        depositNonce = 1;
        minExitBond = 200000;
    }

    // @param blocks       32 byte merkle headers appended in ascending order
    // @param txnsPerBlock number of transactions per block
    // @param feesPerBlock amount of fees the validator has collected per block
    // @param blockNum     the block number of the first header
    // @notice each block is capped at 2**16-1 transactions
    function submitBlock(bytes32[] memory headers, uint256[] memory txnsPerBlock, uint256[] memory feePerBlock, uint256 blockNum)
        public
        onlyOperator
    {
        require(blockNum == lastCommittedBlock.add(1));
        require(headers.length == txnsPerBlock.length && txnsPerBlock.length == feePerBlock.length);

        for (uint i = 0; i < headers.length && lastCommittedBlock <= lastBlockNum; i++) {
            require(headers[i] != bytes32(0) && txnsPerBlock[i] > 0 && txnsPerBlock[i] < feeIndex);

            lastCommittedBlock = lastCommittedBlock.add(1);
            plasmaChain[lastCommittedBlock] = plasmaBlock({
                header: headers[i],
                numTxns: txnsPerBlock[i],
                feeAmount: feePerBlock[i],
                createdAt: block.timestamp,
                ethBlockNum: block.number
            });

            emit BlockSubmitted(headers[i], lastCommittedBlock, txnsPerBlock[i], feePerBlock[i]);
        }
   }

    // @param owner owner of this deposit
    function deposit(address owner)
        public
        payable
    {
        deposits[depositNonce] = depositStruct(owner, msg.value, block.timestamp, block.number);
        emit Deposit(owner, msg.value, depositNonce, block.number);

        depositNonce = depositNonce.add(uint256(1));
    }

    // @param depositNonce the nonce of the specific deposit
    function startDepositExit(uint256 nonce, uint256 committedFee)
        public
        payable
        isBonded
    {
        require(deposits[nonce].owner == msg.sender);
        require(deposits[nonce].amount > committedFee);
        require(depositExits[nonce].state == ExitState.NonExistent);

        address owner = deposits[nonce].owner;
        uint256 amount = deposits[nonce].amount;
        uint256 priority = block.timestamp << 128 | nonce;
        depositExitQueue.insert(priority);
        depositExits[nonce] = exit({
            owner: owner,
            amount: amount,
            committedFee: committedFee,
            createdAt: block.timestamp,
            position: [0,0,0,nonce],
            state: ExitState.Pending
        });

        emit StartedDepositExit(nonce, owner, amount, committedFee);
    }

    // Transaction encoding:
    // [[Blknum1, TxIndex1, Oindex1, DepositNonce1, Input1ConfirmSig,
    //   Blknum2, TxIndex2, Oindex2, DepositNonce2, Input2ConfirmSig,
    //   NewOwner, Denom1, NewOwner, Denom2, Fee],
    //  [Signature1, Signature2]]
    //
    // All integers are padded to 32 bytes. Input's confirm signatures are 130 bytes for each input.
    // Zero bytes if unapplicable (deposit/fee inputs) Signatures are 65 bytes in length
    //
    // @param txBytes rlp encoded transaction
    // @notice this function will revert if the txBytes are malformed
    function decodeTransaction(bytes memory txBytes)
        internal
        pure
        returns (RLPReader.RLPItem[] memory txList, RLPReader.RLPItem[] memory sigList, bytes32 txHash)
    {
        // entire byte length of the rlp encoded transaction.
        require(txBytes.length == 811);

        RLPReader.RLPItem[] memory spendMsg = txBytes.toRlpItem().toList();
        require(spendMsg.length == 2);

        txList = spendMsg[0].toList();
        require(txList.length == 15);

        sigList = spendMsg[1].toList();
        require(sigList.length == 2);

        // bytes the signatures are over
        txHash = keccak256(spendMsg[0].toRlpBytes());
    }


    // @param txPos             location of the transaction [blkNum, txIndex, outputIndex]
    // @param txBytes           transaction bytes containing the exiting output
    // @param proof             merkle proof of inclusion in the plasma chain
    // @param confSig0          confirm signatures sent by the owners of the first input acknowledging the spend.
    // @param confSig1          confirm signatures sent by the owners of the second input acknowledging the spend (if applicable).
    // @notice `confirmSignatures` and `ConfirmSig0`/`ConfirmSig1` are unrelated to each other.
    // @notice `confirmSignatures` is either 65 or 130 bytes in length dependent on if a second input is present
    // @notice `confirmSignatures` should be empty if the output trying to be exited is a fee output
    function startTransactionExit(uint256[3] memory txPos, bytes memory txBytes, bytes memory proof, bytes memory confirmSignatures, uint256 committedFee)
        public
        payable
        isBonded
    {
        require(txPos[1] < feeIndex);
        uint256 position = calcPosition(txPos);
        require(txExits[position].state == ExitState.NonExistent);

        uint256 amount = startTransactionExitHelper(txPos, txBytes, proof, confirmSignatures);
        require(amount > committedFee);

        // calculate the priority of the transaction taking into account the withdrawal delay attack
        // withdrawal delay attack: https://github.com/FourthState/plasma-mvp-rootchain/issues/42
        uint256 createdAt = plasmaChain[txPos[0]].createdAt;
        txExitQueue.insert(SafeMath.max(createdAt.add(1 weeks), block.timestamp) << 128 | position);

        // write exit to storage
        txExits[position] = exit({
            owner: msg.sender,
            amount: amount,
            committedFee: committedFee,
            createdAt: block.timestamp,
            position: [txPos[0], txPos[1], txPos[2], 0],
            state: ExitState.Pending
        });

        emit StartedTransactionExit(txPos, msg.sender, amount, confirmSignatures, committedFee);
    }

    // @returns amount of the exiting transaction
    // @notice the purpose of this helper was to work around the capped evm stack frame
    function startTransactionExitHelper(uint256[3] memory txPos, bytes memory txBytes, bytes memory proof, bytes memory confirmSignatures)
        private
        view
        returns (uint256)
    {
        bytes32 txHash;
        RLPReader.RLPItem[] memory txList;
        RLPReader.RLPItem[] memory sigList;
        (txList, sigList, txHash) = decodeTransaction(txBytes);

        uint base = txPos[2].mul(2);
        require(msg.sender == txList[base.add(10)].toAddress());

        plasmaBlock memory blk = plasmaChain[txPos[0]];

        // Validation

        bytes32 merkleHash = sha256(txBytes);
        require(merkleHash.checkMembership(txPos[1], blk.header, proof, blk.numTxns));

        address recoveredAddress;
        bytes32 confirmationHash = sha256(abi.encodePacked(merkleHash, blk.header));

        bytes memory sig = sigList[0].toBytes();
        require(sig.length == 65 && confirmSignatures.length % 65 == 0 && confirmSignatures.length > 0 && confirmSignatures.length <= 130);
        recoveredAddress = confirmationHash.recover(confirmSignatures.slice(0, 65));
        require(recoveredAddress != address(0) && recoveredAddress == txHash.recover(sig));
        if (txList[5].toUintStrict() > 0 || txList[8].toUintStrict() > 0) { // existence of a second input
            sig = sigList[1].toBytes();
            require(sig.length == 65 && confirmSignatures.length == 130);
            recoveredAddress = confirmationHash.recover(confirmSignatures.slice(65, 65));
            require(recoveredAddress != address(0) && recoveredAddress == txHash.recover(sig));
        }

        // check that the UTXO's two direct inputs have not been previously exited
        require(validateTransactionExitInputs(txList));

        return txList[base.add(11)].toUintStrict();
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
            uint256 base = uint256(5).mul(i);
            uint depositNonce_ = txList[base.add(3)].toUintStrict();
            if (depositNonce_ == 0) {
                uint256 blkNum = txList[base].toUintStrict();
                uint256 txIndex = txList[base.add(1)].toUintStrict();
                uint256 outputIndex = txList[base.add(2)].toUintStrict();
                uint256 position = calcPosition([blkNum, txIndex, outputIndex]);
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
    function startFeeExit(uint256 blockNumber, uint256 committedFee)
        public
        payable
        onlyOperator
        isBonded
    {
        plasmaBlock memory blk = plasmaChain[blockNumber];
        require(blk.header != bytes32(0));

        uint256 feeAmount = blk.feeAmount;

        // nonzero fee and prevent and greater than the committed fee if spent.
        // default value for a fee amount is zero. Will revert if a block for
        // this number has not been committed
        require(feeAmount > committedFee);

        // a fee UTXO has explicitly defined position [blockNumber, 2**16 - 1, 0]
        uint256 position = calcPosition([blockNumber, feeIndex, 0]);
        require(txExits[position].state == ExitState.NonExistent);

        txExitQueue.insert(SafeMath.max(blk.createdAt.add(1 weeks), block.timestamp) << 128 | position);

        txExits[position] = exit({
            owner: msg.sender,
            amount: feeAmount,
            committedFee: committedFee,
            createdAt: block.timestamp,
            position: [blockNumber, feeIndex, 0, 0],
            state: ExitState.Pending
        });

        // pass in empty bytes for confirmSignatures for StartedTransactionExit event.
        emit StartedTransactionExit([blockNumber, feeIndex, 0], operator, feeAmount, "", 0);
}

    // @param exitingTxPos     position of the invalid exiting transaction [blkNum, txIndex, outputIndex]
    // @param challengingTxPos position of the challenging transaction [blkNum, txIndex]
    // @param txBytes          raw transaction bytes of the challenging transaction
    // @param proof            proof of inclusion for this merkle hash
    // @param confirmSignature signature used to invalidate the invalid exit. Signature is over (merkleHash, block header)
    // @notice The operator can challenge an exit which commits an invalid fee by simply passing in empty bytes for confirm signature as they are not needed.
    //         The committed fee is checked againt the challenging tx bytes
    function challengeExit(uint256[4] memory exitingTxPos, uint256[2] memory challengingTxPos, bytes memory txBytes, bytes memory proof, bytes memory confirmSignature)
        public
    {
        bytes32 txHash;
        RLPReader.RLPItem[] memory txList;
        RLPReader.RLPItem[] memory sigList;
        (txList, sigList, txHash) = decodeTransaction(txBytes);

        // `challengingTxPos` is sequentially after `exitingTxPos`
        require(exitingTxPos[0] < challengingTxPos[0] || (exitingTxPos[0] == challengingTxPos[0] && exitingTxPos[1] < challengingTxPos[1]));

        // must be a direct spend
        bool firstInput = exitingTxPos[0] == txList[0].toUintStrict() && exitingTxPos[1] == txList[1].toUintStrict() && exitingTxPos[2] == txList[2].toUintStrict() && exitingTxPos[3] == txList[3].toUintStrict();
        require(firstInput || exitingTxPos[0] == txList[5].toUintStrict() && exitingTxPos[1] == txList[6].toUintStrict() && exitingTxPos[2] == txList[7].toUintStrict() && exitingTxPos[3] == txList[8].toUintStrict());

        // transaction to be challenged should have a pending exit
        exit storage exit_ = exitingTxPos[3] == 0 ? 
            txExits[calcPosition([exitingTxPos[0], exitingTxPos[1], exitingTxPos[2]])] : depositExits[exitingTxPos[3]];
        require(exit_.state == ExitState.Pending);

        plasmaBlock memory blk = plasmaChain[challengingTxPos[0]];

        bytes32 merkleHash = sha256(txBytes);
        require(blk.header != bytes32(0) && merkleHash.checkMembership(challengingTxPos[1], blk.header, proof, blk.numTxns));

        address recoveredAddress;
        // we check for confirm signatures if:
        // The exiting tx is a first input and commits the correct fee
        // OR
        // The exiting tx is the second input in the challenging transaction
        //
        // If this challenge was a fee mismatch, then we check the first transaction signature
        // to prevent the operator from forging invalid inclusions
        //
        // For a fee mismatch, the state becomes `NonExistent` so that the exit can be reopened.
        // Otherwise, `Challenged` so that the exit can never be opened.
        if (firstInput && exit_.committedFee != txList[14].toUintStrict()) {
            bytes memory sig = sigList[0].toBytes();
            recoveredAddress = txHash.recover(sig);
            require(sig.length == 65 && recoveredAddress != address(0) && exit_.owner == recoveredAddress);

            exit_.state = ExitState.NonExistent;
        } else {
            bytes32 confirmationHash = sha256(abi.encodePacked(merkleHash, blk.header));
            recoveredAddress = confirmationHash.recover(confirmSignature);
            require(confirmSignature.length == 65 && recoveredAddress != address(0) && exit_.owner == recoveredAddress);

            exit_.state = ExitState.Challenged;
        }

        // exit successfully challenged. Award the sender with the bond
        balances[msg.sender] = balances[msg.sender].add(minExitBond);
        totalWithdrawBalance = totalWithdrawBalance.add(minExitBond);
        emit AddedToBalances(msg.sender, minExitBond);

        emit ChallengedExit(exit_.position, exit_.owner, exit_.amount - exit_.committedFee);
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
        if (queue.length == 0) return;

        // retrieve the lowest priority and the appropriate exit struct
        uint256 priority = queue[0];
        exit memory currentExit;
        uint256 position;
        // retrieve the right 128 bits from the priority to obtain the position
        assembly {
   	        position := and(priority, div(not(0x0), exp(256, 16)))
		}

        currentExit = isDeposits ? depositExits[position] : txExits[position];

        /*
        * Conditions:
        *   1. Exits exist
        *   2. Exits must be a week old
        *   3. Funds must exist for the exit to withdraw
        */
        uint256 amountToAdd;
        uint256 challengePeriod = isDeposits ? 5 days : 1 weeks;
        while (block.timestamp.sub(currentExit.createdAt) > challengePeriod &&
               plasmaChainBalance() > 0 &&
               gasleft() > 80000) {

            // skip currentExit if it is not in 'started/pending' state.
            if (currentExit.state != ExitState.Pending) {
                queue.delMin();
            } else {
                // reimburse the bond but remove fee allocated for the operator
                amountToAdd = currentExit.amount.add(minExitBond).sub(currentExit.committedFee);
                
                balances[currentExit.owner] = balances[currentExit.owner].add(amountToAdd);
                totalWithdrawBalance = totalWithdrawBalance.add(amountToAdd);

                if (isDeposits)
                    depositExits[position].state = ExitState.Finalized;
                else
                    txExits[position].state = ExitState.Finalized;

                emit FinalizedExit(currentExit.position, currentExit.owner, amountToAdd);
                emit AddedToBalances(currentExit.owner, amountToAdd);

                // move onto the next oldest exit
                queue.delMin();
            }

            if (queue.length == 0) {
                return;
            }

            // move onto the next oldest exit
            priority = queue[0];
            
            // retrieve the right 128 bits from the priority to obtain the position
            assembly {
   			    position := and(priority, div(not(0x0), exp(256, 16)))
		    }
             
            currentExit = isDeposits ? depositExits[position] : txExits[position];
        }
    }

    // @notice will revert if the output index is out of bounds
    function calcPosition(uint256[3] memory txPos)
        private
        view
        returns (uint256)
    {
        require(validatePostion([txPos[0], txPos[1], txPos[2], 0]));

        uint256 position = txPos[0].mul(blockIndexFactor).add(txPos[1].mul(txIndexFactor)).add(txPos[2]);
        require(position <= 2**128-1); // check for an overflow

        return position;
    }

    function validatePostion(uint256[4] memory position)
        private
        view
        returns (bool)
    {
        uint256 blkNum = position[0];
        uint256 txIndex = position[1];
        uint256 oIndex = position[2];
        uint256 depNonce = position[3];

        if (blkNum > 0) { // utxo input
            // uncommitted block
            if (blkNum > lastCommittedBlock)
                return false;
            // txIndex out of bounds for the block
            if (txIndex >= plasmaChain[blkNum].numTxns && txIndex != feeIndex)
                return false;
            // fee input must have a zero output index
            if (txIndex == feeIndex && oIndex > 0)
                return false;
            // deposit nonce must be zero
            if (depNonce > 0)
                return false;
            // only two outputs
            if (oIndex > 1)
                return false;
        } else { // deposit or fee input
            // deposit input must be zero'd output position
            // `blkNum` is not checked as it will fail above
            if (depNonce > 0 && (txIndex > 0 || oIndex > 0))
                return false;
        }

        return true;
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

        // will revert the above deletion if it fails
        msg.sender.transfer(transferAmount);
        return transferAmount;
    }

    /*
    * Getters
    */

    function plasmaChainBalance()
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

    function txQueueLength()
        public
        view
        returns (uint)
    {
        return txExitQueue.length;
    }

    function depositQueueLength()
        public 
        view
        returns (uint)
    {   
        return depositExitQueue.length;
    }
}
