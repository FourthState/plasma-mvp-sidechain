pragma solidity ^0.4.24;

import "./BytesUtil.sol";

// from https://tendermint.com/docs/spec/blockchain/encoding.html#merkle-trees
library TMSimpleMerkleTree {
    using BytesUtil for bytes;

    // @param leaf     a leaf of the tree
    // @param index    position of this leaf in the tree that is zero indexed
    // @param rootHash block header of the merkle tree
    // @param proof    sequence of 32-byte hashes from the leaf up to, but excluding, the root
    // @paramt total   total # of leafs in the tree
    function checkMembership(bytes32 leaf, uint256 index, bytes32 rootHash, bytes proof, uint256 total)
        internal
        pure
        returns (bool)
    {
        // variable size Merkle tree, but proof must consist of 32-byte hashes
        require(proof.length % 32 == 0, "Incorrect proof length");

        bytes32 computedHash = computeHashFromAunts(index, total, leaf, proof);
        return computedHash == rootHash;
    }

    // helper function as described in the tendermint docs
    function computeHashFromAunts(uint256 index, uint256 total, bytes32 leaf, bytes innerHashes)
        private
        pure
        returns (bytes32)
    {
        require(index < total, "Index must be less than total number of leaf nodes");
        require(total > 0, "Must have at least one leaf node");

        if (total == 1) {
            require(innerHashes.length == 0, "Simple Tree with 1 txn should have no innerHashes");
            return leaf;
        }
        require(innerHashes.length != 0, "Simple Tree with > 1 txn should have innerHashes");

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

        if (index < numLeft) {
            bytes32 leftHash = computeHashFromAunts(index, numLeft, leaf, innerHashes.slice( 0, innerHashes.length - 32));
            uint innerHashesMemOffset = innerHashes.length - 32;
            assembly {
                // get the last 32-byte hash from innerHashes array
                proofElement := mload(add(add(innerHashes, 0x20), innerHashesMemOffset))
            }

            return sha256(abi.encodePacked(b, leftHash, b, proofElement));
        } else {
            bytes32 rightHash = computeHashFromAunts(index-numLeft, total-numLeft, leaf, innerHashes.slice(0, innerHashes.length - 32));
            innerHashesMemOffset = innerHashes.length - 32;
            assembly {
                    // get the last 32-byte hash from innerHashes array
                    proofElement := mload(add(add(innerHashes, 0x20), innerHashesMemOffset))
            }
            return sha256(abi.encodePacked(b, proofElement, b, rightHash));
        }
    }
}
