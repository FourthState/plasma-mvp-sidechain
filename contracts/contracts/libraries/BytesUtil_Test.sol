pragma solidity ^0.4.24;

import "./BytesUtil.sol";

contract BytesUtil_Test {
    using BytesUtil for bytes;

    function slice(bytes a , uint start, uint len) public pure returns (bytes) { return a.slice(start, len); }
}
