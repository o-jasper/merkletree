local Statementize = require "merkle.statement.Statementize"
local sha2 = require "merkle.sha2"

return {
   -- Nonced versions.
   Sha224N = Statementize:class_derive{name = "Sha224N", H = sha2.sha224 },
   Sha256N = Statementize:class_derive{name = "Sha256N", H = sha2.sha256 },

   Sha224 = Statementize:class_derive{name = "Sha224", H = sha2.sha224,
                                      gen_nonce = false, nonce_size = 0, always_nonce=false},
   Sha256 = Statementize:class_derive{name = "Sha256", H = sha2.sha256,
                                      gen_nonce = false, nonce_size = 0, always_nonce=false},
}
