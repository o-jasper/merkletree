--  Copyright (C) 23-02-2016 Jasper den Ouden.
--
--  This is free software: you can redistribute it and/or modify
--  it under the terms of the Afrero GNU General Public License as published
--  by the Free Software Foundation, either version 3 of the License, or
--  (at your option) any later version.

-- Note: typically the object is one-use.
-- In case of merkle tree, merkle tree API is available.
local Statementize = {}
Statementize.__index= "i am virtual"

function Statementize:new(new)
   new = setmetatable(new, self)
   new:init()
   return new
end

function Statementize:init()
   error("init needs to be externally defined.")
end

function Statementize:class_derive(...)
   local New = {}
   for k,v in pairs(self)    do New[k] = v end
   for _, Replace in ipairs{...} do
      for k,v in pairs(Replace) do New[k] = v end
   end
   New.__index = New
   assert(New.hash and New.make)
   assert(New.init and New.add_key,
          "Required functions not defined.(:init or :add_key(key,val)")
   return New
end

-- Statementize.hash
Statementize.encode = require("storebin").encode

Statementize.add_n = 0
local function hashtree(self, tree, front)

   local encode = self.encode

   local function hashtree_val(key, val)
      assert(({number=true,string=true})[type(key)])

      if type(val) == "table" then  -- Recurse into branch.
         hashtree(self, val, front .. encode(key))
      else
         self:add_key(front .. encode(key), encode(val))
      end
   end

   local into = { number={}, string={} }
   for k in pairs(tree) do
      local list = into[type(k)]
      assert(list, "Only number or string keys, got; " .. type(k))
      table.insert(list, k)
--      if #list > 0 then  -- Switch to this to check `pairs` indetermism doesnt matter.
--         table.insert(list, math.random(#list), k)
--      end
   end

   if #into.number == 0 and #into.string == 0 then  -- Completely empty. I reserve __em
      self:add_key(front .. encode("__em"), encode(true))
   else
      for _,list in ipairs{into.number, into.string} do -- `pairs(into)` wont do!
         table.sort(list)
         for _, k in ipairs(list) do hashtree_val(k, tree[k]) end
      end
   end
end

function Statementize:hash(tree)
   hashtree(self, tree, "")
   return self:close()
end

local b64 = require "page_html.util.fmt.base64"

function Statementize:hashstr(tree) return b64.enc(self:hash(tree)) end

--Note: these are overly high values for many applications..
Statementize.nonce_size = 16
function Statementize:gen_nonce()
   local fd = io.open("/dev/random")
   local nonce = fd:read(self.nonce_size)
   fd:close()
   return  nonce
end

Statementize.nonce_assert_size = true  -- Advisable.

-- Note always/never also means it cannot be used to verify those without/with nonce.
-- (if neither, it can verify both, makes nonces if `.gen_nonce`, doesnt if not.
Statementize.always_nonce = true
--Statementize.never_nonce   = false

function Statementize:make(tree)
   if self.never_nonce then
      assert( not (self.gen_nonce or tree.nonce) )
      return self.name .. ":" .. self:hashstr(tree)
   else
      local nonce = tree.nonce or (self.gen_nonce and self:gen_nonce())
      if nonce then
         tree.nonce = nil  -- Take it out, otherwise double.
         assert(type(nonce) == "string")
         assert(not self.nonce_assert_size or self.nonce_size == #nonce)

         local ret = self.name .. ":" .. b64.enc(nonce) .. ":" .. self:hashstr(tree)
         tree.nonce = nonce  -- Put it back on.
         return ret
      else
         assert(not self.always_nonce)
         return self.name .. ":" .. self:hashstr(tree)
      end
   end
end

function Statementize:verify(tree, statement_str)
   local nonce = string.match(statement_str, ":([^:]+):")
   if nonce then
      if tree.nonce and tree.nonce ~= nonce then
         return false, "nonces mismatch"
      end
      tree.nonce = nonce
      local got_statement_str = self:make(tree)
      return got_statement_str == statement_str, "nonced result", got_statement_str
   elseif tree.nonce then
      return false, "statement_str no nonce, yet tree does"
   else
      local got_statement_str = self:make(tree)
      return got_statement_str == statement_str, "result", got_statement_str
   end
end

return Statementize
