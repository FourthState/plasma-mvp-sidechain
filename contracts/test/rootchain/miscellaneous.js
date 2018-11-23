let RootChain = artifacts.require("RootChain");

let { catchError, toHex } = require("../utilities.js");
let { mineNBlocks } = require("./rootchain_helpers.js");

contract('[RootChain] Miscellaneous', async (accounts) => {

    let rootchain;
    let authority = accounts[0];
    beforeEach(async () => {
        rootchain = await RootChain.new({from: authority});
    });

    it("Will not revert finalizeExit with an empty queue", async () => {
        await rootchain.finalizeDepositExits();
        await rootchain.finalizeTransactionExits();
    });

    it("Can submit more than one merkle root", async () => {
        let root1 = web3.sha3("root1").slice(2);
        let root2 = web3.sha3("root2").slice(2);
        let roots = root1 + root2;

        // rootchain finality check
        mineNBlocks(6);

        let currentChildBlock = (await rootchain.currentChildBlock.call()).toNumber();
        rootchain.submitBlock(toHex(roots), {from: authority});

        assert.equal((await rootchain.currentChildBlock.call()).toNumber(), currentChildBlock + 2, "blocknum incremented incorrectly");
        assert.equal((await rootchain.getChildBlock.call(currentChildBlock))[0], toHex(root1), "mismatch in block root");
        assert.equal((await rootchain.getChildBlock.call(currentChildBlock+1))[0], toHex(root2), "mismatch in block root");
    });
});
