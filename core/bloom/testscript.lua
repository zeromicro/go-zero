for _, offset in ipairs(ARGV) do
    if tonumber(redis.call("getbit", KEYS[1], offset)) == 0 then
        return false
    end
end
return true