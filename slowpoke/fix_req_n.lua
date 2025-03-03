

local counter = 1
local threads = {}

function setup(thread)
   table.insert(threads, thread)
   thread:set("x", -1)
end

function init(args)
   x = socket.gettime()
   wrk.thread:set("x", x)
end

function response()
   if counter == 12500 then
      wrk.thread:stop()
      local x = wrk.thread:get("x")
      y = socket.gettime() - x
      wrk.thread:set("y", y)
   end
   counter = counter + 1
end

function done(summary, latency, requests)
   io.write("------------------------------\n")
   for i, thread in ipairs(threads) do
      local y = thread:get("y")
      local msg = "stop time: %f"
      print(msg:format(y))
   end
end
