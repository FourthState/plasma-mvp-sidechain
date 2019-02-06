pragma solidity ^0.5.0;

import "./BytesUtil.sol";

contract BytesUtil_Test {
    using BytesUtil for bytes;

    function slice(bytes memory a, uint start, uint len) public pure returns (bytes memory) { return a.slice(start, len); }
}
