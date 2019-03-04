let RLP = require('rlp');
let assert = require('chai').assert;

let PlasmaMVP = artifacts.require("PlasmaMVP");

let { fastForward, proof, zeroHashes, sha256String, generateMerkleRootAndProof, fillTxList } = require('./plasmamvp_helpers.js');
let { catchError, toHex } = require('../utilities.js');

contract('[PlasmaMVP] Deposits', async (accounts) => {
    let instance;
    let oneWeek = 604800; // in seconds
    let minExitBond = 200000;

    let authority = accounts[0];
    before(async () => {
        instance = await PlasmaMVP.deployed();
    });

    it("Will not revert finalizeDeposit on an empty queue", async () => {
        await instance.finalizeDepositExits();
    });

    it("Catches Deposit event", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        let tx = await instance.deposit(accounts[1], {from: accounts[1], value: 100});
        // check Deposit event
        assert.equal(tx.logs[0].args.depositor, accounts[1], "incorrect deposit owner");
        assert.equal(tx.logs[0].args.amount.toNumber(), 100, "incorrect deposit amount");
        assert.equal(tx.logs[0].args.depositNonce, nonce, "incorrect deposit nonce");
    });

    it("Allows deposits of funds into a different address", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        let tx = await instance.deposit(accounts[2], {from: accounts[1], value: 100});

        // check Deposit event
        assert.equal(tx.logs[0].args.depositor, accounts[2], "incorrect deposit owner");
        assert.equal(tx.logs[0].args.amount.toNumber(), 100, "incorrect deposit amount");
        assert.equal(tx.logs[0].args.depositNonce, nonce, "incorrect deposit nonce");

        // check instance deposit mapping
        let deposit = await instance.deposits.call(nonce);
        assert.equal(deposit[0], accounts[2], "incorrect deposit owner");
        assert.equal(deposit[1], 100, "incorrect deposit amount");
    });

    it("Rejects a committed fee larger than the deposit amount", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[1], value: 100});

        let err;
        [err] = await catchError(instance.startDepositExit(nonce, 100, {from: accounts[1], value: minExitBond}))
        if (!err)
            assert.fail("started a deposit exit with a committed fee equal to the amount");
    });

    it("Only allows deposit owner to start a deposit exit", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[1], value: 100});

        let err;
        // accounts[1] cannot start exit because it's not the owner
        [err] = await catchError(instance.startDepositExit(nonce, 0, {from: accounts[1], value: minExitBond}));
        if (!err)
            assert.fail("Non deposit owner allowed to start an exit");

        // accounts[2] should be able to start exit
        await instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond});
    });

    it("Rejects exiting a deposit with a malicious committed fee", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 10});

        let err;
        [err] = await catchError(instance.startDepositExit(nonce, 11, {from: accounts[2], value: minExitBond}));
        if (!err)
            assert.fail("exited a deposit with a committed fee larger than the deposit amount");
    });

    it("Rejects exiting a deposit twice", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});
        await instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond});

        let err;
        [err] = await catchError(instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond}));
        if (!err)
            assert.fail("Started an exit for the same deposit twice.");
    });

    it("Catches StartedDepositExit event", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});
        let tx = await instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond});

        assert.equal(tx.logs[0].args.nonce.toNumber(), nonce, "StartedDepositExit event emits incorrect nonce");
        assert.equal(tx.logs[0].args.owner, accounts[2], "StartedDepositExit event emits incorrect owner");
        assert.equal(tx.logs[0].args.amount.toNumber(), 100, "StartedDepositExit event emits incorrect amount");
    });

    it("Requires sufficient bond and refunds excess if overpayed", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});

        let err;
        [err] = await catchError(instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond-10}));
        if (!err)
            assert.fail("started exit with insufficient bond");

        await instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond+10});

        let balance = (await instance.balanceOf.call(accounts[2])).toNumber();
        assert.equal(balance, 10, "excess for overpayed bond not refunded to sender");
    });

    it("Can start and finalize a deposit exit", async () => {
        instance = await PlasmaMVP.new({from: authority});

        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});

        let plasmaChainBalance = (await instance.plasmaChainBalance.call()).toNumber();
        assert.equal(plasmaChainBalance, 100);

        await instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond});
        await fastForward(oneWeek + 100);

        await instance.finalizeDepositExits();

        plasmaChainBalance = (await instance.plasmaChainBalance.call()).toNumber();
        assert.equal(plasmaChainBalance, 0);

        let balance = (await instance.balanceOf.call(accounts[2])).toNumber();
        assert.equal(balance, 100 + minExitBond, "deposit exit not finalized after a week");

        let exit = await instance.depositExits.call(nonce);
        assert.equal(exit[4], 3, "exit's state not set to finalized");
    });

    it("Cannot reopen a finalized deposit exit", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});
        await instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond});

        await fastForward(oneWeek + 100);

        await instance.finalizeDepositExits();
        let err;
        [err] = await catchError(instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond}));
        if (!err)
            assert.fail("reopened a finalized deposit exit");
    });

    it("Correctly challenge a spent deposit", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});

        // construct transcation with first input as the deposit
        let txList = Array(15).fill(0);
        txList[3] = nonce; txList[10] = accounts[1]; txList[11] = 100;
        txList = fillTxList(txList);
        let txHash = web3.utils.soliditySha3(toHex(RLP.encode(txList).toString('hex')));
        let sigs = [await web3.eth.sign(txHash, accounts[2]), toHex(Buffer.alloc(65).toString('hex'))];
        let txBytes = [txList, sigs];
        txBytes = RLP.encode(txBytes).toString('hex');

        // include this transaction in the next block
        let root;
        [root, proof] = generateMerkleRootAndProof([sha256String(txBytes)], 0);

        let blockNum = (await instance.lastCommittedBlock.call()).toNumber() + 1;
        await instance.submitBlock([toHex(root)], [1], [0], blockNum, {from: authority});

        // create the confirm sig
        let confirmHash = sha256String(sha256String(txBytes) + root.slice(2));
        let confirmSig = await web3.eth.sign(confirmHash, accounts[2]);

        // start the malicious exit
        await instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond});

        // checks matching inputs
        let err;
        [err] = await catchError(instance.challengeExit([0,0,0,nonce-1], [blockNum, 0],
            toHex(txBytes), toHex(sigs), toHex(proof), toHex(confirmSig), {from: accounts[3]}));
        if (!err)
            assert.fail("did not check against matching inputs");

        // invalid confirm signature
        [err] = await catchError(instance.challengeExit([0,0,0,nonce-1], [blockNum, 0],
            toHex(txBytes), toHex(sigs), toHex(proof), toHex(confirmSig.substring(0, 30)), {from: accounts[3]}));
        if (!err)
            assert.fail("challenged with an invalid signature");

        // correctly challenge
        await instance.challengeExit([0,0,0,nonce], [blockNum, 0],
            toHex(txBytes), toHex(proof), toHex(confirmSig), {from: accounts[3]});

        let balance = (await instance.balanceOf.call(accounts[3])).toNumber();
        assert.equal(balance, minExitBond, "challenger not awarded exit bond");

        let exit = await instance.depositExits.call(nonce);
        assert.equal(exit[4], 2, "exit state not changed to challenged");

        // Cannot challenge twice
        [err] = await catchError(instance.challengeExit([0,0,0,nonce], [blockNum, 0],
            toHex(txBytes), toHex(sigs), toHex(proof), toHex(confirmSig), {from: accounts[3]}));
        if (!err)
            assert.fail("Allowed a challenge for an exit already challenged");
    });

    it("Challenge an invalid first input exit fee mismatch and exit the fee", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});
        let nonce2 = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});

        let txList = Array(15).fill(0);
        txList[3] = nonce; txList[14] = 5; // fee
        txList[8] = nonce2; // second input
        txList[10] = accounts[1]; txList[11] = 100;
        txList = fillTxList(txList);
        let txHash = web3.utils.soliditySha3(toHex(RLP.encode(txList).toString('hex')));
        let sigs = [toHex(await web3.eth.sign(txHash, accounts[2])), toHex(Buffer.alloc(65).toString('hex'))];
        let txBytes = [txList, sigs];
        txBytes = RLP.encode(txBytes).toString('hex');

        // submit the block
        let [merkleRoot, proof] = generateMerkleRootAndProof([sha256String(txBytes)], 0);
        let blockNum = (await instance.lastCommittedBlock.call()).toNumber() + 1;
        await instance.submitBlock([toHex(merkleRoot)], [1], [5], blockNum, {from: authority});

        // exit the first deposit input. commit incorrect fee
        await instance.startDepositExit(nonce, 1, {from: accounts[2], value: minExitBond});

        // exit the second deposit input
        await instance.startDepositExit(nonce2, 0, {from: accounts[2], value: minExitBond});

        // challenge the first input with a fee mismatch
        let tx = await instance.challengeExit([0,0,0,nonce], [blockNum, 0], toHex(txBytes), toHex(proof), toHex(""));

        let exit = await instance.depositExits.call(nonce);
        assert.equal(exit[4].toNumber(), 0, "exit state not changed to non existent");

        // second input should not be able to be challenged with a fee mismatch
        let err;
        [err] = await catchError(instance.challengeExit([0,0,0,nonce2], [blockNum, 0], toHex(txBytes), toHex(proof), toHex("")));
        if (!err)
            assert.fail("challenged second input with a fee mismatch");

        // start the exit again with the correct committed fee
        await instance.startDepositExit(nonce, 5, {from: accounts[2], value: minExitBond});

        // try challenge the exit with a fee mismatch
        [err] = await catchError(instance.challengeExit([0,0,0,nonce], [blockNum, 0], toHex(txBytes), toHex(proof), toHex("")));
        if (!err)
            assert.fail("operator challenged exit with correct committed fee");

        // start a fee exit
        await instance.startFeeExit(blockNum, 0, {from: authority, value: minExitBond});
    });

    it("Attempts a withdrawal delay attack on exiting deposits", async () => {

        /* 1. Start exit for nonce_2 (newest) 
         * 2. Start exit for nonce_1, nonce_0 in the same eth block
         * 3. Check exit ordering:  nonce_2, nonce_0, nonce_1
         */

        instance = await PlasmaMVP.new({from: authority});

        let nonce_0 = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});

        let nonce_1 = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});

        let nonce_2 = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});

        // exit nonce_2 
        await instance.startDepositExit(nonce_2, 0, {from: accounts[2], value: minExitBond});
        await fastForward(10);

        // exit nonce_1 then nonce_0 in the same ethereum block
        await instance.startDepositExit(nonce_0, 0, {from: accounts[2], value: minExitBond});
        await instance.startDepositExit(nonce_1, 0, {from: accounts[2], value: minExitBond});

        await fastForward(oneWeek + 10);
        let depositExits = await instance.finalizeDepositExits({from: authority});
        // every even index to skip the `AddedToBalances` event
        assert.equal(depositExits.logs[0].args.position.toString(), [0, 0, 0, nonce_2].toString(), "nonce_2 was not finalized first");
        assert.equal(depositExits.logs[2].args.position.toString(), [0, 0, 0, nonce_0].toString(), "nonce_0 was not finalized second");
        assert.equal(depositExits.logs[4].args.position.toString(), [0, 0, 0, nonce_1].toString(), "nonce_1 was not finalized last");
    });
});
