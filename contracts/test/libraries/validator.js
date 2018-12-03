let assert = require('chai').assert;
let RLP = require('rlp');

let Validator_Test = artifacts.require("Validator_Test");
let { catchError, toHex } = require('../utilities.js');
let { generateMerkleRootAndProof } = require('../rootchain/rootchain_helpers.js');

contract('Validator', async (accounts) => {
    let instance;
    beforeEach (async () => {
        instance = await Validator_Test.new();
    });

    it("Verifies the membership in a merkle tree with 7 leaves (hardcoded)", async () => {
        // this utxo input and the merkle root of its block were generated
        // by the side chain
        // this side chain block contains 7 txns
        let rootHash = "0x4edd08572735720e2769fa536c03e5d2ed29ff2ba97fb102cfe71fdae3c30428";
        let leaf = "0x714c0269d202e4302fadab4d62a4e9171fbf09508cb589616342cce45981d329";
        let proof = "0xfbec635c930057fdc76939052216ed2aed7af618109b983ee1ae7b13c909f2dd1a61753dc9eccd5506144081ba000a8e44c9262be17ab934dda9e4fa10495fccdbb108450dad46789d67f5cfb5ee4be7f505bd8835a0867410412f22cfda8ad5";
        let total = 7; // Merkle tree contains 7 leaf nodes (transactions)
        let index = 2; // index of the leaf we want to prove

        // checking membership of 3rd leaf
        assert.isTrue(await instance.checkMembership.call(toHex(leaf), index, toHex(rootHash), toHex(proof), total), "Didn't prove membership");


        // checking membership of 4th leaf
        leaf = "0xfbec635c930057fdc76939052216ed2aed7af618109b983ee1ae7b13c909f2dd";
        proof = "0x714c0269d202e4302fadab4d62a4e9171fbf09508cb589616342cce45981d3291a61753dc9eccd5506144081ba000a8e44c9262be17ab934dda9e4fa10495fccdbb108450dad46789d67f5cfb5ee4be7f505bd8835a0867410412f22cfda8ad5";
        index = 3;
        assert.isTrue(await instance.checkMembership.call(toHex(leaf), index, toHex(rootHash), toHex(proof), total), "Didn't prove membership");
    });

    it("Verifies the membership in a merkle tree with only one transaction", async () => {
        let leafHash = web3.sha3("inputSeed");

        let root, proof;
        [root, proof] = generateMerkleRootAndProof([leafHash], 0);

        assert.isTrue(await instance.checkMembership.call(toHex(leafHash), 0, toHex(root), toHex(proof), 1), "Didn't prove membership");
    });

    it("Catches bad input on checkMembership", async () => {
        let leafHash1 = web3.sha3("inputSeed1");
        let leafHash2 = web3.sha3("inputSeed2");
        let leafHash3 = web3.sha3("inputSeed3");

        let root, proof;
        [root, proof] = generateMerkleRootAndProof([leafHash1, leafHash2, leafHash3], 0);

        let badLeafHash = web3.sha3("wrongInputSeed", {encoding: 'hex'});
        assert.isFalse(await instance.checkMembership.call(toHex(badLeafHash), 0, toHex(root), toHex(proof), 3), "Returned true on wrong leaf");

        assert.isFalse(await instance.checkMembership.call(toHex(leafHash1), 1, toHex(root), toHex(proof), 3), "Returned true on wrong index");

        let badRoot = web3.sha3("wrongRoot", {encoding: 'hex'});
        assert.isFalse(await instance.checkMembership.call(toHex(leafHash1), 0, toHex(badRoot), toHex(proof), 3), "Returned true on wrong root");

        let badProof = "a".repeat(proof.length - 2);
        assert.isFalse(await instance.checkMembership.call(toHex(leafHash1), 0, toHex(root), toHex(badProof), 3), "Returned true on wrong proof");

        let err;
        [err] = await catchError(instance.checkMembership.call(toHex(leafHash1), 0, toHex(root), toHex(proof + "0000"), 3));
        if (!err)
            assert.fail("Didn't revert on an proof with the bad size");
    });

    it("Verifies membership in a merkle tree with multiple transactions", async () => {
        let leafHash1 = web3.sha3("inputSeed1");
        let leafHash2 = web3.sha3("inputSeed2");
        let leafHash3 = web3.sha3("inputSeed3");
        let leafHash4 = web3.sha3("inputSeed4");
        let leafHash5 = web3.sha3("inputSeed5");

        let root, proof;
        [root, proof] = generateMerkleRootAndProof([leafHash1, leafHash2, leafHash3, leafHash4, leafHash5], 0);
        assert.isTrue(await instance.checkMembership.call(toHex(leafHash1), 0, toHex(root), toHex(proof), 5), "Didn't prove membership");

        [root, proof] = generateMerkleRootAndProof([leafHash1, leafHash2, leafHash3, leafHash4, leafHash5], 1);
        assert.isTrue(await instance.checkMembership.call(toHex(leafHash2), 1, toHex(root), toHex(proof), 5), "Didn't prove membership");

        [root, proof] = generateMerkleRootAndProof([leafHash1, leafHash2, leafHash3, leafHash4, leafHash5], 2);
        assert.isTrue(await instance.checkMembership.call(toHex(leafHash3), 2, toHex(root), toHex(proof), 5), "Didn't prove membership");

        [root, proof] = generateMerkleRootAndProof([leafHash1, leafHash2, leafHash3, leafHash4, leafHash5], 3);
        assert.isTrue(await instance.checkMembership.call(toHex(leafHash4), 3, toHex(root), toHex(proof), 5), "Didn't prove membership");

        [root, proof] = generateMerkleRootAndProof([leafHash1, leafHash2, leafHash3, leafHash4, leafHash5], 4);
        assert.isTrue(await instance.checkMembership.call(toHex(leafHash5), 4, toHex(root), toHex(proof), 5), "Didn't prove membership");
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

        // assert valid confirmSignatures will pass checkSigs
        assert.isTrue(await instance.checkSigs.call(txHash, toHex(confirmationHash), false, toHex(sig0), toHex(sig1), toHex(confirmSignatures)),
            "checkSigs should pass");

        // assert invalid confirmSignatures will not pass checkSigs
        assert.isFalse(await instance.checkSigs.call(txHash, toHex(confirmationHash), false, toHex(sig0), toHex(sig1), toHex(invalidConfirmSignatures)),
            "checkSigs should not pass given invalid confirmSignatures");
    });

    it("Correctly handles empty signatures", async () => {
        let singleEmptyConfirmSig = Buffer.alloc(65).toString('hex');
        let doubleEmptyConfirmSigs = Buffer.alloc(130).toString('hex');
        let emptySig0 = Buffer.alloc(65).toString('hex');
        let emptySig1 = Buffer.alloc(65).toString('hex');

        let txHash = web3.sha3(Buffer.alloc(65).toString('hex'), {encoding: 'hex'});
        let confirmationHash = web3.sha3(Buffer.alloc(65).toString('hex'), {encoding: 'hex'});

        assert.isFalse(await instance.checkSigs.call(txHash, toHex(confirmationHash), false, toHex(emptySig0), toHex(emptySig1), toHex(singleEmptyConfirmSig)),
            "checkSigs should not pass given empty tx sigs and confirm signatures");

        assert.isFalse(await instance.checkSigs.call(txHash, toHex(confirmationHash), true, toHex(emptySig0), toHex(emptySig1), toHex(doubleEmptyConfirmSigs)),
            "checkSigs should not pass given empty tx sigs and confirm signatures");
    });

    it("Checks incorrect signature lengths", async () => {
        let confirmSignatures = Buffer.alloc(65).toString('hex');
        let sig0 = Buffer.alloc(65).toString('hex');
        let emptySig1 = Buffer.alloc(65).toString('hex');

        let txHash = web3.sha3(Buffer.alloc(65).toString('hex'), {encoding: 'hex'});
        let confirmationHash = web3.sha3(Buffer.alloc(65).toString('hex'), {encoding: 'hex'});

        let err;
        [err] = await catchError(instance.checkSigs.call(txHash, toHex(confirmationHash), false, toHex(sig0 + "0000"), toHex(emptySig1), toHex(confirmSignatures)));
        if (!err)
            assert.fail("Didn't revert on signature of wrong size");

        [err] = await catchError(instance.checkSigs.call(txHash, toHex(confirmationHash), false, toHex(sig0), toHex(emptySig1), toHex(confirmSignatures + "0000")));
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

        assert.isTrue(await instance.checkSigs.call(txHash, toHex(confirmationHash), false, toHex(sigOverTxHash), toHex(sig1), toHex(confirmSignatures)),
            "checkSigs should pass");
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

        assert.isFalse(await instance.checkSigs.call(txHash, toHex(confirmationHash), true, toHex(sig0), toHex(sigOverTxHash), toHex(confirmSignatures)),
            "checkSigs should not pass given an empty first confirmsig and non-empty second confirmsig");
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

        assert.isFalse(await instance.checkSigs.call(txHash, toHex(confirmationHash), true, toHex(sig0), toHex(invalidSig), toHex(confirmSignatures)),
            "checkSigs should not pass given invalid transaction sigs");
        assert.isFalse(await instance.checkSigs.call(txHash, toHex(confirmationHash), true, toHex(sig0), toHex(validSig), toHex(invalidConfirmSignatures)),
            "checkSigs should not pass given invalid transaction sigs");
        assert.isTrue(await instance.checkSigs.call(txHash, toHex(confirmationHash), true, toHex(sig0), toHex(validSig), toHex(confirmSignatures)),
            "checkSigs should pass for valid transaction sigs");
    });
});
