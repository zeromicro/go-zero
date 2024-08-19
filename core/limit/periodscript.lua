-- to be compatible with aliyun redis, we cannot use `local key = KEYS[1]` to reuse the key
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCRBY", KEYS[1], 1)
if current == 1 then
    redis.call("expire", KEYS[1], window)
end
if current < limit then
    return 1
elseif current == limit then
    return 2
else
    return 0
end