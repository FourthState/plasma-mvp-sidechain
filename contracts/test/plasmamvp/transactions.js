let RLP = require('rlp');
let assert = require('chai').assert

let PlasmaMVP = artifacts.require('PlasmaMVP');

let {
    fastForward,
    sha256String,
    generateMerkleRootAndProof,
    fillTxList,
} = require('./plasmamvp_helpers.js');

let { toHex, catchError } = require('../utilities.js');

contract('[PlasmaMVP] Transactions', async (accounts) => {
    let instance;
    let oneWeek = 604800; // in seconds
    let authority = accounts[0];
    let minExitBond = 200000;

    // deploy the instance contract before each test.
    // deposit from authority and mine the first block which
    // includes a spend of that full deposit to account[1] (first input)
    let amount = 100; let feeAmount = 10;
    let depositNonce;
    let txPos, txBytes;
    let proof, feeProof;
    let sigs, confirmSignatures;
    beforeEach(async () => {
        instance = await PlasmaMVP.new({from: authority});

        depositNonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(authority, {from: authority, value: amount*2 + 10});

        // deposit is the first input. authority creates two outputs.
        // split equally with a fee amount of 10
        let txList = Array(15).fill(0);
        txList[3] = depositNonce;
        txList[10] = authority; txList[11] = amount;
        txList[12] = authority; txList[13] = amount;
        txList[14] = feeAmount;
        txList = fillTxList(txList);
        let txHash = web3.utils.soliditySha3(toHex(RLP.encode(txList).toString('hex')));
        let sigs = [toHex(await web3.eth.sign(txHash, authority)), toHex(Buffer.alloc(65).toString('hex'))];
        txBytes = [txList, sigs];
        txBytes = RLP.encode(txBytes).toString('hex');

        let merkleHash = sha256String(txBytes);

        // submit the block
        let merkleRoot;
        [merkleRoot, proof] = generateMerkleRootAndProof([merkleHash], 0);
        let blockNum = (await instance.lastCommittedBlock.call()).toNumber() + 1;
        await instance.submitBlock([toHex(merkleRoot)], [1], [10], blockNum, {from: authority});

        // construct the confirm signature
        let confirmHash = sha256String(merkleHash + merkleRoot.slice(2));
        confirmSignatures = await web3.eth.sign(confirmHash, authority);

        txPos = [blockNum, 0, 0];
    });

    it("Will not revert finalizeExit with an empty queue", async () => {
        await instance.finalizeTransactionExits();
    });

    it("Will reject an exit with a committedFee larger than the transaction amount", async () => {
        let err;
        [err] = await catchError(instance.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSignatures), amount+1, {from: accounts[1], value: minExitBond}));

        if (!err)
            assert.fail("Started a transaction exit with a committed fee larger than the tx amount");
    });

    it("Allows only the utxo owner to start an exit", async () => {
        let err;
        [err] = await catchError(instance.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSignatures), 0, {from: accounts[3], value: minExitBond}));
        if (!err)
            assert.fail("address exited an output it does not own");
    });

    it("Can challenge a spend of a utxo with only a valid confirm signature", async () => {
        // spend the first output back to authority
        let txList = Array(15).fill(0);
        txList[0] = txPos[0]; txList[1] = txPos[1]; txList[2] = txPos[2];
        txList[10] = authority; txList[11] = amount;
        txList = fillTxList(txList);
        txHash = web3.utils.soliditySha3(toHex(RLP.encode(txList).toString('hex')));
        let sigs = [await web3.eth.sign(txHash, authority), toHex(Buffer.alloc(65).toString('hex'))];
        let challengingTxBytes = [txList, sigs];
        challengingTxBytes = RLP.encode(challengingTxBytes).toString('hex');
        let merkleHash = sha256String(challengingTxBytes);

        // submit the block
        let challengingBlockNum = (await instance.lastCommittedBlock.call()).toNumber() + 1;
        let [merkleRoot, challengingProof] = generateMerkleRootAndProof([merkleHash], 0);
        await instance.submitBlock([merkleRoot], [1], [0], challengingBlockNum, {from: authority});

        let confirmationHash = sha256String(merkleHash + merkleRoot.slice(2));
        let challengingConfirmSig = await web3.eth.sign(confirmationHash, authority);

        // start an transaction exit of the input
        await instance.startTransactionExit(txPos, toHex(txBytes), toHex(proof), toHex(confirmSignatures), 0, {from: authority, value: minExitBond});

        // attempt to challenge with an invalid confirm signature
        let wrongConfirmSig = await web3.eth.sign(confirmationHash, accounts[1]);
        let [err] = await catchError(instance.challengeExit([...txPos, 0], [challengingBlockNum, 0], toHex(challengingTxBytes), toHex(challengingProof),
            toHex(wrongConfirmSig), {from: accounts[2]}));
        
        // correctly challenge the exit above
        await instance.challengeExit([...txPos, 0], [challengingBlockNum, 0], toHex(challengingTxBytes), toHex(challengingProof),
            toHex(challengingConfirmSig), {from: accounts[2]});

        // check that the bond has been rewarded to the challenger
        let balance = (await instance.balanceOf.call(accounts[2])).toNumber();
        assert.equal(balance, minExitBond, "exit bond not rewarded to challenger");

        let position = 1000000*txPos[0] + 10*txPos[1];
        let exit = await instance.txExits.call(position);
        assert.equal(exit[4].toNumber(), 2, "Fee exit state is not Challenged");

        // cannot reopen the exit
        [err] = await catchError(instance.startTransactionExit(txPos, toHex(txBytes), toHex(proof),
            toHex(confirmSignatures), 0, {from: authority, value: minExitBond}));
        if (!err)
            assert.fail("reopened challenged exit");
    });

    it("Catches StartedTransactionExit event", async () => {
        let tx = await instance.startTransactionExit(txPos, toHex(txBytes),
            toHex(proof), toHex(confirmSignatures), 0, {from: authority, value: minExitBond});

        assert.equal(tx.logs[0].args.position.toString(), txPos.toString(), "StartedTransactionExit event emits incorrect priority");
        assert.equal(tx.logs[0].args.owner, authority, "StartedTransactionExit event emits incorrect owner");
        assert.equal(tx.logs[0].args.amount.toNumber(), amount, "StartedTransactionExit event emits incorrect amount");
        assert.equal(tx.logs[0].args.confirmSignatures, toHex(confirmSignatures), "StartedTransactionExit event does not emit confirm signatures");
    });

    it("Can start and finalize a transaction exit", async () => {
        await instance.startTransactionExit(txPos, toHex(txBytes),
            toHex(proof), toHex(confirmSignatures), 0, {from: authority, value: minExitBond});

        await fastForward(oneWeek + 10);

        await instance.finalizeTransactionExits();

        let balance = (await instance.balanceOf.call(authority)).toNumber();
        assert.equal(balance, amount + minExitBond);

        let position = 1000000*txPos[0];
        let exit = await instance.txExits.call(position);
        assert.equal(exit[4].toNumber(), 3, "exit's state not set to finalized");
    });

    it("Only authority to start a fee withdrawal exit", async () => {
        // non-authoritys cannot start fee exits
        let err;
        [err] = await catchError(instance.startFeeExit(txPos[0], 0, {from: accounts[1], value: minExitBond}));
        if (!err)
            assert.fail("fee exit start from non-authority");

        // authority cannot start a fee exit without putting a sufficient exit bond
        [err] = await catchError(instance.startFeeExit(txPos[0], 0, {from: authority, value: minExitBond - 100}));
        if (!err)
            assert.fail("started fee exit with insufficient bond");

        // the committed fee must be less than the fee amount
        [err] = await catchError(instance.startFeeExit(txPos[0], feeAmount+10, {from: authority, value: minExitBond}));
        if (!err)
            assert.fail("started fee exit with a commited fee larger than the fee amount");

        // cannot start a fee exit for a non-existent block
        let nonExistentBlockNum = txPos[0] + 100;
        [err] = await catchError(instance.startFeeExit(nonExistentBlockNum, 0, {from: authority, value: minExitBond}));
        if (!err)
            assert.fail("started fee exit for non-existent block");

        // authority can start a fee exit with sufficient exit bond
        await instance.startFeeExit(txPos[0], 0, {from: authority, value: minExitBond});

        let position = 1000000*txPos[0] + 10*(Math.pow(2, 16) - 1);
        let feeExit = await instance.txExits.call(position);
        assert.equal(feeExit[0].toNumber(), 10, "Incorrect fee exit amount");
        assert.equal(feeExit[3], authority, "Incorrect fee exit owner");
        assert.equal(feeExit[4].toNumber(), 1, "Incorrect exit state.");

        // can only start a fee exit for any particular block once
        [err] = await catchError(instance.startFeeExit(txPos[0], 0, {from: authority, value: minExitBond}));
        if (!err)
            assert.fail("attempted the same exit while a pending one existed");

        // fast forward one week
        await fastForward(oneWeek+10);

        // finalize exits
        await instance.finalizeTransactionExits();

        let balance = (await instance.balanceOf.call(authority)).toNumber();
        assert.equal(balance, feeAmount + minExitBond, "authority has incorrect balance");

        feeExit = await instance.txExits.call(position);
        assert.equal(feeExit[4].toNumber(), 3, "Fee exit state is not Finalized");
    });

    it("Allows authority to challenge only a incorrect committed fee", async () => {
        // spend both inputs with a fee. fee should only come from the first input
        let txList2 = Array(15).fill(0);
        txList2[0] = txPos[0]; txList2[1] = txPos[1]; txList2[2] = txPos[2];
        txList2[5] = txPos[0]; txList2[6] = txPos[1]; txList2[7] = 1;
        txList2[10] = authority; txList2[11] = amount - 5;
        txList2[12] = authority; txList2[13] = amount;
        txList2[14] = 5; // fee
        txList2 = fillTxList(txList2);
        let txHash2 = web3.utils.soliditySha3(toHex(RLP.encode(txList2).toString('hex')));
        let sigs2 = [toHex(await web3.eth.sign(txHash2, authority)), toHex(await web3.eth.sign(txHash2, authority))];
        let txBytes2 = [txList2, sigs2];
        txBytes2 = RLP.encode(txBytes2).toString('hex');

        // submit the block and claim the fee
        let merkleHash2 = sha256String(txBytes2);
        let [merkleRoot2, proof2] = generateMerkleRootAndProof([merkleHash2], 0);
        let blockNum2 = (await instance.lastCommittedBlock.call()).toNumber() + 1;
        await instance.submitBlock([toHex(merkleRoot2)], [1], [5], blockNum2, {from: authority});

        // first input exits without committing to the fee
        await instance.startTransactionExit(txPos, toHex(txBytes), toHex(proof),
            toHex(confirmSignatures), 1, {from: authority, value: minExitBond});

        // second input will exit (not forced to commit the fee)
        let secondOutput = [txPos[0], txPos[1], 1];
        await instance.startTransactionExit(secondOutput, toHex(txBytes), toHex(proof),
            toHex(confirmSignatures), 0, {from: authority, value: minExitBond});

        // operator will try incorrectly challenge with second output
        let err;
        [err] = await catchError(instance.challengeExit([...secondOutput, 0], [blockNum2, 0], toHex(txBytes2), toHex(proof2), toHex("")));
        if (!err)
            assert.fail("operator challenged with the second input");

        // operator will challenge the exit with the first input
        await instance.challengeExit([...txPos, 0], [blockNum2, 0], toHex(txBytes2), toHex(proof2), toHex(""));

        let exit = await instance.txExits.call(1000000*txPos[0] + 10*txPos[1]);
        assert.equal(exit[4].toNumber(), 0, "exit with incorrect fee not changed to non existent after challenge");

        // should not be able to challenge an exit which does not exist
        [err] = await catchError(instance.challengeExit([...txPos, 0], [blockNum2, 0], toHex(txBytes2), toHex(proof2), toHex("")));
        if (!err)
            assert.fail("operator challenged an exit which does not exist");

        // first input will exit with the correct committed fee
        await instance.startTransactionExit(txPos, toHex(txBytes), toHex(proof),
            toHex(confirmSignatures), 5, {authority, value: minExitBond});

        // operator will challenge the exit and will fail without a confirm signature
        [err] = await catchError(instance.challengeExit([...txPos, 0], [blockNum2, 0], toHex(txBytes2), toHex(proof2), toHex("")));
        if (!err)
            assert.fail("operator challenged an exit with the correct committed fee");
    });
    
    it("Transaction signatures are checked during a fee mismatch challenge", async () => {
        // spend both inputs with a fee. fee should only come from the first input
        let txList2 = Array(15).fill(0);
        txList2[0] = txPos[0]; txList2[1] = txPos[1]; txList2[2] = txPos[2];
        txList2[5] = txPos[0]; txList2[6] = txPos[1]; txList2[7] = 1;
        txList2[10] = authority; txList2[11] = amount - 5;
        txList2[12] = authority; txList2[13] = amount;
        txList2[14] = 5; // fee
        txList2 = fillTxList(txList2);
        let txHash2 = web3.utils.soliditySha3(toHex(RLP.encode(txList2).toString('hex')));
        // incorrect sig
        let sigs2 = [toHex(await web3.eth.sign(txHash2, accounts[1])), toHex(await web3.eth.sign(txHash2, authority))];
        let txBytes2 = [txList2, sigs2];
        txBytes2 = RLP.encode(txBytes2).toString('hex');

        // submit the block and claim the fee
        let merkleHash2 = sha256String(txBytes2);
        let [merkleRoot2, proof2] = generateMerkleRootAndProof([merkleHash2], 0);
        let blockNum2 = (await instance.lastCommittedBlock.call()).toNumber() + 1;
        await instance.submitBlock([toHex(merkleRoot2)], [1], [5], blockNum2, {from: authority});

        // first input exits not committing the fee
        await instance.startTransactionExit(txPos, toHex(txBytes), toHex(proof),
            toHex(confirmSignatures), 0, {from: authority, value: minExitBond});

        // operator will try challenge with the invalid spend
        let [err] = await catchError(instance.challengeExit([...txPos, 0], [blockNum2, 0],
            toHex(txBytes2), toHex(proof2), toHex("")));
        if (!err)
            assert.fail("challenge fee mismatch with an invalid inclusion");
    });

    it("Challenge a fee withdrawal exit", async () => {
        // spend the fee in txPos[0]
        let txList = Array(15).fill(0);
        txList[0] = txPos[0]; txList[1] = Math.pow(2, 16) - 1;
        txList[10] = authority; txList[11] = feeAmount;
        txList = fillTxList(txList);
        let txHash = sha256String(toHex(RLP.encode(txList).toString('hex')));
        let sigs = [await web3.eth.sign(txHash, authority), toHex(Buffer.alloc(65).toString('hex'))];
        let txBytes = [txList, sigs];
        txBytes = RLP.encode(txBytes).toString('hex');

        // submit the block
        let merkleHash = sha256String(txBytes);
        let [merkleRoot, proof] = generateMerkleRootAndProof([merkleHash], 0);
        let challengingBlockNum = (await instance.lastCommittedBlock.call()).toNumber() + 1;
        await instance.submitBlock([toHex(merkleRoot)], [1], [0], challengingBlockNum, {from: authority});

        // create the confirm sig
        let confirmHash = sha256String(merkleHash + merkleHash.slice(2));
        let challengingConfirmSig = await web3.eth.sign(confirmHash, authority);

        // start fee exit
        await instance.startFeeExit(txPos[0], 0, {from: authority, value: minExitBond});

        let position = 1000000*txPos[0] + 10*(Math.pow(2, 16) - 1);
        let feeExit = await instance.txExits.call(position);
        assert.equal(feeExit[0].toNumber(), feeAmount, "Incorrect fee exit amount");
        assert.equal(feeExit[3], authority, "Incorrect fee exit owner");
        assert.equal(feeExit[4].toNumber(), 1, "Fee exit state is not Pending");

        // challenge fee exit
        await instance.challengeExit([txPos[0], Math.pow(2, 16) - 1, 0, 0], [challengingBlockNum, 0],
            toHex(txBytes), toHex(proof), toHex(challengingConfirmSig), {from: accounts[2]});

        let balance = (await instance.balanceOf.call(accounts[2])).toNumber();
        assert.equal(balance, minExitBond, "exit bond not rewarded to challenger");

        feeExit = await instance.txExits.call(position);
        assert.equal(feeExit[4].toNumber(), 2, "Fee exit state is not Challenged");

        // cannot reopen the fee exit
        let [err] = await catchError(instance.startFeeExit(txPos[0], 0, {from: authority, value: minExitBond}));
        if (!err)
            assert.fail("reopend a challenged fee exit");
    });

    it("Requires sufficient bond and refunds excess if overpayed", async () => {
        let [err] = await catchError(instance.startTransactionExit(txPos, toHex(txBytes),
            toHex(proof), toHex(confirmSignatures), 0, {from: authority, value: minExitBond - 100}));
        if (!err)
            assert.fail("started exit with insufficient bond");

        await instance.startTransactionExit(txPos, toHex(txBytes),
            toHex(proof), toHex(confirmSignatures), 0, {from: authority, value: minExitBond + 100});

        let balance = (await instance.balanceOf(authority)).toNumber();
        assert.equal(balance, 100, "excess funds not repayed back to caller");
    });

    it("Cannot exit a utxo with an input pending an exit", async () => {
        await instance.startDepositExit(depositNonce, 0, {from: authority, value: minExitBond});

        let err;
        [err] = await catchError(instance.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSignatures), 0,
            {from: accounts[1], value: minExitBond}));

        if (!err)
            assert.fail("started an exit with an input who has a pending exit state");
    });

    it("Attempt a withdrawal delay attack", async () => {
        let fiveDays = 432000 // in seconds

        // authority spends (txPos[0], 0, 1) utxo, sends 1 utxo to themself and the other to accounts[1]
        let txList2 = Array(15).fill(0);
        txList2[0] = txPos[0]; txList2[2] = 1; // first input
        txList2[10] = authority; txList2[11] = amount / 2; // first output
        txList2[12] = accounts[1]; txList2[13] = amount / 2; // second output
        txList2 = fillTxList(txList2);
        let txHash2 = web3.utils.soliditySha3(toHex(RLP.encode(txList2).toString('hex')));
        let sigs2 = [toHex(await web3.eth.sign(txHash2, authority)), toHex(Buffer.alloc(65).toString('hex'))];
        let txBytes2 = RLP.encode([txList2, sigs2]).toString('hex');
        let merkleHash2 = sha256String(txBytes2);

        let root2, proof2;
        [root2, proof2] = generateMerkleRootAndProof([merkleHash2], 0);
        let blockNum2 = (await instance.lastCommittedBlock.call()).toNumber() + 1;
        await instance.submitBlock([toHex(root2)], [1], [0], blockNum2, {from: authority});

        // create confirmation signature
        let confirmationHash2 = sha256String(merkleHash2 + root2.slice(2));
        let confirmSigs2 = await web3.eth.sign(confirmationHash2, authority);

        // make utxos > 1 week old
        await fastForward(oneWeek + 100);

        // start exit for accounts[1], last utxo to be created
        await instance.startTransactionExit([blockNum2, 0, 1],
            toHex(txBytes2), toHex(proof2), toHex(confirmSigs2), 0, {from: accounts[1], value: minExitBond});

        // increase time slightly, so exit by accounts[1] has better priority than authority
        await fastForward(10);

        // start exit for authority utxo
        await instance.startTransactionExit([blockNum2, 0, 0],
            toHex(txBytes2), toHex(proof2), toHex(confirmSigs2), 0, {from: authority, value: minExitBond});

        // Fast Forward ~5 days
        await fastForward(fiveDays);

        // Check to make sure challenge period has not ended
        let position = 1000000 * blockNum2 + 1;
        let currExit = await instance.txExits.call(position);
        assert.ok((currExit[2] + 604800) > (await web3.eth.getBlock("latest")).timestamp);

        // start exit for oldest utxo avaliable
        await instance.startTransactionExit([txPos[0], 0, 0],
            toHex(txBytes), toHex(proof), toHex(confirmSignatures), 0, {from: authority, value: minExitBond});

        // Fast Forward < 1 week
        await fastForward(fiveDays);

        // finalize exits should finalize accounts[1] then authority
        let finalizedExits = await instance.finalizeTransactionExits({from: authority});
        let finalizedExit = await instance.txExits.call(position);
        assert.equal(finalizedExits.logs[0].args.position.toString(), [blockNum2, 0, 1, 0].toString(), "Incorrect position for finalized tx");
        assert.equal(finalizedExits.logs[0].args.owner, accounts[1], "Incorrect finalized exit owner");
        assert.equal(finalizedExits.logs[0].args.amount.toNumber(), amount/2 + minExitBond, "Incorrect finalized exit amount.");
        assert.equal(finalizedExit[4].toNumber(), 3, "Incorrect finalized exit state.");

        // Check other exits
        position = 1000000 * blockNum2;
        finalizedExit = await instance.txExits.call(position);
        assert.equal(finalizedExits.logs[2].args.position.toString(), [blockNum2, 0, 0, 0].toString(), "Incorrect position for finalized tx");
        assert.equal(finalizedExits.logs[2].args.owner, authority, "Incorrect finalized exit owner");
        assert.equal(finalizedExits.logs[2].args.amount.toNumber(), amount/2 + minExitBond, "Incorrect finalized exit amount.");
        assert.equal(finalizedExit[4].toNumber(), 3, "Incorrect finalized exit state.");

        // Last exit should still be pending
        position = 1000000 * txPos[0];
        let pendingExit = await instance.txExits.call(position);
        assert.equal(pendingExit[3], authority, "Incorrect pending exit owner");
        assert.equal(pendingExit[0].toNumber(), amount, "Incorrect pending exit amount");
        assert.equal(pendingExit[4].toNumber(), 1, "Incorrect pending exit state.");

        // Fast Forward rest of challenge period
        await fastForward(oneWeek + 1000);
        await instance.finalizeTransactionExits({from: authority});
        // Check that last exit was processed
        finalizedExit = await instance.txExits.call(position);
        assert.equal(finalizedExit[4].toNumber(), 3, "Incorrect finalized exit state.");
    });

    it("Reverts if finalizeExit runs out of gas", async () => {
        // exit both outputs
        let txPos2 = [txPos[0], 0, 1];
        await instance.startTransactionExit(txPos, toHex(txBytes), toHex(proof),
            toHex(confirmSignatures), 0, {from: authority, value: minExitBond});
        await instance.startTransactionExit(txPos2, toHex(txBytes), toHex(proof),
            toHex(confirmSignatures), 0, {from: authority, value: minExitBond});

        await fastForward(oneWeek + 1000);

        // Only provide enough gas for 1 txn to be finalized
        await instance.finalizeTransactionExits({gas: 120000});

        // The first utxo should have been exited correctly
        let balance = (await instance.balanceOf.call(authority)).toNumber();
        assert.equal(balance, amount + minExitBond);

        let position = 1000000*txPos[0];
        let exit = await instance.txExits.call(position);
        assert.equal(exit[4].toNumber(), 3, "first exit's state not set to finalized");

        position = 1000000*txPos2[0] + 1;
        exit = await instance.txExits.call(position);
        assert.equal(exit[4].toNumber(), 1, "second exit should still be pending");
    });

    it("Requires two correct confirm signatures with two inputs", async () => {
        // spend both outputs
        let txList = Array(15).fill(0);
        txList[0] = txPos[0]; txList[5] = txPos[0]; txList[8] = 1;
        txList[10] = authority; txList[11] = amount*2;
        txList = fillTxList(txList);
        let txHash = web3.utils.soliditySha3(toHex(RLP.encode(txList).toString('hex')));
        let sig = toHex(await web3.eth.sign(txHash, authority));
        let sigs = [sig, sig];
        let txBytes = [txList, sigs];
        txBytes = RLP.encode(txBytes).toString('hex');
        let merkleHash = sha256String(txBytes);

        let blockNum = (await instance.lastCommittedBlock.call()).toNumber() + 1;
        let [merkleRoot, proof] = generateMerkleRootAndProof([merkleHash], 0);
        await instance.submitBlock([toHex(merkleRoot)], [1], [0], blockNum, {from: authority});

        let confirmHash = sha256String(merkleHash + merkleRoot.slice(2));
        let confirmSig = (await web3.eth.sign(confirmHash, authority)).slice(2);
        let incorrectConfirmSig = (await web3.eth.sign(confirmHash, accounts[1])).slice(2);

        // exit the new output with incorrect sigs
        let [err] = await catchError(instance.startTransactionExit([blockNum, 0, 0], toHex(txBytes),
            toHex(proof), toHex(incorrectConfirmSig + confirmSig), 0, {from: authority, value: minExitBond}));
        if (!err)
            assert.fail("started exit with incorrect first confirm sig");
        
        [err] = await catchError(instance.startTransactionExit([blockNum, 0, 0], toHex(txBytes),
            toHex(proof), toHex(confirmSig + incorrectConfirmSig), 0, {from: authority, value: minExitBond}));
        if (!err)
            assert.fail("started exit with incorrect second confirm sig");

        // start successfully
        await instance.startTransactionExit([blockNum, 0, 0], toHex(txBytes), toHex(proof),
            toHex(confirmSig + confirmSig), 0, {from: authority, value: minExitBond});
    });
});
