let RLP = require('rlp');
let assert = require('chai').assert;

let RootChain = artifacts.require("RootChain");

let { fastForward, mineNBlocks, proof, zeroHashes } = require('./rootchain_helpers.js');
let { catchError, toHex } = require('../utilities.js');

contract('[RootChain] Deposits', async (accounts) => {
    let rootchain;
    let one_week = 604800; // in seconds
    let minExitBond = 10000;

    let authority = accounts[0];
    beforeEach(async () => {
        rootchain = await RootChain.new({from: authority});
    });

    it("Catches Deposit event", async () => {
        let nonce = (await rootchain.depositNonce.call()).toNumber();
        let tx = await rootchain.deposit(accounts[1], {from: accounts[1], value: 100});
        // check Deposit event
        assert.equal(tx.logs[0].args.depositor, accounts[1], "incorrect deposit owner");
        assert.equal(tx.logs[0].args.amount.toNumber(), 100, "incorrect deposit amount");
        assert.equal(tx.logs[0].args.depositNonce, nonce, "incorrect deposit nonce");
    });

    it("Allows deposits of funds into a different address", async () => {
        let nonce = (await rootchain.depositNonce.call()).toNumber();
        let tx = await rootchain.deposit(accounts[2], {from: accounts[1], value: 100});

        // check Deposit event
        assert.equal(tx.logs[0].args.depositor, accounts[2], "incorrect deposit owner");
        assert.equal(tx.logs[0].args.amount.toNumber(), 100, "incorrect deposit amount");
        assert.equal(tx.logs[0].args.depositNonce, nonce, "incorrect deposit nonce");

        // check rootchain deposit mapping
        let deposit = await rootchain.getDeposit.call(nonce);
        assert.equal(deposit[0], accounts[2], "incorrect deposit owner");
        assert.equal(deposit[1], 100, "incorrect deposit amount");
    });

    it("Only allows deposit owner to start a deposit exit", async () => {
        let nonce = (await rootchain.depositNonce.call()).toNumber();
        await rootchain.deposit(accounts[2], {from: accounts[1], value: 100});
        let err;

        // accounts[1] cannot start exit because it's not the owner
        [err] = await catchError(rootchain.startDepositExit(nonce, {from: accounts[1], value: minExitBond}));
        if (!err)
            assert.fail("Non deposit owner allowed to start an exit");

        //accounts[2] should be able to start exit
        await rootchain.startDepositExit(nonce, {from: accounts[2], value: minExitBond});
    });

    it("Rejects exiting a deposit twice", async () => {
        let nonce = (await rootchain.depositNonce.call()).toNumber();
        await rootchain.deposit(accounts[2], {from: accounts[2], value: 100});
        await rootchain.startDepositExit(nonce, {from: accounts[2], value: minExitBond});

        let err;
        [err] = await catchError(rootchain.startDepositExit(nonce, {from: accounts[2], value: minExitBond}));
        if (!err)
            assert.fail("Started an exit for the same deposit twice.");
    });

    it("Catches StartedDepositExit event", async () => {
        let nonce = (await rootchain.depositNonce.call()).toNumber();
        await rootchain.deposit(accounts[2], {from: accounts[2], value: 100});
        let tx = await rootchain.startDepositExit(nonce, {from: accounts[2], value: minExitBond});

        assert.equal(tx.logs[0].args.nonce.toNumber(), nonce, "StartedDepositExit event emits incorrect nonce");
        assert.equal(tx.logs[0].args.owner, accounts[2], "StartedDepositExit event emits incorrect owner");
        assert.equal(tx.logs[0].args.amount.toNumber(), 100, "StartedDepositExit event emits incorrect amount");
    });

    it("Requires sufficient bond and refunds excess if overpayed", async () => {
        let nonce = (await rootchain.depositNonce.call()).toNumber();
        await rootchain.deposit(accounts[2], {from: accounts[2], value: 100});

        let err;
        [err] = await catchError(rootchain.startDepositExit(nonce, {from: accounts[2], value: minExitBond-10}));
        if (!err)
            assert.fail("started exit with insufficient bond");

        await rootchain.startDepositExit(nonce, {from: accounts[2], value: minExitBond+10});

        let balance = (await rootchain.balanceOf.call(accounts[2])).toNumber();
        assert.equal(balance, 10, "excess for overpayed bond not refunded to sender");
    });

    it("Can start and finalize a deposit exit. Child chain balance should reflect accordingly", async () => {
        let nonce = (await rootchain.depositNonce.call()).toNumber();
        await rootchain.deposit(accounts[2], {from: accounts[2], value: 100});

        let childChainBalance = (await rootchain.childChainBalance.call()).toNumber();
        assert.equal(childChainBalance, 100);

        await rootchain.startDepositExit(nonce, {from: accounts[2], value: minExitBond});
        await fastForward(one_week + 100);
        await rootchain.finalizeDepositExits();

        childChainBalance = (await rootchain.childChainBalance.call()).toNumber();
        assert.equal(childChainBalance, 0);

        let balance = (await rootchain.balanceOf.call(accounts[2])).toNumber();
        assert.equal(balance, 100 + minExitBond, "deposit exit not finalized after a week");

        let exit = await rootchain.getDepositExit.call(nonce);
        assert.equal(exit[3], 3, "exit's state not set to finalized");
    });

    it("Cannot reopen a finalized deposit exit", async () => {
        let nonce = (await rootchain.depositNonce.call()).toNumber();
        await rootchain.deposit(accounts[2], {from: accounts[2], value: 100});
        await rootchain.startDepositExit(nonce, {from: accounts[2], value: minExitBond});

        await fastForward(one_week + 100);

        await rootchain.finalizeDepositExits();
        let err;
        [err] = await catchError(rootchain.startDepositExit(nonce, {from: accounts[2], value: minExitBond}));
        if (!err)
            assert.fail("reopened a finalized deposit exit");
    });

    it("Correctly challenge a spent deposit", async () => {
        let nonce = (await rootchain.depositNonce.call()).toNumber();
        await rootchain.deposit(accounts[2], {from: accounts[2], value: 100});

        // construct transcation with first input as the deposit
        let msg = Array(17).fill(0);
        msg[3] = nonce; msg[12] = accounts[1]; msg[13] = 100;
        let encodedMsg = RLP.encode(msg);
        let hashedEncodedMsg = web3.sha3(encodedMsg.toString('hex'), {encoding: 'hex'});

        // create signature by deposit owner. Second signature should be zero
        let sigList = Array(2).fill(0);
        sigList[0] = (await web3.eth.sign(accounts[2], hashedEncodedMsg));

        let txBytes = Array(2).fill(0);
        txBytes[0] = msg; txBytes[1] = sigList;
        txBytes = RLP.encode(txBytes);

        // create signature by deposit owner. Second signature should be zero
        let sigs = (await web3.eth.sign(accounts[2], hashedEncodedMsg));
        sigs = sigs + Buffer.alloc(65).toString('hex');

        let merkleHash = web3.sha3(txBytes.toString('hex'), {encoding: 'hex'});

        // include this transaction in the next block
        let root = merkleHash;
        for (let i = 0; i < 16; i++)
            root = web3.sha3(root + zeroHashes[i], {encoding: 'hex'}).slice(2)
        let blockNum = (await rootchain.currentChildBlock.call()).toNumber();
        mineNBlocks(5); // presumed finality before submitting the block
        await rootchain.submitBlock(toHex(root), {from: authority});

        // create the confirm sig
        let confirmHash = web3.sha3(merkleHash.slice(2) + root, {encoding: 'hex'});
        let confirmSig = await web3.eth.sign(accounts[2], confirmHash);

        // start the malicious exit
        await rootchain.startDepositExit(nonce, {from: accounts[2], value: minExitBond});

        // checks matching inputs
        let err;
        [err] = await catchError(rootchain.challengeDepositExit(nonce-1, [blockNum, 0, 0],
            toHex(txBytes), toHex(sigs), toHex(proof), toHex(confirmSig), {from: accounts[3]}));
        if (!err)
            assert.fail("did not check against matching inputs");

        // correctly challenge
        await rootchain.challengeDepositExit(nonce, [blockNum, 0, 0],
            toHex(txBytes), toHex(proof), toHex(confirmSig), {from: accounts[3]});

        let balance = (await rootchain.balanceOf.call(accounts[3])).toNumber();
        assert.equal(balance, minExitBond, "challenger not awarded exit bond");

        let exit = await rootchain.getDepositExit.call(nonce);
        assert.equal(exit[3], 2, "exit state not changed to challenged");

        // Cannot challenge twice
        [err] = await catchError(rootchain.challengeDepositExit(nonce, [blockNum, 0, 0],
            toHex(txBytes), toHex(sigs), toHex(proof), toHex(confirmSig), {from: accounts[3]}));
        if (!err)
            assert.fail("Allowed a challenge for an exit already challenged");
    });
});
