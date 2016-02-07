local sha2 = require 'merkle.sha2'
local enhex = require "merkle.enhex"

local fd = io.open("/dev/stdin")
local str = fd:read("l") or ""
fd:close()

-- Check consistency with class.
if math.random(4) == 1 then
   assert( sha2.hash256(str) == sha2.Hash256:new{str}:close() )
else
   local i = math.random(#str)
   assert( sha2.hash256(str) ==
              sha2.Hash256:new{string.sub(str, 1, i), string.sub(str, i+1)}:close() )
end
print(enhex(sha2.hash256(str)))
