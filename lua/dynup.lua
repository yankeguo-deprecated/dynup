-- dynup.lua
-- dynamic upstream resolving

local redis = require "resty.redis"

-- helper functions
function dynup_error(reason)
    ngx.log(ngx.ERR, reason)
    ngx.status = 503
    ngx.header["Content-Type"] = "text/html; charset=utf-8"
    ngx.say("<h1>Dynup Error:</h1><h2>"..reason.."</h2>")
    ngx.exit(503)
end

function dynup_warn(reason)
    ngx.log(ngx.ERR, reason)
end

function dynup_trim(s)
    return s:match'^%s*(.*%S)' or ''
end

local project = ngx.var.dynup_project
local redis_host = ngx.var.dynup_redis_host
local redis_port = ngx.var.dynup_redis_port;

-- check variables
if not project then
    return dynup_error("$dynup_project not set in nginx.conf")
end

if not redis_host then
    dynup_warn("$dynup_redis_host not set, default to '127.0.0.1'")
    redis_host = "127.0.0.1"
end

if not redis_port then
    dynup_warn("$dynup_redis_port not set, default to 6379")
    redis_port = 6379
else
    redis_port = tonumber(redis_port)
end

-- redis client
local red = redis:new()
red:set_timeout(1000) -- 1 sec

local ok, err = red:connect(redis_host, redis_port)
if not ok then
    return dynup_error("failed to connect redis "..redis_host..":"..tostring(redis_port))
end

-- fetch frontend rules
local res, err = red:get("dynup.projects."..project..".frontend-rules")
if not res then
    return dynup_error("failed to fetch frontend rules")
end
if res == ngx.null then
    return dynup_error("frontend rules not set")
end

-- eval frontend rules
res = dynup_trim(res)

-- fetch backends
-- TODO:

-- set the $dynup_upstream variable
-- TODO:
ngx.var.dynup_upstream = "127.0.0.1:8081"