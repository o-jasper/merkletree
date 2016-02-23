local Statementize = require "merkle.statement.Statementize"
local sha2 = require "merkle.sha2"

local MerkleTree = require "merkle.Merkle.Tree"

return {
   -- Nonced versions.
   Sha224N = Statementize:class_derive(
      MerkleTree,
      {name = "MerkleSha224N", H=sha2.sha224}),
   Sha256N = Statementize:class_derive(
      MerkleTree,
      {name = "MerkleSha256N", H=sha2.sha256}),

   Sha224 = Statementize:class_derive(
      MerkleTree,
      {name = "MerkleSha224", H=sha2.sha224,
       gen_nonce = false, nonce_size = 0, always_nonce=false}),
   Sha256 = Statementize:class_derive(
      MerkleTree,
      {name = "MerkleSha256", H=sha2.sha256,
       gen_nonce = false, nonce_size = 0, always_nonce=false}),
}
