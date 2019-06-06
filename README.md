# dynup

Dynamic Upstream with OpenResty and Redis

## Integration

```nginx
location / {
    set $dynup_project      test;
    set $dynup_redis_host   127.0.0.1;
    set $dynup_redis_port   6379;
    set $dynup_redis_pass   '';         # must exist, empty for null
    set $dynup_upstream     '';         # must pre-define
    access_by_lua_file ../lua/dynup.lua;
    proxy_pass http://$dynup_upstream;
}
```

## Redis Keys

### Projects (for Web UI Only)

Project names are stored as a set

```plain
> SMEMBERS dynup.projects

1) foo
2) bar
```

### Backends

```plain
> GET gateway-backends-foo

{
  "canary": [
    "127.0.0.1:8082"
  ],
  "default": [
    "127.0.0.1:8081"
  ]
}
```

Backend group `default` MUST exist

### Rules

Frontend rules will be evaluated by Lua script

```plain
> GET gateway-rules-foo

[
  {
    "type": "query",
    "fields": [
      "canary",
      "is_canary"
    ],
    "pattern": "true",          // exact string match
    // "pattern": "/[Tt].+/",   // lua regexp
    // "pattern": "[1,2,3]",    // list
    // "pattern": "<1,10>",     // range of number
    "target": "canary"
  }
]
```

Other supported types are `header`

If no rule matched, will fallback to `default` group

## License

Yanke Guo <guoyk@pagoda.com.cn>, MIT License
