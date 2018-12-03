let RLP = require('rlp');
let assert = require('chai').assert

let RootChain = artifacts.require('RootChain');

let {
    fastForward,
    sha256String,
    generateMerkleRootAndProof
} = require('./rootchain_helpers.js');

let { toHex, catchError } = require('../utilities.js');

contract('[RootChain] Transactions', async (accounts) => {
    let rootchain;
    let one_week = 604800; // in seconds
    let authority = accounts[0];
    let minExitBond = 10000;

    // deploy the rootchain contract before each test.
    // deposit from accounts[0] and mine the first block which
    // includes a spend of that full deposit to account[1] (first input)
    let amount = 100;
    let depositNonce;
    let txPos, txBytes;
    let proof;
    let sigs, confirmSignatures;
    beforeEach(async () => {
        rootchain = await RootChain.new({from: authority});

        depositNonce = (await rootchain.depositNonce.call()).toNumber();
        await rootchain.deposit(accounts[0], {from: accounts[0], value: amount});

        // deposit is the first input. accounts[0] sends entire deposit to accounts[1]
        let txList = Array(17).fill(0);
        txList[3] = depositNonce; txList[12] = accounts[1]; txList[13] = amount;
        let txHash = web3.sha3(RLP.encode(txList).toString('hex'), {encoding: 'hex'});

        let sigs = [toHex(await web3.eth.sign(accounts[0], txHash)), toHex(Buffer.alloc(65).toString('hex'))];

        txBytes = [txList, sigs];
        txBytes = RLP.encode(txBytes).toString('hex');

        // submit the block
        let merkleHash = sha256String(txBytes);
        let merkleRoot;
        [merkleRoot, proof] = generateMerkleRootAndProof([merkleHash], 0);
        let blockNum = (await rootchain.lastCommittedBlock.call()).toNumber() + 1;
        await rootchain.submitBlock(toHex(merkleRoot), [1], blockNum, {from: authority});

        // construct the confirm signature
        let confirmHash = sha256String(merkleHash + merkleRoot.slice(2));
        confirmSignatures = await web3.eth.sign(accounts[0], confirmHash);

        txPos = [blockNum, 0, 0];
    });

    it("Will not revert finalizeExit with an empty queue", async () => {
        await rootchain.finalizeDepositExits();
        await rootchain.finalizeTransactionExits();
    });

    it("Allows only the utxo owner to start an exit (hardcoded)", async () => {
        rootchain = await RootChain.new({from: authority});

        // utxo information
        // this utxo input and the merkle root of its block were generated
        // by the side chain
        let txPos = [2, 1, 0];
        let txBytes = "0xf8ebf86180808002940e02ce999156cf4e5a30d91b79329e1f01d61379c080808080940000000000000000000000000000000000000000c09453bb5e06573dbd3baeff3710c860f09f06c4c8a4329400000000000000000000000000000000000000008032f886b841288caa04324245958feb44f9ef5d483618b2cfea74622af8a1075a4089be513001cca34ed1230a20849b5c2b4ae33b3e24f4b36cf1d2d2dc45f1485c6f2c03a600b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000";
        let proof = "0xf17d0ac90940e6055a992ac3f76742a2ab47c504b495c6b1accdf839ac018814";
        let confirmSigs = "0xff05e0519b90b7b3f0d9d8a73a2792d55413f0f5626901d45ab0c8adedb668b638d1eac62cd88cd0719d8bfe0c47357c27463677c0997ad3337a0baed7fd6d6600";
        let total = 2;

        // submit block roots
        let root1 = web3.sha3('1234').slice(2);
        // this side chain block contains 2 txns
        let root2 = "783842a0f2aacc2f988d0d9736aac13a0530f1c78d55ab468a1debcd6b42f109";
        let roots = root1 + root2;

        await rootchain.submitBlock(toHex(roots), [1, total], 1, {from: authority});

        let newOwner = "0x53bB5E06573dbD3baEFF3710c860F09F06C4C8A4";

        // attempt to start an transaction exit
        await rootchain.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSigs), {from: newOwner, value: minExitBond});
    });

    it("Can challenge a spend of a utxo (hardcoded)", async () => {
        rootchain = await RootChain.new({from: authority});

        // utxo information
        // this utxo input and the merkle root of its block were generated
        // by the side chain
        let txPos = [2, 1, 0];
        let txBytes = "0xf8ebf86180808002940e02ce999156cf4e5a30d91b79329e1f01d61379c080808080940000000000000000000000000000000000000000c09453bb5e06573dbd3baeff3710c860f09f06c4c8a4329400000000000000000000000000000000000000008032f886b841288caa04324245958feb44f9ef5d483618b2cfea74622af8a1075a4089be513001cca34ed1230a20849b5c2b4ae33b3e24f4b36cf1d2d2dc45f1485c6f2c03a600b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000";
        let proof = "0xf17d0ac90940e6055a992ac3f76742a2ab47c504b495c6b1accdf839ac018814";
        let confirmSigs = "0xff05e0519b90b7b3f0d9d8a73a2792d55413f0f5626901d45ab0c8adedb668b638d1eac62cd88cd0719d8bfe0c47357c27463677c0997ad3337a0baed7fd6d6600";
        let total = 2;

        // spend the utxo
        // this utxo input and the merkle root of its block were generated
        // by the side chain
        let newTxPostxPos = [3, 0, 0];
        let newTxBytes = "0xf9012ff8a5020180809453bb5e06573dbd3baeff3710c860f09f06c4c8a4f843b841ff05e0519b90b7b3f0d9d8a73a2792d55413f0f5626901d45ab0c8adedb668b638d1eac62cd88cd0719d8bfe0c47357c27463677c0997ad3337a0baed7fd6d660080808080940000000000000000000000000000000000000000c0940e02ce999156cf4e5a30d91b79329e1f01d61379199400000000000000000000000000000000000000008019f886b841ff4bafce58d0752731eec71617c68d256058b518e015bb7b0f85a053e491868964be5cd62b744362cb744e20de55f7123ca80115058ff9c1b76d7f31589bf99b01b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000";
        let newProof = "";
        let newConfirmSigs = "0x3dd7595c79ddbf6ea25ccb64f96218a80ecb435369fedcb865e54a87bd8464282210148512342574a6bf8e6f7a916bf4f0e65e1601cbe0b8124718b85860bb7e01";
        let newTotal = 1;

        // submit block roots
        let root1 = web3.sha3('1234').slice(2);
        // this side chain block contains 2 txns
        let root2 = "783842a0f2aacc2f988d0d9736aac13a0530f1c78d55ab468a1debcd6b42f109";
        // this side chain block contains 1 txn
        let root3 = "0501f4b09300d277cdfedb8c6d4747919bbbf454ef6ba9d300796e2703bf444c";
        let roots = root1 + root2 + root3;

        let blockNum = (await rootchain.lastCommittedBlock.call()).toNumber() + 1;
        await rootchain.submitBlock(toHex(roots), [1, total, newTotal], blockNum, {from: authority});

        let newOwner = "0x53bB5E06573dbD3baEFF3710c860F09F06C4C8A4";

        // attempt to start an transaction exit
        let tx = await rootchain.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSigs), {from: newOwner, value: minExitBond});

        // challenge the exit above
        await rootchain.challengeTransactionExit(txPos, newTxPostxPos,
            toHex(newTxBytes), toHex(newProof), toHex(newConfirmSigs),
            {from: accounts[2]});

        // check that the bond has been rewarded to the challenger
        let balance = (await rootchain.balanceOf.call(accounts[2])).toNumber();
        assert.equal(balance, minExitBond, "exit bond not rewarded to challenger");
    });

    it("Catches StartedTransactionExit event", async () => {
        let tx = await rootchain.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSignatures),
            {from: accounts[1], value: minExitBond});

        assert.equal(tx.logs[0].args.position.toString(), txPos.toString(), "StartedTransactionExit event emits incorrect priority");
        assert.equal(tx.logs[0].args.owner, accounts[1], "StartedTransactionExit event emits incorrect owner");
        assert.equal(tx.logs[0].args.amount.toNumber(), amount, "StartedTransactionExit event emits incorrect amount");
        assert.equal(tx.logs[0].args.confirmSignatures, toHex(confirmSignatures), "StartedTransactionExit event does not emit confirm signatures");
    });

    it("Can start and finalize a transaction exit", async () => {
        await rootchain.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSignatures),
            {from: accounts[1], value: minExitBond});

        fastForward(one_week + 1000);

        await rootchain.finalizeTransactionExits();

        let balance = (await rootchain.balanceOf.call(accounts[1])).toNumber();
        assert.equal(balance, amount + minExitBond);

        let position = 1000000*txPos[0];
        let exit = await rootchain.txExits.call(position);
        assert.equal(exit[3].toNumber(), 3, "exit's state not set to finalized");
    });

    it("Requires sufficient bond and refunds excess if overpayed", async () => {
        let err;
        [err] = await catchError(rootchain.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSignatures),
            {from: accounts[1], value: minExitBond - 100}));
        if (!err)
            assert.fail("started exit with insufficient bond");

        await rootchain.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSignatures),
            {from: accounts[1], value: minExitBond + 100});

        let balance = (await rootchain.balanceOf(accounts[1])).toNumber();
        assert.equal(balance, 100, "excess funds not repayed back to caller");
    });

    it("Only allows exiting a utxo once", async () => {
        await rootchain.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSignatures),
            {from: accounts[1], value: minExitBond});

        let err;
        [err] = await catchError(rootchain.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSignatures),
            {from: accounts[1], value: minExitBond}));

        if (!err)
            assert.fail("reopened the same exit while already a pending one existed");

        fastForward(one_week + 100);

        [err] = await catchError(rootchain.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSignatures),
            {from: accounts[1], value: minExitBond}));

        if (!err)
            assert.fail("reopened the same exit after already finalized");
    });

    it("Cannot exit a utxo with an input pending an exit", async () => {
        await rootchain.startDepositExit(depositNonce, {from: accounts[0], value: minExitBond});

        let err;
        [err] = await catchError(rootchain.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSignatures),
            {from: accounts[1], value: minExitBond}));

        if (!err)
            assert.fail("started an exit with an input who has a pending exit state");
    });

    it("Can challenge a spend of a utxo", async () => {
        // spend all funds to account[2] and mine the block
        // deposit is the first input. spending entire deposit to accounts[1]
        let txList2 = Array(17).fill(0);
        txList2[0] = txPos[0]; txList2[1] = txPos[1]; txList2[2] = txPos[2]; // first input
        txList2[12] = accounts[2]; txList2[13] = amount; // first output

        // create signature by deposit owner. Second signature should be zero
        let txHash = web3.sha3(RLP.encode(txList2).toString('hex'), {encoding: 'hex'});
        let sigs = [toHex(await web3.eth.sign(accounts[1], txHash)), toHex(Buffer.alloc(65).toString('hex'))]

        let newTxBytes = [txList2, sigs];
        newTxBytes = RLP.encode(newTxBytes).toString('hex');

        // include this transaction in the next block
        let merkleHash = sha256String(newTxBytes);
        let root, proof2;
        [root, proof2] = generateMerkleRootAndProof([merkleHash], 0);
        let blockNum = (await rootchain.lastCommittedBlock.call()).toNumber() + 1;
        await rootchain.submitBlock(toHex(root), [1], blockNum, {from: authority});

        // create the confirm sig
        let confirmHash = sha256String(merkleHash + root.slice(2));
        let newConfirmSignatures = await web3.eth.sign(accounts[1], confirmHash);

        // start an exit of the original utxo
        await rootchain.startTransactionExit(txPos,
            toHex(txBytes), toHex(proof), toHex(confirmSignatures),
            {from: accounts[1], value: minExitBond});

        // try to exit this new utxo and realize it cannot. child has a pending exit
        let err;
        [err] = await catchError(rootchain.startTransactionExit([blockNum, 0, 0],
            toHex(newTxBytes), toHex(proof2), toHex(newConfirmSignatures),
            {from: accounts[2], value: minExitBond}));
        if (!err)
            assert.fail("started exit when the child has a pending exit");

        // matching input required
        [err] = await catchError(rootchain.challengeTransactionExit([txPos[0], 0, 1], [blockNum, 0, 0],
            toHex(newTxBytes), toHex(proof2), toHex(newConfirmSignatures.substring(0,65),
            {from: accounts[2]})));
        if (!err)
            assert.fail("challenged with transaction that is not a direct child");

        // challenge
        await rootchain.challengeTransactionExit(txPos, [blockNum, 0, 0],
            toHex(newTxBytes), toHex(proof2), toHex(newConfirmSignatures),
            {from: accounts[2]});

        let balance = (await rootchain.balanceOf.call(accounts[2])).toNumber();
        assert.equal(balance, minExitBond, "exit bond not rewarded to challenger");

        // start an exit of the new utxo after successfully challenging
        await rootchain.startTransactionExit([blockNum, 0, 0],
            toHex(newTxBytes), toHex(proof2), toHex(newConfirmSignatures),
            {from: accounts[2], value: minExitBond});
    });

    it("Rejects exiting a transaction whose sole input is the second", async () => {
        let nonce = (await rootchain.depositNonce.call()).toNumber();
        await rootchain.deposit(accounts[2], {from: accounts[2], value: 100});

        // construct transcation with second input as the deposit
        let txList2 = Array(17).fill(0);
        txList2[9] = nonce; txList2[12] = accounts[1]; txList2[13] = 100;
        let txHash = web3.sha3(RLP.encode(txList2, {encoding: 'hex'}));

        // create signature by deposit owner. Second signature should be zero
        let sigs = [toHex(Buffer.alloc(65).toString('hex')), toHex(await web3.eth.sign(accounts[2], txHash))];

        let newTxBytes = RLP.encode([txList2, sigs]).toString('hex');

        let merkleHash = sha256String(newTxBytes);

        // include this transaction in the next block
        let root, proof2;
        [root, proof2] = generateMerkleRootAndProof([merkleHash], 0);
        let blockNum = (await rootchain.lastCommittedBlock.call()).toNumber() + 1;
        await rootchain.submitBlock(toHex(root), [1], blockNum, {from: authority});

        // create the confirm sig
        let confirmHash = sha256String(merkleHash + root.slice(2));
        let confirmSig = await web3.eth.sign(accounts[2], confirmHash);

        let err;
        [err] = await catchError(rootchain.startTransactionExit([blockNum, 0, 0],
            toHex(newTxBytes), toHex(proof2), toHex(confirmSig), {from: accounts[1], value: minExitBond}));
        if (!err)
            assert.fail("With a nonzero second input, two confirm signatures should have been required");

    });

    it("Cannot challenge with an incorrect transaction", async () => {
        // account[1] spends deposit and creates two utxos for themselves
        let txList1 = Array(17).fill(0);
        txList1[0] = txPos[0]; txList1[1] = txPos[1]; txList1[2] = txPos[2]; // first input
        txList1[12] = accounts[1]; txList1[13] = amount/2; // first utxo
        txList1[14] = accounts[1]; txList1[15] = amount/2; // second utxo

        // include this tx the next block
        let txHash1 = web3.sha3(RLP.encode(txList1).toString('hex'), {encoding: 'hex'});
        let sigs1 = [toHex(await web3.eth.sign(accounts[1], txHash1)), toHex(Buffer.alloc(65).toString('hex'))];

        let txBytes1 = RLP.encode([txList1, sigs1]).toString('hex');

        let merkleHash1 = sha256String(txBytes1);
        let root1, proof1;
        [root1, proof1] = generateMerkleRootAndProof([merkleHash1], 0);
        let blockNum1 = (await rootchain.lastCommittedBlock.call()).toNumber() + 1;
        await rootchain.submitBlock(toHex(root1), [1], blockNum1, {from: authority});

        // create confirmation signature
        let confirmationHash1 = sha256String(merkleHash1.slice(2) + root1.slice(2));
        let confirmSigs1 = await web3.eth.sign(accounts[1], confirmationHash1);

        // accounts[1] spends the first output to accounts[2]
        let txList2 = Array(17).fill(0);
        txList2[0] = blockNum1; txList2[12] = accounts[2]; txList2[12] = amount/2;
        let txHash2 = web3.sha3(RLP.encode(txList2).toString('hex'), {encoding: 'hex'});

        // include this tx the next block
        let sigs2 = [toHex(await web3.eth.sign(accounts[1], txHash2)), toHex(Buffer.alloc(65).toString('hex'))];

        let txBytes2 = RLP.encode([txList2, sigs2]).toString('hex');

        let merkleHash2 = sha256String(txBytes2);
        let root2, proof2;
        [root2, proof2] = generateMerkleRootAndProof([merkleHash2], 0);
        let blockNum2 = (await rootchain.lastCommittedBlock.call()).toNumber() + 1;
        await rootchain.submitBlock(toHex(root2), [1], blockNum2, {from: authority});

        // create confirmation signature
        let confirmationHash2 = sha256String(merkleHash2.slice(2) + root2.slice(2));
        let confirmSigs2 = await web3.eth.sign(accounts[2], confirmationHash2);

        // accounts[1] exits the second output
        await rootchain.startTransactionExit([blockNum1, 0, 1], toHex(txBytes1),
            toHex(proof1), toHex(confirmSigs1), {from: accounts[1], value: minExitBond});

        // try to challenge with the spend of the first output
        let err;
        [err] = await catchError(rootchain.challengeTransactionExit([blockNum1, 0, 1], [blockNum2, 0, 0],
            toHex(txBytes2), toHex(proof2), toHex(confirmSigs2)))
        if (!err)
            assert.fail("Challenged with incorrect transaction")
    });

    it("Attempt a withdrawal delay attack", async () => {
        let five_days = 432000 // in seconds
        // accounts[1] spends deposit and creates 2 new utxos for themself
        let txList1 = Array(17).fill(0);
        txList1[0] = txPos[0]; txList1[1] = txPos[1]; txList1[2] = txPos[2]; // first input
        txList1[12] = accounts[1]; txList1[13] = amount/2; // first output
        txList1[14] = accounts[1]; txList1[15] = amount/2; // second output
        let txHash1 = web3.sha3(RLP.encode(txList1).toString('hex'), {encoding: 'hex'});
        let sigs1 = [toHex(await web3.eth.sign(accounts[1], txHash1)), toHex(Buffer.alloc(65).toString('hex'))];
        let txBytes1 = RLP.encode([txList1, sigs1]).toString('hex');

        let merkleHash1 = sha256String(txBytes1);
        let root1, proof1;
        [root1, proof1] = generateMerkleRootAndProof([merkleHash1], 0);
        let blockNum1 = (await rootchain.lastCommittedBlock.call()).toNumber() + 1;
        await rootchain.submitBlock(toHex(root1), [1], blockNum1, {from: authority});

        // create confirmation signature
        let confirmationHash1 = sha256String(merkleHash1.slice(2) + root1.slice(2));
        let confirmSigs1 = await web3.eth.sign(accounts[1], confirmationHash1);

        // accounts[1] spends (blockNum1, 0, 1) utxo, sends 1 utxo to themself and the other to accounts[2]
        let txList2 = Array(17).fill(0);
        txList2[0] = blockNum1; txList2[2] = 1; // first input
        txList2[12] = accounts[1]; txList2[13] = amount / 4; // first output
        txList2[14] = accounts[2]; txList2[15] = amount / 4; // second output
        let txHash2 = web3.sha3(RLP.encode(txList2).toString('hex'), {encoding: 'hex'});
        let sigs2 = [toHex(await web3.eth.sign(accounts[1], txHash2)), toHex(Buffer.alloc(65).toString('hex'))];
        let txBytes2 = RLP.encode([txList2, sigs2]).toString('hex');

        let merkleHash2 = sha256String(txBytes2);
        let root2, proof2;
        [root2, proof2] = generateMerkleRootAndProof([merkleHash2], 0);
        let blockNum2 = (await rootchain.lastCommittedBlock.call()).toNumber() + 1;
        await rootchain.submitBlock(toHex(root2), [1], blockNum2, {from: authority});

        // create confirmation signature
        let confirmationHash2 = sha256String(merkleHash2.slice(2) + root2.slice(2));
        let confirmSigs2 = await web3.eth.sign(accounts[1], confirmationHash2);

        // make utxos > 1 week old
        fastForward(one_week + 100);

        // start exit for accounts[2], last utxo to be created
        await rootchain.startTransactionExit([blockNum2, 0, 1],
            toHex(txBytes2), toHex(proof2), toHex(confirmSigs2), {from: accounts[2], value: minExitBond});

        // increase time slightly, so exit by accounts[2] has better priority than accounts[1]
        fastForward(10);

        // start exit for accounts[1] utxo
        await rootchain.startTransactionExit([blockNum2, 0, 0],
            toHex(txBytes2), toHex(proof2), toHex(confirmSigs2), {from: accounts[1], value: minExitBond});

        // Fast Forward ~5 days
        fastForward(five_days);

        // Check to make sure challenge period has not ended
        let position = 1000000 * blockNum2 + 1;
        let currExit = await rootchain.txExits.call(position);
        assert.ok((currExit[2] + 604800) > (await web3.eth.getBlock(await web3.eth.blockNumber)).timestamp);

        // start exit for accounts[1], oldest utxo avaliable
        await rootchain.startTransactionExit([blockNum1, 0, 0],
            toHex(txBytes1), toHex(proof1), toHex(confirmSigs1), {from: accounts[1], value: minExitBond});

        // Fast Forward < 1 week
        fastForward(five_days);

        // finalize exits should finalize accounts[2] then accounts[1]
        let finalizedExits = await rootchain.finalizeTransactionExits({from: authority});
        let finalizedExit = await rootchain.txExits.call(position);
        assert.equal(finalizedExits.logs[0].args.position.toString(), [blockNum2, 0, 1, 0].toString(), "Incorrect position for finalized tx");
        assert.equal(finalizedExits.logs[0].args.owner, accounts[2], "Incorrect finalized exit owner");
        assert.equal(finalizedExits.logs[0].args.amount.toNumber(), 25 + minExitBond, "Incorrect finalized exit amount.");
        assert.equal(finalizedExit[3].toNumber(), 3, "Incorrect finalized exit state.");

        // Check other exits
        position = 1000000 * blockNum2;
        finalizedExit = await rootchain.txExits.call(position);
        assert.equal(finalizedExits.logs[2].args.position.toString(), [blockNum2, 0, 0, 0].toString(), "Incorrect position for finalized tx");
        assert.equal(finalizedExits.logs[2].args.owner, accounts[1], "Incorrect finalized exit owner");
        assert.equal(finalizedExits.logs[2].args.amount.toNumber(), 25 + minExitBond, "Incorrect finalized exit amount.");
        assert.equal(finalizedExit[3].toNumber(), 3, "Incorrect finalized exit state.");

        // Last exit should still be pending
        position = 1000000 * blockNum1;
        let pendingExit = await rootchain.txExits.call(position);
        assert.equal(pendingExit[2], accounts[1], "Incorrect pending exit owner");
        assert.equal(pendingExit[0], 50, "Incorrect pending exit amount");
        assert.equal(pendingExit[3].toNumber(), 1, "Incorrect pending exit state.");

        // Fast Forward rest of challenge period
        fastForward(one_week);
        await rootchain.finalizeTransactionExits({from: authority});
        // Check that last exit was processed
        finalizedExit = await rootchain.txExits.call(position);
        assert.equal(finalizedExit[2], accounts[1], "Incorrect finalized exit owner");
        assert.equal(finalizedExit[0], 50, "Incorrect finalized exit amount");
        assert.equal(finalizedExit[3].toNumber(), 3, "Incorrect finalized exit state.");
    });
});
