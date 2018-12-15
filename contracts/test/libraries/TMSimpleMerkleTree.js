let assert = require('chai').assert;

let TMSimpleMerkleTree_Test = artifacts.require("TMSimpleMerkleTree_Test");

let { catchError, toHex } = require('../utilities.js');
let { generateMerkleRootAndProof } = require('../plasmamvp/plasmamvp_helpers.js');

contract('TMSimpleMerkleTree', async (accounts) => {
    let instance;
    before(async () => {
        instance = await TMSimpleMerkleTree_Test.new();
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
});
