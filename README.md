# dynup

Dynamic Upstream with OpenResty and Redis

## Redis Keys

### Projects

Project names are stored as a set

```plain
> SMEMBERS dynup.projects

1) foo
2) bar
```

### Backend Rules

Backend rules are stored as a tab (`\t`) seperated multiline text

```plain
> GET dynup.projects.foo.backend-rules

# source    rules                   port    group
consul      /dev\.foo\.\d+/         8081    default # consul host matching regexp (with leading '/' and trailing '/')
consul      dev.foo.v-1,dev.foo.v-2 8082    default # consul host, support comma
manual      10.10.10.2,10.10.10.12  8080    canary  # manual ip, support comma
```

Backend group `default` MUST exist

### Backend Hosts

Backends are stored as a comma seperated list

`dynupd` will continusly resolve `backend-rules` into `backends`

```plain
> GET dynup.projects.foo.backends.default

10.10.10.3:8081,10.10.10.4:8081,10.10.10.5:8081

> GET dynup.projects.foo.backends.canary

10.10.10.2:8080,10.10.10.12:8080
```

### Backend RR Cursor

Backend RR cursor can be retrieved with `INCR` command

```plain
> INCR GET dynup.projects.foo.backends.default.rr-cur

80
```

### Frontend Rules

Frontend rules are stored as a tab (`\t`) seperated multiline text

Frontend rules will be evaluated by Lua script

```plain
> GET dynup.projects.foo.frontend-rules

# type  field               value           target
header  X-Canary,Canary     true            canary  # force canary  on X-Canary = true
header  X-Canary,Canary     false           default # force default on X-Canary = false
header  X-User-ID,X-USERID  /100\d+/        canary  # header matching regexp (with leading '/' and trailing '/')
header  X-User-ID,X-USERID  [1,2,3,4,5]     canary  # header matching comma seperated list (with leading '[' and trailing ']')
header  X-User-ID,X-USERID  <12,38>         canary  # header matching inclusive range (with leading '<' and trailing '>'), will convert to number 
```

Other supported types are `query` and `cookie`

If no rule matched, will fallback to `default` group
