
local H = require("merkle.sha2").hash224

local Tree = require "merkle.Merkle.Tree"
local verify = require("merkle.Merkle.Verify"):new{ H = H }

local function test_example()
   local tree = Tree:new{ H = H }

   tree:add("one")
   local provable = tree:add("two", true)
   tree:add("three")
   tree:add("four")

   local root_H = tree:finish()

   local proof = provable:produce_proof()

   assert( verify:verify(root_H, proof, "two") )
   assert( not verify:verify(root_H, proof, "wrong") )
end

test_example()

local function rand_str(lmin, lmax)
   local i, n, ret = 0, math.random(lmin or 4, lmax or 10), {}
   while i < n do
      table.insert(ret, string.char(math.random(33, 126)))
      i = i + 1
   end
   return table.concat(ret)
end

local function test_blast(n, prove_p)
   local tree = Tree:new{ H = H }

   local provables = {}
   while n > 0 do
      local data = rand_str()
      if math.random() < prove_p then
         local node = tree:add(data, true)
         node.data = data
         assert(node.H == H(data))
         table.insert(provables, node)
      else
         tree:add(data)
      end
      n = n - 1
   end
   local root_H = tree:finish()

   for i, node in ipairs(provables) do
      local proof = node:produce_proof()
      assert(verify:verify(root_H, proof, node.data), string.format("Failed %d", i))
   end
end

local m = 10
while m > 0 do
   print(m)
   test_blast(19, 0.1)
   m = m - 1
end
