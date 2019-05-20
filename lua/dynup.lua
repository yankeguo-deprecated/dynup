-- dynup.lua
-- dynamic upstream resolving

local redis = require 'resty.redis'
local cjson = require 'cjson'

-- helper functions
function dynup_error(reason)
    ngx.log(ngx.ERR, reason)
    ngx.status = 503
    ngx.header['Content-Type'] = 'text/html; charset=utf-8'
    ngx.say('<h1>Dynup Error:</h1><h2>' .. reason .. '</h2>')
    ngx.exit(503)
end

function dynup_warn(reason)
    ngx.log(ngx.ERR, reason)
end

function dynup_match(v, p)
    if not v then
        return false
    end
    v = tostring(v)
    if p:sub(0, 1) == '/' and p:sub(-1) == '/' then
        p = p:sub(2, -2)
        return not (not v:find(p))
    end
    if p:sub(0, 1) == '[' and p:sub(-1) == ']' then
        p = p:sub(2, -2)
        for sub in p:gmatch('[^%s,]+') do
            if v:lower() == sub:lower() then
                return true
            end
        end
        return false
    end
    if p:sub(0, 1) == '<' and p:sub(-1) == '>' then
        p = p:sub(2, -2)
        for l, h in p:gmatch('([^,]+),%s*([^,]+)') do
            return tonumber(v) >= tonumber(l) and tonumber(v) <= tonumber(h)
        end
        return false
    end
    return v:gsub('%s', ''):lower() == p:lower()
end

local project = ngx.var.dynup_project
local redis_host = ngx.var.dynup_redis_host
local redis_port = ngx.var.dynup_redis_port

-- check variables
if not project then
    return dynup_error('$dynup_project not set in nginx.conf')
end

local dynup_key_rules = 'dynup.projects.' .. project .. '.rules'
local dynup_key_backends = 'dynup.projects.' .. project .. '.backends'

if not redis_host then
    dynup_warn("$dynup_redis_host not set, default to '127.0.0.1'")
    redis_host = '127.0.0.1'
end

if not redis_port then
    dynup_warn('$dynup_redis_port not set, default to 6379')
    redis_port = 6379
else
    redis_port = tonumber(redis_port)
end

-- redis client
local red = redis:new()
red:set_timeout(3000)

local ok, err = red:connect(redis_host, redis_port)
if not ok then
    return dynup_error('failed to connect redis ' .. redis_host .. ':' .. tostring(redis_port))
end

-- fetch frontend rules
local res, err = red:get(dynup_key_rules)
if not res or err then
    return dynup_error('failed to fetch rules for project: ' .. project)
end
if res == ngx.null then
    return dynup_error('rules not set for project: ' .. project)
end

-- default backend
local backend = 'default'
local backend_set = false

-- determine backend
local rules = cjson.decode(res)

for _, v in pairs(rules) do
    if v.type == 'header' then
        for _, f in pairs(v.fields) do
            if dynup_match(ngx.req.get_headers()[f], v.pattern) then
                backend = v.target
                backend_set = true
                break
            end
        end
    end
    if v.type == 'query' then
        for _, f in pairs(v.fields) do
            if dynup_match(ngx.req.get_uri_args()[f], v.pattern) then
                backend = v.target
                backend_set = true
                break
            end
        end
    end

    if backend_set then
        break
    end
end

-- backend rr
local dynup_key_backend_rr = 'dynup.projects.' .. project .. '.backends.' .. backend .. '.rr-cur'

-- fetch backends
res, err = red:get(dynup_key_backends)
if not res or err then
    return dynup_error('failed to fetch backends for project: ' .. project)
end
if res == ngx.null then
    return dynup_error('backends not set for project: ' .. project)
end

-- fetch targets
local targets = cjson.decode(res)[backend]

if not targets then
    return dynup_error('backend: ' .. backend .. ' not set for project: ' .. project)
end

-- rotate the rr
local cursor, err = red:incr(dynup_key_backend_rr)
if not cursor or err then
    return dynup_error('failed to fetch cursor for backend: ' .. backend .. ' for project: ' .. project)
end
if cursor == ngx.null then
    return dynup_error('cursor for backend: ' .. backend .. ' not set for project:' .. project)
end

cursor = tonumber(cursor)

if cursor > 999999999 then
    red:set(dynup_key_backend_rr, 0)
end

-- set the $dynup_upstream variable
ngx.var.dynup_upstream = targets[1 + cursor % table.maxn(targets)]
