let assert = require('chai').assert;

let BytesUtil_Test = artifacts.require("BytesUtil_Test");

let { catchError, toHex } = require('../utilities.js');

contract('BytesUtil', async (accounts) => {
    let instance;
    before(async () => {
        instance = await BytesUtil_Test.new();
    });

    it("Slices bytes correctly", async () => {
        let inputHash = web3.sha3("inputSeed");

        assert.equal((await instance.slice.call(toHex(inputHash), 0, 32)).toString(), inputHash, "Slice didn't get entire substring");

        assert.equal((await instance.slice.call(toHex(inputHash), 0, 16)).toString(), toHex(inputHash.substring(2,34)), "Didn't get first half of the hash");
        assert.equal((await instance.slice.call(toHex(inputHash), 16, 16)).toString(), toHex(inputHash.substring(34)), "Didn't get second half of the hash");

        assert.equal((await instance.slice.call(toHex(inputHash), 0, 8)).toString(), toHex(inputHash.substring(2,18)), "Didn't get first quarter of the hash");
        assert.equal((await instance.slice.call(toHex(inputHash), 8, 24)).toString(), toHex(inputHash.substring(18)), "Didn't get rest of the hash");
    });

    it("Reverts if trying to slice out of range", async () => {
        let inputHash = web3.sha3("inputSeed");

        // sha3 maps input to a 32 byte hash (64 charac
        let err;
        [err] = await catchError(instance.slice.call(toHex(inputHash), 1, 32));
        if (!err)
            assert.fail("slice did not revert when inputs produce an out of bounds error");
    });

    it("Can slice bytes larger than a evm word size", async () => {
        let input = "0x";
        for (let i = 0; i < 100; i++) { // 50 bytes
            input += Math.floor(Math.random()*10) // include a random hex digit from 0-9
        }

        assert.equal((await instance.slice.call(toHex(input), 1, 40)).toString(), toHex(input.substring(4, 84)), "Didn't copy over a whole word size then left over bytes");
    });
});
