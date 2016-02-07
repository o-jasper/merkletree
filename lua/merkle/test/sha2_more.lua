local sha2 = require 'merkle.sha2'

assert(sha2.Hash224.fun == sha2.hash224)
assert(sha2.Hash224.__name == "Hash224", sha2.Hash224.__name)

assert(sha2.Hash256.fun == sha2.hash256)
assert(sha2.Hash256.__name == "Hash256")

local function gstr()
   local fd = io.open("/dev/random")
   local str = fd:read("l") or ""
   fd:close()
   return str
end

local function t1(Class, fun, n)
   local hasher = Class:new()
   local i, strlist = n or 10, {}
   while i > 0 do
      local str = gstr()
      table.insert(strlist, str)
      hasher:add(str)
      i = i - 1
   end
   local h1, h2 = hasher:close(), fun(table.concat(strlist))
   assert( h1 == h2, string.format("Not same(%s)\n%s\n%s\ndata\n%s", Class.__name,
                                   h1,h2, table.concat(strlist)))
end

local j = 10

while j > 0 do
   t1(sha2.Hash224, sha2.hash224, math.random(4,10))
   t1(sha2.Hash256, sha2.hash256, math.random(4,10))
   j = j - 1
end
