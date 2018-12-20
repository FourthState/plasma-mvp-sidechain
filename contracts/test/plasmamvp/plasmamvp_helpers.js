let RLP = require('rlp');
let ethjs_util = require('ethereumjs-util');

let { toHex } = require('../utilities.js');

// Fast forward 1 week
let fastForward = async function(time) {
    await web3.currentProvider.send({jsonrpc: "2.0", method: "evm_mine", params: [], id: 0});
    let oldTime = (await web3.eth.getBlock(await web3.eth.blockNumber)).timestamp;

    // fast forward
    await web3.currentProvider.send({jsonrpc: "2.0", method: "evm_increaseTime", params: [time], id: 0});

    await web3.currentProvider.send({jsonrpc: "2.0", method: "evm_mine", params: [], id: 0});
    let currTime = (await web3.eth.getBlock(await web3.eth.blockNumber)).timestamp;

    assert.isAtLeast(currTime - oldTime, time, `Block time was not fast forwarded by at least ${time} seconds`);
}

// SHA256 hash the input and returns it in string form.
// Expects a hex input.
let sha256String = function(input) {
    return toHex(ethjs_util.sha256(toHex(input)).toString('hex'));
};

// SHA256 hashes together 2 inputs and returns it in string form.
// Expects hex inputs, and prepend each input with a 0x20 byte literal.
// Tendermint prefixes intermediate hashes with 0x20 bytes literals 
// before hashing them.
let sha256StringMultiple = function(input1, input2) {
    let toHash = "0x20" + input1.slice(2) + "20" + input2.slice(2);
    return toHex(ethjs_util.sha256(toHash).toString('hex'));
};

// For a given list of leaves, this function constructs a simple merkle tree.
// It returns the merkle root and the merkle proof for the txn at index.
// @param leaves The leaves for which this function generates a merkle root and proof
// @param txIndex The leaf for which this function generates a merkle proof
//
// Simple Tree: https://tendermint.com/docs/spec/blockchain/encoding.html#merkle-trees
let generateMerkleRootAndProof = function(leaves, index) {
    if (leaves.length == 0) { // If there are no leaves, then we can't generate anything
        return ["", ""];
    } else if (leaves.length == 1) { // If there's only 1 leaf, return it with and empty proof
        return [leaves[0], ""];
    } else {
        let pivot = Math.floor((leaves.length + 1) / 2);

        let left, right;
        let proof = "";

        // If the index will be in the left subtree (index < pivot), then we
        // need to generate the proof using the intermediary hash from the right
        // side. Otherwise, do the reverse.
        if (index < pivot) {
            // recursively call the function on the leaves that will be in the
            // left and right sub trees.
            left = generateMerkleRootAndProof(leaves.slice(0, pivot), index);
            right = generateMerkleRootAndProof(leaves.slice(pivot, leaves.length), -1);

            // add current level's right intermediary hash to the proof
            if (index >= 0) {
                proof = left[1] + right[0].slice(2);
            }
        } else {
            // recursively call the function on the leaves that will be in the
            // left and right sub trees.
            // since the index will be in the right sub tree, we need to update
            // it's value.
            left = generateMerkleRootAndProof(leaves.slice(0, pivot), -1);
            right = generateMerkleRootAndProof(leaves.slice(pivot, leaves.length), index - pivot);

            // add current level's left intermediary hash to the proof
            if (index >= 0) {
                proof = right[1] + left[0].slice(2);
            }
        }
        return [sha256StringMultiple(left[0], right[0]), toHex(proof)];
    }
};


module.exports = {
    fastForward,
    sha256String,
    generateMerkleRootAndProof
};
