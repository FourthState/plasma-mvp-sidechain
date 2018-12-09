pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/ECRecovery.sol";

library Validator {
    uint8 constant WORD_SIZE = 32;

    // @param leaf     a leaf of the tree
    // @param index    position of this leaf in the tree that is zero indexed
    // @param rootHash block header of the merkle tree
    // @param proof    sequence of 32-byte hashes from the leaf up to, but excluding, the root
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

    // from https://tendermint.com/docs/spec/blockchain/encoding.html#merkle-trees
    function computeHashFromAunts(uint256 index, uint256 total, bytes32 leaf, bytes innerHashes)
        internal
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
            bytes32 leftHash = computeHashFromAunts(index, numLeft, leaf, slice(innerHashes, 0, innerHashes.length - 32));
            uint innerHashesMemOffset = innerHashes.length - 32;
            assembly {
                // get the last 32-byte hash from innerHashes array
                proofElement := mload(add(add(innerHashes, 0x20), innerHashesMemOffset))
            }

            return sha256(abi.encodePacked(b, leftHash, b, proofElement));
        } else {
            bytes32 rightHash = computeHashFromAunts(index-numLeft, total-numLeft, leaf, slice(innerHashes, 0, innerHashes.length - 32));
            innerHashesMemOffset = innerHashes.length - 32;
            assembly {
                    // get the last 32-byte hash from innerHashes array
                    proofElement := mload(add(add(innerHashes, 0x20), innerHashesMemOffset))
            }
            return sha256(abi.encodePacked(b, proofElement, b, rightHash));
        }
    }

    // @param txHash      transaction hash
    // @param rootHash    block header of the merkle tree
    // @param input1      indicator for the second input
    // @param sigs        transaction signatures
    // @notice            when one input is present, we require it to be the first input by convention
    function checkSigs(bytes32 txHash, bytes32 confirmationHash, bool input1, bytes sig0, bytes sig1, bytes confirmSignatures)
        internal
        pure
        returns (bool)
    {
        require(sig0.length == 65 && sig1.length == 65, "signatures must be 65 bytes in length");

        if (input1) {
            require(confirmSignatures.length == 130, "two confirm signatures required with two inputs");

            address recoveredAddr0 = recover(txHash, sig0);
            address recoveredAddr1 = recover(txHash, sig1);

            return recoveredAddr0 == recover(confirmationHash, slice(confirmSignatures, 0, 65))
                   && recoveredAddr1 == recover(confirmationHash, slice(confirmSignatures, 65, 65))
                   && recoveredAddr0 != address(0) && recoveredAddr1 != address(0);
        }

        // only 1 input present
        require(confirmSignatures.length == 65, "one confirm signature required with one input present");

        address recoveredAddr = recover(txHash, sig0);
        return recoveredAddr == recover(confirmationHash, confirmSignatures) && recoveredAddr != address(0);
    }

    function recover(bytes32 hash, bytes sig)
        internal
        pure
        returns (address)
    {

        hash = ECRecovery.toEthSignedMessageHash(hash);
        return ECRecovery.recover(hash, sig);
    }

    /* Helpers */

    // @param _bytes raw bytes that needs to be slices
    // @param start  start of the slice relative to `_bytes`
    // @param len    length of the sliced byte array
    function slice(bytes _bytes, uint start, uint len)
            internal
            pure
            returns (bytes)
        {
            require(_bytes.length - start >= len, "slice out of bounds");

            if (_bytes.length == len)
                return _bytes;

            bytes memory result;
            uint src;
            uint dest;
            assembly {
                // memory & free memory pointer
                result := mload(0x40)
                mstore(result, len) // store the size in the prefix
                mstore(0x40, add(result, and(add(add(0x20, len), 0x1f), not(0x1f)))) // padding

                // pointers
                src := add(start, add(0x20, _bytes))
                dest := add(0x20, result)
            }

            // copy as many word sizes as possible
            for(; len >= WORD_SIZE; len -= WORD_SIZE) {
                assembly {
                    mstore(dest, mload(src))
                }

                src += WORD_SIZE;
                dest += WORD_SIZE;
            }

            // copy remaining bytes
            uint mask = 256 ** (WORD_SIZE - len) - 1;
            assembly {
                let srcpart := and(mload(src), not(mask)) // zero out src
                let destpart := and(mload(dest), mask) // retrieve the bytes
                mstore(dest, or(destpart, srcpart))
            }

            return result;
    }
}
