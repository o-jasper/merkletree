local This = { __constant=true }
This.__index = This

function This:new(new)
   return setmetatable(new, self)
end

-- Instead of the "partial tree" of proofs, makes a single list.
function This:produce_proof(ret_list)
   local cur, ret_list = self, ret_list or {}
   while cur.parent do  -- Go up, collecting the opposite side hashes.
      if cur.parent.left == cur then
         table.insert(ret_list, cur.parent.right.H)
      else
         assert(cur.parent.right == cur)
         table.insert(ret_list, cur.parent.left.H)
      end
      cur = cur.parent
   end
   return ret_list
end

return This
