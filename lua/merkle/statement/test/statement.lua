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


local function rand_i()
   local r = 0
   for _ = 1,6 do r = 256*r + string.byte(rand_str(1,1)) end
   return r
end

math.randomseed(os.time() + 10000*os.clock())
math.randomseed(rand_i())

local gen_tree = require("storebin.test.lib")[2]

local statement = require "merkle.statement"

local function run_1(Which, tree)
   tree = tree or gen_tree(true, 5, {mini=1, maxi=102, no_boolean=true})
   tree.nonce = nil  -- Just in case.

   local r1 = Which:new():make_text(tree)
-- Really kindah only test the function, only more if `pairs` nondeterministic.
   assert(Which:new():verify_from_text(tree, r1))
   print(r1)
end

-- sha2 does not work with incorrect init.
local sha2 = require "merkle.sha2"
for _,k in ipairs{"Sha256", "Sha224"} do
   assert(sha2[k].init == statement[k].init )
   assert(sha2[k].init == statement[k .. "N"].init )
end

run_1(statement.Sha256)
run_1(statement.Sha256N)
run_1(statement.Sha224)
run_1(statement.Sha224N)

local statement_merkle = require "merkle.statement.merkle"
--for _,k in ipairs{"Sha256", "Sha224"} do assert( sha2[k].init == statement[k].init ) end

run_1(statement_merkle.Sha256, {})
run_1(statement_merkle.Sha256, {x={y={}}})

run_1(statement_merkle.Sha256)
run_1(statement_merkle.Sha256N)
run_1(statement_merkle.Sha224)
run_1(statement_merkle.Sha224N)

-- "Damage" one to look if `pairs` indeterminism can break it.
print("INDETERMINISM")
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
      if #list > 0 then  -- TEST INDETERMINISM
         table.insert(list, math.random(#list), k)
      end
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

function statement_merkle.Sha224N:hash(tree)
   hashtree(self, tree, "")
   return self:close()
end
-- Note: probably good for optimizability of lua if this werent possible :)

run_1(statement_merkle.Sha224N)
