pragma solidity ^0.5.0;

import "./TMSimpleMerkleTree.sol";

contract TMSimpleMerkleTree_Test {
    using TMSimpleMerkleTree for bytes32;

    function checkMembership(bytes32 leaf, uint256 index, bytes32 rootHash, bytes memory proof, uint256 total)
        public
        pure
        returns (bool)
    {
        return leaf.checkMembership(index, rootHash, proof, total);
    }
}
