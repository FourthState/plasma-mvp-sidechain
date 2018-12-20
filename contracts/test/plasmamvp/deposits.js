let RLP = require('rlp');
let assert = require('chai').assert;

let PlasmaMVP = artifacts.require("PlasmaMVP");

let { fastForward, proof, zeroHashes, sha256String, generateMerkleRootAndProof } = require('./plasmamvp_helpers.js');
let { catchError, toHex } = require('../utilities.js');

contract('[PlasmaMVP] Deposits', async (accounts) => {
    let instance;
    let one_week = 604800; // in seconds
    let minExitBond = 10000;

    let authority = accounts[0];
    beforeEach(async () => {
        instance = await PlasmaMVP.new({from: authority});
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

    it("Only allows deposit owner to start a deposit exit", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[1], value: 100});
        let err;

        // accounts[1] cannot start exit because it's not the owner
        [err] = await catchError(instance.startDepositExit(nonce, 0, {from: accounts[1], value: minExitBond}));
        if (!err)
            assert.fail("Non deposit owner allowed to start an exit");

        //accounts[2] should be able to start exit
        await instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond});
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

    it("Can start and finalize a deposit exit. Child chain balance should reflect accordingly", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});

        let childChainBalance = (await instance.childChainBalance.call()).toNumber();
        assert.equal(childChainBalance, 100);

        await instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond});
        await fastForward(one_week + 100);
        await instance.finalizeDepositExits();

        childChainBalance = (await instance.childChainBalance.call()).toNumber();
        assert.equal(childChainBalance, 0);

        let balance = (await instance.balanceOf.call(accounts[2])).toNumber();
        assert.equal(balance, 100 + minExitBond, "deposit exit not finalized after a week");

        let exit = await instance.depositExits.call(nonce);
        assert.equal(exit[4], 3, "exit's state not set to finalized");
    });

    it("Cannot reopen a finalized deposit exit", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});
        await instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond});

        await fastForward(one_week + 100);

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

        let merkleHash = sha256String(txBytes.toString('hex'));

        // include this transaction in the next block
        let root;
        [root, proof] = generateMerkleRootAndProof([merkleHash], 0);

        let blockNum = (await instance.lastCommittedBlock.call()).toNumber() + 1;
        await instance.submitBlock([toHex(root)], [1], [0], blockNum, {from: authority});

        // create the confirm sig
        let confirmHash = sha256String(merkleHash + root.slice(2));
        let confirmSig = await web3.eth.sign(accounts[2], confirmHash);

        // start the malicious exit
        await instance.startDepositExit(nonce, 0, {from: accounts[2], value: minExitBond});

        // checks matching inputs
        let err;
        [err] = await catchError(instance.challengeExit([0,0,0,nonce-1], [blockNum, 0],
            toHex(txBytes), toHex(sigs), toHex(proof), toHex(confirmSig), {from: accounts[3]}));
        if (!err)
            assert.fail("did not check against matching inputs");

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

    it("Allows operator to challenge a deposit spend committing to an incorrect fee", async () => {
        let nonce = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});

        let txList = Array(17).fill(0);
        txList[3] = nonce; txList[16] = 5; // fee
        txList[12] = accounts[1]; txList[13] = 100;
        let txHash = web3.sha3(RLP.encode(txList).toString('hex'), {encoding: 'hex'});
        let sigs = [toHex(await web3.eth.sign(accounts[2], txHash)), toHex(Buffer.alloc(65).toString('hex'))];

        let txBytes = [txList, sigs];
        txBytes = RLP.encode(txBytes).toString('hex');

        // submit the block
        let merkleHash = sha256String(txBytes);
        let [merkleRoot, proof] = generateMerkleRootAndProof([merkleHash], 0);
        let blockNum = (await instance.lastCommittedBlock.call()).toNumber() + 1;
        await instance.submitBlock([toHex(merkleRoot)], [1], [5], blockNum, {from: authority});

        // exit the deposit not committing to the fee
        await instance.startDepositExit(nonce, 1, {from: accounts[2], value: minExitBond});

        // challenge the exit
        await instance.challengeFeeMismatch([0,0,0,nonce], [blockNum, 0], toHex(txBytes), toHex(proof));

        let exit = await instance.depositExits.call(nonce);
        assert.equal(exit[4].toNumber(), 0, "exit state not changed to non existent");

        // start the exit again with the correct committed fee
        await instance.startDepositExit(nonce, 5, {from: accounts[2], value: minExitBond});

        // try challenge the exit
        let err;
        [err] = await catchError(instance.challengeFeeMismatch([0,0,0,nonce], [blockNum, 0], toHex(txBytes), toHex(proof)));
        if (!err)
            assert.fail("operator challenged exit with correct committed fee");
    });

    it("Attempts a withdrawal delay attack on exiting deposits", async () => {
        
        /* 1. Start exit for nonce_2 (newest) 
         * 2. Start exit for nonce_1, nonce_0 in the same eth block
         * 3. Check exit ordering:  nonce_2, nonce_0, nonce_1
         */

        let nonce_0 = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});

        let nonce_1 = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});

        let nonce_2 = (await instance.depositNonce.call()).toNumber();
        await instance.deposit(accounts[2], {from: accounts[2], value: 100});

        // exit nonce_2 
        await instance.startDepositExit(nonce_2, 0, {from: accounts[2], value: minExitBond});
        
        // first exit should be in a different eth block
        fastForward(10);
        
        // exit nonce_1 then nonce_0 in the same ethereum block
        async function exits(nonce_1, nonce_0) {
            let p1 = instance.startDepositExit(nonce_1, 0, {from: accounts[2], value: minExitBond});
            let p2 = instance.startDepositExit(nonce_0, 0, {from: accounts[2], value: minExitBond});
            return Promise.all([p1, p2]);
        }
        await exits(nonce_1, nonce_0);
       
        fastForward(one_week);
        let depositExits = await instance.finalizeDepositExits({from: authority});
        assert.equal(depositExits.logs[0].args.position.toString(), [0, 0, 0, nonce_2].toString(), "nonce_2 was not finalized first");
        
        assert.equal(depositExits.logs[2].args.position.toString(), [0, 0, 0, nonce_0].toString(), "nonce_0 was not finalized second");
        assert.equal(depositExits.logs[4].args.position.toString(), [0, 0, 0, nonce_1].toString(), "nonce_1 was not finalized last");
    });
});
