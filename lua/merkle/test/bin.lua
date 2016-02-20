local sha2 = require 'merkle.sha2'
local enhex = require "merkle.enhex"

local hash, Hash = sha2["sha" .. arg[1]], sha2["Sha" .. arg[1]]

local fd = io.open("/dev/stdin")
local str = fd:read("l") or ""
fd:close()

local fh = hash(str)

-- Check consistency with class.
if math.random(4) == 1 then
   assert( fh == Hash:new{str}:close() )
else
   local i = math.random(#str)
   assert( fh ==
           Hash:new{string.sub(str, 1, i), string.sub(str, i+1)}:close() )
end
print(enhex(fh))
