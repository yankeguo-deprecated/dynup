# 如何集成 Dynup

## 使用 OpenResty 载入 Lua 脚本

```nginx
location / {
    set $dynup_project      test;       # 项目名，以 test 为例
    set $dynup_redis_host   127.0.0.1;  # Redis IP
    set $dynup_redis_port   6379;       # Redis 端口
    set $dynup_redis_pass   '';         # Redis 密码，没有则必须留空字符串
    set $dynup_upstream     '';         # 必须预先定义空值，才能在脚本中写入该值
    access_by_lua_file /path/to/lua/dynup.lua;  # 执行脚本，写入 $dynup_stream
    proxy_pass http://$dynup_upstream;  # 使用 $dynup_upstream
}
```

## 使用 Redis 写入配置

### 后端服务器记录

将后端服务器记录写入 `gateway-backends-test` 键值

Dynup 使用分组概念，后端服务器被分为多个组。其中 `default` 组为默认组，必须存在。

```
> SET gateway-backends-test

{
  "canary": [
    "127.0.0.1:8082"
  ],
  "default": [
    "127.0.0.1:8081"
  ]
}
```

### 前端路由规则

将前端路由规则写入 `gateway-rules-test` 键值

```plain
> SET gateway-rules-test

[
  {
    "type": "query",            // 支持 query 和 header
    "fields": [                 // 填写需要匹配的字段
      "canary",
      "is_canary"
    ],
    "pattern": "true",          // 字符串匹配
    // "pattern": "/[Tt].+/",   // Lua 正则表达式，使用 // 包裹
    // "pattern": "[1,2,3]",    // 列表匹配，使用 [] 包裹
    // "pattern": "<1,10>",     // 数字范围，使用 <> 包裹
    "target": "canary"          // 目标后端服务器分组
  }
]
```
