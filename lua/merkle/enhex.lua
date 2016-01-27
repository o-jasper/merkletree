
-- Bytes to hex.
local function enhex(s)
   local str = string.gsub(s, ".", function(c) 
                            return string.format("%02x", string.byte(c))
   end)
   return str
end

return enhex
