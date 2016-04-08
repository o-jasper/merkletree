local Public = { __constant=true }

for k,v in pairs(require "merkle.statement") do
   if type(v) == "table" then
      Public[v.name] = v
   end
end

for k,v in pairs(require "merkle.statement.merkle") do
   if type(v) == "table" then
      Public[v.name] = v
   end
end

return Public
