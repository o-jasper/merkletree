--  Copyright (C) 27-01-2016 Jasper den Ouden.
--
--  This is free software: you can redistribute it and/or modify
--  it under the terms of the GNU General Public License as published
--  by the Free Software Foundation, either version 3 of the License, or
--  (at your option) any later version.

local This = {}
This.__index = This

function This:new(new)
   setmetatable(new, self)
   new:init()
   return new
end

function This:init()
   assert(self.H,
          "Need a pair hash function! Can use Class:class_pairify to use a hash")
   if not self.H2 then self.H2 = This:class_pairify() end
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

function This:verify_H(root, list, leaf)
   assert(self.H)
   return self:verify_H(root, list, self.H(leaf))
end

function This:verify_H(root, list, leaf_H)
   local cur_H = leaf_H
   for _, el in ipairs(list) do
      cur_H = self.H2(cur_H, el)
   end
   return cur_H
end

return This
