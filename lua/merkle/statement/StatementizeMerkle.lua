--  Copyright (C) 23-02-2016 Jasper den Ouden.
--
--  This is free software: you can redistribute it and/or modify
--  it under the terms of the Afrero GNU General Public License as published
--  by the Free Software Foundation, either version 3 of the License, or
--  (at your option) any later version.

-- Note unlike plain statementize, this one is FULLY A CLASS.

local Statementize = require "merkle.statement.Statementize"

local MerkleTree = require "merkle.Merkle.Tree"

local This = Statementize:class_derive(MerkleTree)

local encode = require("storebin").encode

This.encode = nil

local function hashtree(self, tree, front)

   local function hashtree_val(key, val)
      if type(val) == "table" then  -- Recurse into branch.
         hashtree(self, val, front .. encode(key))
      else
         self:add_key(front .. encode(key), encode(val))
      end
   end

   local into = { number={}, string={} }
   for k in pairs(tree) do
      local list = into[type(k)]
      assert(list, "Only number or string keys")
      table.insert(list, k)
   end

   for _,list in ipairs{into.number, into.string} do -- `pairs(into)` wont do!
      table.sort(list)
      for _, k in ipairs(list) do hashtree_val(k, tree[k]) end
   end
end

function This:hash(tree)
   assert(type(tree) == "table")
   hashtree(self, tree, front)
   return self:close()
end

return This
