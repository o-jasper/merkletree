local sha2 = require 'merkle.sha2'

local function gstr(n)
   local n = n or 300
   if io then
      local fd = io.open("/dev/random")
      local str = fd:read(math.random(n)) or ""
      fd:close()
      return str
   else
      if os then math.randomseed(os.time() + 100*os.clock()) end
      local str, len = "", math.random(n)
      while #str < len do
         str = str .. string.char(math.random(256) - 1)
      end
      return str
   end
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
   t1(sha2.Sha224, sha2.sha224, math.random(4, 40))
   t1(sha2.Sha256, sha2.sha256, math.random(4, 40))
   j = j - 1
end
