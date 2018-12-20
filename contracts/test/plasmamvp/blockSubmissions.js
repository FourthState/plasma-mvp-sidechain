let assert = require('chai').assert;

let PlasmaMVP = artifacts.require("PlasmaMVP");

let { toHex, catchError } = require('../utilities.js');

contract('[PlasmaMVP] Block Submissions', async (accounts) => {
    let instance;
    let authority = accounts[0];
    let minExitBond = 10000;
    beforeEach(async () => {
        instance = await PlasmaMVP.new({from: authority});
    });

    it("Submit block from authority", async () => {
        let root = web3.sha3('1234');
        let tx = await instance.submitBlock([root], [1], [0], 1, {from: authority});

        // BlockSubmitted event
        assert.equal(tx.logs[0].args.root, root, "incorrect block root in BlockSubmitted event");
        assert.equal(tx.logs[0].args.blockNumber.toNumber(), 1, "incorrect block number in BlockSubmitted event");
        assert.equal(tx.logs[0].args.numTxns.toNumber(), 1, "incorrect block size in BlockSubmitted event");
        assert.equal(tx.logs[0].args.feeAmount.toNumber(), 0, "incorrect block fee amount in BlockSubmitted event");

        assert.equal((await instance.childChain.call(1))[0], root, 'Child block merkle root does not match submitted merkle root.');
    });

    it("Submit block from someone other than authority", async () => {
        let prev = (await instance.lastCommittedBlock.call()).toNumber();

        let [err] = await catchError(instance.submitBlock([web3.sha3('578484785954')], [1], [0], 1, {from: accounts[1]}));
        if (!err)
            assert.fail("Submitted blocked without being the authority");

        let curr = (await instance.lastCommittedBlock.call()).toNumber();
        assert.equal(prev, curr, "Child blocknum incorrectly changed");
    });

    it("Can submit more than one merkle root", async () => {
        let root1 = web3.sha3("root1").slice(2);
        let root2 = web3.sha3("root2").slice(2);

        await instance.submitBlock([toHex(root1), toHex(root2)], [1, 2], [0, 0], 1, {from: authority});

        assert.equal((await instance.lastCommittedBlock.call()).toNumber(), 2, "blocknum incremented incorrectly");
        assert.equal((await instance.childChain.call(1))[0], toHex(root1), "mismatch in block root");
        assert.equal((await instance.childChain.call(2))[0], toHex(root2), "mismatch in block root");
    });

    it("Can submit fees for a block", async () => {
        let root1 = web3.sha3("root1").slice(2);
        let root2 = web3.sha3("root2").slice(2);

        let fees = [100, 200];

        await instance.submitBlock([toHex(root1), toHex(root2)], [1, 2], fees, 1, {from: authority});

        assert.equal((await instance.lastCommittedBlock.call()).toNumber(), 2, "blocknum incremented incorrectly");
        assert.equal((await instance.childChain.call(1))[0], toHex(root1), "mismatch in block root");
        assert.equal((await instance.childChain.call(2))[0], toHex(root2), "mismatch in block root");
        assert.equal((await instance.childChain.call(1))[2], fees[0], "mismatch in block fees");
        assert.equal((await instance.childChain.call(2))[2], fees[1], "mismatch in block fees");
    });

    it("Cannot exceed size limits for a block", async () => {
        let root = web3.sha3("root1").slice(2);

        let fees = [100];
        let numTxns0 = [65535]
        let numTxns1 = [65536];

        let err;
        // numTxns1 exceeds block size limit
        [err] = await catchError(instance.submitBlock([toHex(root)], numTxns1, fees, 1, {from: authority}));
        if (!err)
            assert.fail("Allowed submission of a block that exceeds block size limit");

        // numTxns0 does not exceed block size limit
        await instance.submitBlock([toHex(root)], numTxns0, fees, 1, {from: authority});
    });

    it("Enforces block number ordering", async () => {
        let root1 = web3.sha3("root1").slice(2)
        let root3 = web3.sha3("root3").slice(2)

        await instance.submitBlock([toHex(root1)], [1], [0], 1);
        let err;
        [err] = await catchError(instance.submitBlock([toHex(root3)], [1], [0],3));
        if (!err)
            assert.fail("Allowed block submission with inconsistent ordering");
    });
});
