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

   local inst = Which:new()
   print(inst:make(tree))
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
