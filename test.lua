-- http://lua-users.org/wiki/FileInputOutput

function file_exists(file)
   local f = io.open(file, 'rb')
   if f then f:close() end
   return f ~= nil
end

function lines_from(file)
   if not file_exists(file) then return {} end
   lines = {}
   for line in io.lines(file) do
      lines[#lines+1] = line
   end
   return lines
end

local file = 'urls.txt'
local addrs = lines_from(file)

request = function()
   local addr = 'https://'..addrs[math.random(#addrs)]
   wrk.headers['Host'] = addr
   return wrk.format('GET', '/')
end
