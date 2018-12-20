let assert = require('chai').assert;

let Validator_Test = artifacts.require("Validator_Test");
let { catchError, toHex } = require('../utilities.js');
let { generateMerkleRootAndProof } = require('../plasmamvp/plasmamvp_helpers.js');

contract('Validator', async (accounts) => {
    let instance;
    before(async () => {
        instance = await Validator_Test.new();
    });

    it("Correctly recovers the signee of a signature", async () => {
        // create tx hash
        let txHash = web3.sha3("inputSeed");

        let signer1 = accounts[1];
        // create tx sigs
        let txSigs1 = await web3.eth.sign(signer1, txHash);

        let signer2 = accounts[2];
        // create tx sigs
        let txSigs2 = await web3.eth.sign(signer2, txHash);

        assert.equal((await instance.recover.call(txHash, txSigs1)).toString(), signer1, "Recovered incorrect address");
        assert.equal((await instance.recover.call(txHash, txSigs2)).toString(), signer2, "Recovered incorrect address");
        assert.notEqual((await instance.recover.call(txHash, txSigs1)).toString(), (await instance.recover.call(txHash, txSigs2)).toString(),
            "Recovered the same address");
    });

    it("Correctly checks signatures", async () => {
        let signer = accounts[5];
        let invalidSigner = accounts[6];

        let txHash = web3.sha3("tx bytes to be hashed");
        let sig0 = await web3.eth.sign(signer, txHash);
        let sig1 = Buffer.alloc(65).toString('hex');

        let confirmationHash = web3.sha3("merkle leaf hash concat with root hash");

        let confirmSignatures = await web3.eth.sign(signer, confirmationHash);

        let invalidConfirmSignatures = await web3.eth.sign(invalidSigner, confirmationHash);

        // assert valid confirmSignatures will pass checkSignatures
        assert.isTrue(await instance.checkSignatures.call(txHash, toHex(confirmationHash), false, toHex(sig0), toHex(sig1), toHex(confirmSignatures)),
            "checkSignatures should pass");

        // assert invalid confirmSignatures will not pass checkSignatures
        assert.isFalse(await instance.checkSignatures.call(txHash, toHex(confirmationHash), false, toHex(sig0), toHex(sig1), toHex(invalidConfirmSignatures)),
            "checkSignatures should not pass given invalid confirmSignatures");
    });

    it("Correctly handles empty signatures", async () => {
        let singleEmptyConfirmSig = Buffer.alloc(65).toString('hex');
        let doubleEmptyConfirmSigs = Buffer.alloc(130).toString('hex');
        let emptySig0 = Buffer.alloc(65).toString('hex');
        let emptySig1 = Buffer.alloc(65).toString('hex');

        let txHash = web3.sha3(Buffer.alloc(65).toString('hex'), {encoding: 'hex'});
        let confirmationHash = web3.sha3(Buffer.alloc(65).toString('hex'), {encoding: 'hex'});

        assert.isFalse(await instance.checkSignatures.call(txHash, toHex(confirmationHash), false, toHex(emptySig0), toHex(emptySig1), toHex(singleEmptyConfirmSig)),
            "checkSignatures should not pass given empty tx sigs and confirm signatures");

        assert.isFalse(await instance.checkSignatures.call(txHash, toHex(confirmationHash), true, toHex(emptySig0), toHex(emptySig1), toHex(doubleEmptyConfirmSigs)),
            "checkSignatures should not pass given empty tx sigs and confirm signatures");
    });

    it("Checks incorrect signature lengths", async () => {
        let confirmSignatures = Buffer.alloc(65).toString('hex');
        let sig0 = Buffer.alloc(65).toString('hex');
        let emptySig1 = Buffer.alloc(65).toString('hex');

        let txHash = web3.sha3(Buffer.alloc(65).toString('hex'), {encoding: 'hex'});
        let confirmationHash = web3.sha3(Buffer.alloc(65).toString('hex'), {encoding: 'hex'});

        let err;
        [err] = await catchError(instance.checkSignatures.call(txHash, toHex(confirmationHash), false, toHex(sig0 + "0000"), toHex(emptySig1), toHex(confirmSignatures)));
        if (!err)
            assert.fail("Didn't revert on signature of wrong size");

        [err] = await catchError(instance.checkSignatures.call(txHash, toHex(confirmationHash), false, toHex(sig0), toHex(emptySig1), toHex(confirmSignatures + "0000")));
        if (!err)
            assert.fail("Didn't revert on confirm signature of wrong size");
    });

    it("Allows for only the first signature to be present", async () => {
        // create txHash
        let txBytes = web3.sha3("inputSeed");
        let txHash = web3.sha3(txBytes.toString('hex'), {encoding: 'hex'});

        // create sigs
        let signer = accounts[4];
        let sigOverTxHash = await web3.eth.sign(signer, txHash);
        let sig1 = Buffer.alloc(65).toString('hex');

        // create confirmationHash
        let merkleHash = web3.sha3(txHash.slice(2) + sigOverTxHash.slice(2), {encoding: 'hex'});
        let rootHash = generateMerkleRootAndProof([merkleHash], 0)[0];
        let confirmationHash = web3.sha3(merkleHash.slice(2) + rootHash, {encoding: 'hex'});

        // create confirmSignatures
        let confirmSignatures = await web3.eth.sign(signer, confirmationHash);

        assert.isTrue(await instance.checkSignatures.call(txHash, toHex(confirmationHash), false, toHex(sigOverTxHash), toHex(sig1), toHex(confirmSignatures)),
            "checkSignatures should pass");
    });

    it("Asserts that the first input cannot be empty", async () => {
        // create txHash
        let txBytes = web3.sha3("inputSeed");
        let txHash = web3.sha3(txBytes.toString('hex'), {encoding: 'hex'});

        // create sigs
        let signer = accounts[4];
        let sig0 = Buffer.alloc(65).toString('hex');
        let sigOverTxHash = (await web3.eth.sign(signer, txHash)).slice(2);

        // create confirmationHash
        let merkleHash = web3.sha3(txHash.slice(2) + sigOverTxHash, {encoding: 'hex'});
        let rootHash = generateMerkleRootAndProof([merkleHash], 0)[0];
        let confirmationHash = web3.sha3(merkleHash.slice(2) + rootHash, {encoding: 'hex'});

        // create confirmSignatures
        let confirmSignatures = Buffer.alloc(65).toString('hex');
        confirmSignatures += (await web3.eth.sign(signer, confirmationHash)).slice(2);

        assert.isFalse(await instance.checkSignatures.call(txHash, toHex(confirmationHash), true, toHex(sig0), toHex(sigOverTxHash), toHex(confirmSignatures)),
            "checkSignatures should not pass given an empty first confirmsig and non-empty second confirmsig");
    });

    it("Handles incorrect transaction signatures", async () => {
        // create txHash
        let txBytes = web3.sha3("inputSeed");
        let txHash = web3.sha3(txBytes.toString('hex'), {encoding: 'hex'});

        // create sigs
        let signer0 = accounts[4];
        let signer1 = accounts[5];
        let invalidSigner = accounts[6];
        let invalidSigner2 = accounts[7];

        // second tx sig is invalid
        let sig0 = await web3.eth.sign(signer0, txHash);
        let validSig = await web3.eth.sign(signer1, txHash).slice(2);
        let invalidSig = await web3.eth.sign(invalidSigner, txHash).slice(2);

        // create confirmationHash
        let merkleHash = web3.sha3(txHash.slice(2) + validSig.slice(2), {encoding: 'hex'});
        let rootHash = generateMerkleRootAndProof([merkleHash], 0)[0];
        let confirmationHash = web3.sha3(merkleHash.slice(2) + rootHash, {encoding: 'hex'});
        // create confirmSignatures
        let confirmSignatures = await web3.eth.sign(signer0, confirmationHash);
        confirmSignatures += await web3.eth.sign(signer1, confirmationHash).slice(2);
        // create invalid confirmSignatures
        let invalidConfirmSignatures = await web3.eth.sign(invalidSigner, confirmationHash);
        invalidConfirmSignatures += await web3.eth.sign(invalidSigner2, confirmationHash).slice(2);

        assert.isFalse(await instance.checkSignatures.call(txHash, toHex(confirmationHash), true, toHex(sig0), toHex(invalidSig), toHex(confirmSignatures)),
            "checkSignatures should not pass given invalid transaction sigs");
        assert.isFalse(await instance.checkSignatures.call(txHash, toHex(confirmationHash), true, toHex(sig0), toHex(validSig), toHex(invalidConfirmSignatures)),
            "checkSignatures should not pass given invalid transaction sigs");
        assert.isTrue(await instance.checkSignatures.call(txHash, toHex(confirmationHash), true, toHex(sig0), toHex(validSig), toHex(confirmSignatures)),
            "checkSignatures should pass for valid transaction sigs");
    });
});
