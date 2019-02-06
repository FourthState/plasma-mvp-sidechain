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
        let header = web3.utils.keccak256('1234');
        let tx = await instance.submitBlock([header], [1], [0], 1, {from: authority});

        // BlockSubmitted event
        assert.equal(tx.logs[0].args.header, header, "incorrect block header in BlockSubmitted event");
        assert.equal(tx.logs[0].args.blockNumber.toNumber(), 1, "incorrect block number in BlockSubmitted event");
        assert.equal(tx.logs[0].args.numTxns.toNumber(), 1, "incorrect block size in BlockSubmitted event");

        assert.equal((await instance.plasmaChain.call(1))[0], header, 'child block merkle header does not match submitted merkle header.');
    });

    it("Submit block from someone other than authority", async () => {
        let prev = (await instance.lastCommittedBlock.call()).toNumber();

        let [err] = await catchError(instance.submitBlock([web3.utils.keccak256('578484785954')], [1],
            [0], 1, {from: accounts[1]}));
        if (!err)
            assert.fail("Submitted blocked without being the authority");

        let curr = (await instance.lastCommittedBlock.call()).toNumber();
        assert.equal(prev, curr, "Child blocknum incorrectly changed");
    });

    it("Can submit more than one merkle header", async () => {
        let header1 = web3.utils.keccak256("header1").slice(2);
        let header2 = web3.utils.keccak256("header2").slice(2);

        await instance.submitBlock([toHex(header1), toHex(header2)], [1, 2], [0, 0], 1, {from: authority});

        assert.equal((await instance.lastCommittedBlock.call()).toNumber(), 2, "lastCommittedBlock not incremented incorrectly");
        assert.equal((await instance.plasmaChain.call(1))[0], toHex(header1), "mismatch in block header");
        assert.equal((await instance.plasmaChain.call(2))[0], toHex(header2), "mismatch in block header");
    });

    it("Enforces block number ordering", async () => {
        let header1 = web3.utils.keccak256("header1").slice(2)
        let header3 = web3.utils.keccak256("header3").slice(2)

        await instance.submitBlock([toHex(header1)], [1], [0], 1, {from: authority});
        let err;
        [err] = await catchError(instance.submitBlock([toHex(header3)], [1], 3, {from: authority}));
        if (!err)
            assert.fail("allowed block submission with inconsistent ordering");
    });

    it("Enforces block size capacity", async () => {
        let header = web3.utils.keccak256("header");
        let [err] = await catchError(instance.submitBlock([header], [Math.pow(2, 16)], [0], 1, {from: authority}))
        if (!err)
            assert.fail("block size capacity not enforeced");

        [err] = await catchError(instance.submitBlock([header], [0], [0], 1, {from: authority}))
        if (!err)
            assert.fail("block size capacity not enforced");
    });
});
