pragma solidity ^0.5.0;

import "./BytesUtil.sol";

// from https://tendermint.com/docs/spec/blockchain/encoding.html#merkle-trees
library TMSimpleMerkleTree {
    using BytesUtil for bytes;

    // @param leaf     a leaf of the tree
    // @param index    position of this leaf in the tree that is zero indexed
    // @param rootHash block header of the merkle tree
    // @param proof    sequence of 32-byte hashes from the leaf up to, but excluding, the root
    // @paramt total   total # of leafs in the tree
    function checkMembership(bytes32 leaf, uint256 index, bytes32 rootHash, bytes memory proof, uint256 total)
        internal
        pure
        returns (bool)
    {
        // variable size Merkle tree, but proof must consist of 32-byte hashes
        require(proof.length % 32 == 0); // incorrect proof length

        bytes32 computedHash = computeHashFromAunts(index, total, leaf, proof);
        return computedHash == rootHash;
    }

    // helper function as described in the tendermint docs
    function computeHashFromAunts(uint256 index, uint256 total, bytes32 leaf, bytes memory innerHashes)
        private
        pure
        returns (bytes32)
    {
        require(index < total); // index must be within bound of the # of leave
        require(total > 0); // must have one leaf node

        if (total == 1) {
            require(innerHashes.length == 0); // 1 txn has no proof
            return leaf;
        }
        require(innerHashes.length != 0); // >1 txns should have a proof

        uint256 numLeft = (total + 1) / 2;
        bytes32 proofElement;

        // prepend 0x20 byte literal to hashes
        // tendermint prefixes intermediate hashes with 0x20 bytes literals
        // before hashing them.
        bytes memory b = new bytes(1);
        assembly {
            let memPtr := add(b, 0x20)
            mstore8(memPtr, 0x20)
        }

        uint innerHashesMemOffset = innerHashes.length - 32;
        if (index < numLeft) {
            bytes32 leftHash = computeHashFromAunts(index, numLeft, leaf, innerHashes.slice(0, innerHashes.length - 32));
            assembly {
                // get the last 32-byte hash from innerHashes array
                proofElement := mload(add(add(innerHashes, 0x20), innerHashesMemOffset))
            }

            return sha256(abi.encodePacked(b, leftHash, b, proofElement));
        } else {
            bytes32 rightHash = computeHashFromAunts(index-numLeft, total-numLeft, leaf, innerHashes.slice(0, innerHashes.length - 32));
            assembly {
                    // get the last 32-byte hash from innerHashes array
                    proofElement := mload(add(add(innerHashes, 0x20), innerHashesMemOffset))
            }
            return sha256(abi.encodePacked(b, proofElement, b, rightHash));
        }
    }
}
