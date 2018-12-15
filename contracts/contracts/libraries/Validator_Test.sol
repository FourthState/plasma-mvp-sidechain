pragma solidity ^0.4.24;

import "./Validator.sol";

/*
* Used to proxy function calls to the Validator for testing
*/

contract Validator_Test {
  using Validator for bytes32;

  function checkSignatures(bytes32 txHash, bytes32 confirmationHash, bool input1, bytes sig0, bytes sig1, bytes confirmSignatures)
      public
      pure
      returns (bool)
  {
      return txHash.checkSignatures(confirmationHash, input1, sig0, sig1, confirmSignatures);
  }

  function recover(bytes32 hash, bytes sig)
      public
      pure
      returns (address)
  {
      return hash.recover(sig);
  }
}
