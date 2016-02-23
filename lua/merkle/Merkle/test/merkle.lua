
local H = require("merkle.sha2").sha224

local Tree = require "merkle.Merkle.Tree"
local verify = require("merkle.Merkle.Verify"):new{ H = H }

local function test_example()
   local tree = Tree:new{ H = H }

   tree:add("one")
   local provable = tree:add("two", true)
   tree:add("three")
   tree:add("four")

   local root_H = tree:close()

   local proof = provable:produce_proof()

   assert( verify:verify(root_H, proof, "two") )
   assert( not verify:verify(root_H, proof, "wrong") )
end

test_example()

local random_fd

local function rand_str(lmin, lmax)
   local n = math.random(lmin or 4, lmax or 10)
   if not arg[1] then
      local i, ret = 0, {}
      while i < n do
         table.insert(ret, string.char(math.random(33, 126)))
         i = i + 1
      end
      return table.concat(ret)
   else
      random_fd = random_fd or io.open(arg[1])
      return random_fd:read(n)
   end
end

local function change_random_bit(proof)
   local i = math.random(#proof)
   local j = math.random(#proof[i])

   -- Change one bit, now the proof should be invalid.
   local n, x = math.random(7), string.byte(proof[i], j)
   local mod = string.char((x + (2^n))%256)
   local new_val = string.sub(proof[i], 1, j-1) .. mod .. string.sub(proof[i], j + 1)
   assert( new_val ~= proof[i], string.format("Didnt change the proof(%d,%s)?\n%s\n%s",
                                              n,x, proof[i], new_val))
   proof[i] = new_val
end

-- TODO Proofs that _do_not_ work out.
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
   local root_H = tree:close()

   for i, node in ipairs(provables) do
      local proof = node:produce_proof()
      assert(verify:verify(root_H, proof, node.data), string.format("Failed %d", i))
      -- Modify a bit, and check it is wrong.
      change_random_bit(proof)
      assert(not verify:verify(root_H, proof, node.data),
             string.format("False positive(alter 1 bit) %d", i))

      -- Make up nonsense, 
      local total_nonsense = {}
      for _,el in ipairs(proof) do table.insert(total_nonsense, rand_str(#el, #el)) end
      assert(not verify:verify(root_H, total_nonsense, node.data),
             string.format("False positive(nonsense) %d", i))
   end
end

local m = 10
while m > 0 do
   print(m)
   test_blast(19, 0.1)
   m = m - 1
end
