# router

This router is something to create the web simply.

## Router

```
    r := router.New(router.Template,nil)
    r.Add("/pattern",handler,http.MethodGet)
    http.Handle("/",r)
```

### Parameter

error handling

### LoginHandler


### Pattern

```
r.Add("/url/{arg}",handler)
p.Get("arg")

```


## Customize


### ErrorFunc

### TemplateFunc

### JSONFunc
