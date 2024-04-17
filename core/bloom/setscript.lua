for _, offset in ipairs(ARGV) do
    redis.call("setbit", KEYS[1], offset, 1)
end