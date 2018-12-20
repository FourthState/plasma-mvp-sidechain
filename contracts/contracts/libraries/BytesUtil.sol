pragma solidity ^0.4.24;

library BytesUtil {
    uint8 constant WORD_SIZE = 32;

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
