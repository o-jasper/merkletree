local sha2 = require 'merkle.sha2'
local enhex = require "merkle.enhex"

local fd = io.open("/dev/stdin")
local str = fd:read("l")
fd:close()
print(enhex(sha2.hash256(str or "")))
