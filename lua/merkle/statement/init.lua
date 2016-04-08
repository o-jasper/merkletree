local Statementize = require "merkle.statement.Statementize"
local sha2 = require "merkle.sha2"

return {
   __constant = true,

   -- Nonced versions.
   Sha224N = Statementize:class_derive(sha2.Sha224, {name = "Sha224N"}),
   Sha256N = Statementize:class_derive(sha2.Sha256, {name = "Sha256N"}),

   -- No nonce. (never_nonce off; _may_ check nonced data.
   Sha224 = Statementize:class_derive(
      sha2.Sha224,
      {name = "Sha224", gen_nonce = false, nonce_size = 0, always_nonce=false}),
   Sha256 = Statementize:class_derive(
      sha2.Sha256,
      {name = "Sha256", gen_nonce = false, nonce_size = 0, always_nonce=false}),
}
