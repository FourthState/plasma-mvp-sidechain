pragma solidity ^0.5.0;

import "./MinPriorityQueue.sol";

// Purpose of this contract is to forward calls to the library for testing
contract MinPriorityQueue_Test {
    using MinPriorityQueue for uint256[];

    uint256[] heapList;

    function insert(uint256 k) public { heapList.insert(k); }
    function delMin() public { heapList.delMin(); }
    function currentSize() public view returns (uint256) { return heapList.length; }
    function getMin() public view returns (uint256) {
        require(heapList.length != 0, "empty queue");
        return heapList[0]; 
    }
}
