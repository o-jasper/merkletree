--  Copyright (C) 27-01-2016 Jasper den Ouden.
--
--  This is free software: you can redistribute it and/or modify
--  it under the terms of the GNU General Public License as published
--  by the Free Software Foundation, either version 3 of the License, or
--  (at your option) any later version.

local This = {}
This.__index = This

function This:new(new)
   new = setmetatable(new or {}, self)
   new:init()
   return new
end

function This:init()
   assert(self.H,
          "Need a pair hash function! Can use Class:class_pairify to use a hash")
   if not self.H2 then self.H2 = self:class_pairify() end
end

function This:class_pairify(H)
   local H = H or self.H
   return function(x, y)
      if x > y then
         return H(y .. x)
      else
         return H(x .. y)
      end
   end
end

function This:expect_root(proof, leaf)
   assert(self.H)
   return self:expect_root_H(proof, self.H(leaf))
end

function This:expect_root_H(proof, leaf_H)
   local cur_H = leaf_H
   for _, el in ipairs(proof) do
      cur_H = self.H2(cur_H, el)
   end
   return cur_H 
end


function This:verify(root, proof, leaf)
   assert(self.H)
   return self:verify_H(root, proof, self.H(leaf))
end

function This:verify_H(root, proof, leaf_H)
   return self:expect_root_H(proof, leaf_H) == root
end

function This:verify_key(root, proof, key, leaf)
   return self:verify(root, proof, key .. leaf)
end

return This
