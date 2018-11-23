let RLP = require('rlp');

let { toHex } = require('../utilities.js');

// Wait for n blocks to pass
let mineNBlocks = async function(numBlocks) {
    for (i = 0; i < numBlocks; i++) {
    await web3.currentProvider.send({jsonrpc: "2.0", method: "evm_mine", params: [], id: 0});
    }
}

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

// For a given list of leaves, this function generates a merkle root. It assumes
// the merkle tree is of depth 16. If there are less than 2^16 leaves, the
// list is padded with 0x0 transactions. The function also generates a merkle
// proof for the leaf at txIndex.
// @param leaves The leaves for which this function generates a merkle root and proof
// @param txIndex The leaf for which this function gneerates a merkle proof
let generateMerkleRootAndProof = function(leaves, txIndex) {
    return generateMerkleRootAndProofHelper(leaves, 16, txIndex, 0);
};

// This helper function recursively generates a merkle root and merkle proof for
// a given list of leaves and a leaf's txIndex.
let generateMerkleRootAndProofHelper = function(leaves, depth, txIndex, zeroHashesIndex) {
    // If the depth is 0, then we are already at the root. This means that we
    // expect there to only be one leaf, which is the root.
    if (depth == 0) {
        if (leaves.length == 1) {
            return [leaves[0], ""];
        }
        else {
            return ["", ""];
        }
    }
    else {
        let newLeaves = [];
        let proof = "";

        // For each pair of leaves, concat them together and hash the result
        let i = 0;
        while (i + 2 <= leaves.length) {
            let mergedHash = web3.sha3(leaves[i].slice(2) + leaves[i + 1].slice(2), {encoding: 'hex'});
            newLeaves.push(mergedHash);

            // For the txIndex of interest, we want to generate a merkle proof,
            // which means that we need to keep track of the other leaf in the
            // pair.
            if (txIndex == i) {
                proof = leaves[i + 1].slice(2);
            }
            else if (txIndex == i + 1) {
                proof = leaves[i].slice(2);
            }

            i += 2;
        }

        // If i < leaves.length, then that means there's an odd number of leaves
        // In this case, we need to hash the remaining leaf with the zeroHash of
        // the current depth, which has been hardcoded in "rootchain_helpers"
        if (i < leaves.length) {
            let mergedHash = web3.sha3(leaves[i].slice(2) + zeroHashes[zeroHashesIndex], {encoding: 'hex'});
            // For the txIndex of interest, we want to generate a merkle proof,
            // which means that we need to keep track of the other leaf in the
            // pair.
            if (txIndex == i) {
                proof = zeroHashes[zeroHashesIndex];
            }
            newLeaves.push(mergedHash);
        }

        // Recursively call the helper function, updating the variables we pass in
        // We expect to see the number of leaves to decrease by 1/2
        // This would be the next layer up in the merkle tree.
        let result = generateMerkleRootAndProofHelper(newLeaves, depth - 1, Math.floor(txIndex/2), zeroHashesIndex + 1);

        result[1] = proof + result[1];

        return result;
    }
};


// 512 bytes
let proof = '0000000000000000000000000000000000000000000000000000000000000000ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5b4c11951957c6f8f642c4af61cd6b24640fec6dc7fc607ee8206a99e92410d3021ddb9a356815c3fac1026b6dec5df3124afbadb485c9ba5a3e3398a04b7ba85e58769b32a1beaf1ea27375a44095a0d1fb664ce2dd358e7fcbfb78c26a193440eb01ebfc9ed27500cd4dfc979272d1f0913cc9f66540d7e8005811109e1cf2d887c22bd8750d34016ac3c66b5ff102dacdd73f6b014e710b51e8022af9a1968ffd70157e48063fc33c97a050f7f640233bf646cc98d9524c6b92bcf3ab56f839867cc5f7f196b93bae1e27e6320742445d290f2263827498b54fec539f756afcefad4e508c098b9a7e1d8feb19955fb02ba9675585078710969d3440f5054e0f9dc3e7fe016e050eff260334f18a5d4fe391d82092319f5964f2e2eb7c1c3a5f8b13a49e282f609c317a833fb8d976d11517c571d1221a265d25af778ecf8923490c6ceeb450aecdc82e28293031d10c7d73bf85e57bf041a97360aa2c5d99cc1df82d9c4b87413eae2ef048f94b4d3554cea73d92b0f7af96e0271c691e2bb5c67add7c6caf302256adedf7ab114da0acfe870d449a3a489f781d659e8beccda7bce9f4e8618b6bd2f4132ce798cdc7a60e7e1460a7299e3c6342a579626d2';

let zeroHashes = [ '0000000000000000000000000000000000000000000000000000000000000000',
    'ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5',
    'b4c11951957c6f8f642c4af61cd6b24640fec6dc7fc607ee8206a99e92410d30',
    '21ddb9a356815c3fac1026b6dec5df3124afbadb485c9ba5a3e3398a04b7ba85',
    'e58769b32a1beaf1ea27375a44095a0d1fb664ce2dd358e7fcbfb78c26a19344',
    '0eb01ebfc9ed27500cd4dfc979272d1f0913cc9f66540d7e8005811109e1cf2d',
    '887c22bd8750d34016ac3c66b5ff102dacdd73f6b014e710b51e8022af9a1968',
    'ffd70157e48063fc33c97a050f7f640233bf646cc98d9524c6b92bcf3ab56f83',
    '9867cc5f7f196b93bae1e27e6320742445d290f2263827498b54fec539f756af',
    'cefad4e508c098b9a7e1d8feb19955fb02ba9675585078710969d3440f5054e0',
    'f9dc3e7fe016e050eff260334f18a5d4fe391d82092319f5964f2e2eb7c1c3a5',
    'f8b13a49e282f609c317a833fb8d976d11517c571d1221a265d25af778ecf892',
    '3490c6ceeb450aecdc82e28293031d10c7d73bf85e57bf041a97360aa2c5d99c',
    'c1df82d9c4b87413eae2ef048f94b4d3554cea73d92b0f7af96e0271c691e2bb',
    '5c67add7c6caf302256adedf7ab114da0acfe870d449a3a489f781d659e8becc',
    'da7bce9f4e8618b6bd2f4132ce798cdc7a60e7e1460a7299e3c6342a579626d2' ];

module.exports = {
    fastForward,
    mineNBlocks,
    proof,
    zeroHashes,
    generateMerkleRootAndProof
};
