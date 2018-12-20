pragma solidity ^0.4.24;

import "./Validator.sol";

/*
* Used to proxy function calls to the Validator for testing
*/

contract Validator_Test {
<<<<<<< HEAD

  using Validator for bytes32;

  function checkMembership(bytes32 leaf, uint256 index, bytes32 rootHash, bytes proof, uint256 total)
      public
      pure
      returns (bool)
  {
      return leaf.checkMembership(index, rootHash, proof, total);
  }

  function checkSigs(bytes32 txHash, bytes32 confirmationHash, bool input1, bytes sig0, bytes sig1, bytes confirmSignatures)
=======
  using Validator for bytes32;

  function checkSignatures(bytes32 txHash, bytes32 confirmationHash, bool input1, bytes sig0, bytes sig1, bytes confirmSignatures)
>>>>>>> b3167013cb609ec55bd2a944e44a4d169ed332c9
      public
      pure
      returns (bool)
  {
<<<<<<< HEAD
      return txHash.checkSigs(confirmationHash, input1, sig0, sig1, confirmSignatures);
=======
      return txHash.checkSignatures(confirmationHash, input1, sig0, sig1, confirmSignatures);
>>>>>>> b3167013cb609ec55bd2a944e44a4d169ed332c9
  }

  function recover(bytes32 hash, bytes sig)
      public
      pure
      returns (address)
  {
      return hash.recover(sig);
  }
<<<<<<< HEAD

  function slice(bytes _bytes, uint start, uint len)
      public
      pure
      returns (bytes)
  {
      return Validator.slice(_bytes, start, len);
  }
=======
>>>>>>> b3167013cb609ec55bd2a944e44a4d169ed332c9
}
