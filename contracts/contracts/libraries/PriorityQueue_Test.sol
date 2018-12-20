pragma solidity ^0.4.24;

import "./PriorityQueue.sol";

// Purpose of this contract is to forward calls to the library for testing
contract PriorityQueue_Test {
    using PriorityQueue for uint256[];

    uint256[] heapList;

    function insert(uint256 k) public { heapList.insert(k); }
    function getMin() public view returns (uint256) { return heapList.getMin(); }
    function delMin() public { heapList.delMin(); }
    function currentSize() public view returns (uint256) { return heapList.length; }
}
