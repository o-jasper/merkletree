local StatementizeMerkle = require "merkle.statement.Statementize"
local sha2 = require "merkle.sha2"

return {
   -- Nonced versions.
   Sha224N = StatementizeMerkle:class_derive{name = "MerkleSha224N", H = sha2.sha224 },
   Sha256N = StatementizeMerkle:class_derive{name = "MerkleSha256N", H = sha2.sha256 },

   Sha224 = StatementizeMerkle:class_derive{name = "MerkleSha224", H = sha2.sha224,
                        gen_nonce = false, nonce_size = 0, always_nonce=false},
   Sha256 = StatementizeMerkle:class_derive{name = "MerkleSha256", H = sha2.sha256,
                        gen_nonce = false, nonce_size = 0, always_nonce=false},
}
