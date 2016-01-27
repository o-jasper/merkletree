
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
end

test_example()
