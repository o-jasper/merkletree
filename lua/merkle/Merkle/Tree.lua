--  Copyright (C) 27-01-2016 Jasper den Ouden.
--
--  This is free software: you can redistribute it and/or modify
--  it under the terms of the GNU General Public License as published
--  by the Free Software Foundation, either version 3 of the License, or
--  (at your option) any later version.

-- Running-adding hash values.

local MerkleNode = require "merkle.Merkle.Node"
local MerkleVerify = require "merkle.Merkle.Verify"

local This = {}

for _,k in ipairs{"new", "class_pairify"} do This[k] = MerkleVerify[k] end

This.__index = This

function This:init()
   MerkleVerify.init(self)
   self.tops = {}
   self.kept_keys = {}
end

This.keep_proof_default = false
This.keep_data = true

function This:add(data, keep_proof)
   assert(self.H)
   local ret = self:add_H(self.H(data), keep_proof)
   ret.data = self.keep_data and data or nil
   return ret
end

function This:add_key(keydata, data, keep_proof)  -- Adds with key.
   local keep_proof = keep_proof or self.keep_proof_default or self.kept_keys[keydata]
   local ret = self:add(keydata .. data, keep_proof)
   if keep_proof then
      self.kept_keys[keydata] = ret
   end
   return ret
end

function This:add_H(H, keep_proof) return self:_add_H(H, 1, keep_proof) end

local tab_insert = table.insert

This.cnt_added = 0
function This:_add_H(H, n, keep_proof)
   self.cnt_added = self.cnt_added + 1
   assert(not self.finished)

   local keep_proof = (keep_proof==nil and self.keep_proof_default) or keep_proof
   local new = MerkleNode:new{ H=H, n=n, keep=keep_proof }
   tab_insert(self.tops, 1, new)
   self:_re_merge(false)
   return new
end

function This:_re_merge(super_force)
   local left, right = self.tops[1], self.tops[2]
   -- Keep putting two together if there are two of same depth, or using force.
   while right and (super_force or left.n == right.n) do
      local new_H, keep_either = self.H2(left.H, right.H), left.keep or right.keep
      local new = MerkleNode:new{ H=new_H, n=left.n+1, keep=keep_either,
                                  left=left, right=right }
      if keep_either then  -- Keep what is needed for proofs.
         left.parent  = new
         right.parent = new
      end
      table.remove(self.tops, 1)
      self.tops[1] = new  -- Replace the node with deeper node.

      left, right = self.tops[1], self.tops[2]
   end
end

-- Tentative final version,
function This:root_H_if_single()
   if #self.tops == 1 then return self.tops[1].H end
end

-- Finish it. (again, `:add` changes it again.
function This:close()
   if not self.finished then
      self:_re_merge(true) -- Force-merge everything.
      self.finished = true
   end
   assert(#self.tops == 1,
          string.format("It didnt re-merge it into one result. Got %s(%s)",
                        #self.tops, self.cnt_added))
   return self.tops[1].H
end

return This
