# braza


## [See the Documentation](https://github.com/5tkgarage/braza/blob/main/docs)

## Features
    - File Watcher (hot reload)
    - Error management & custom response errors
    - Dynamic route paths & subdomains
    - Schema Validator (with c3po package)
    - Rendering built-in (template/html)
    - Implements net/http
    - ...

    - Supports
        - Jwt 
        - Cors 
        - Sessions
        - Middlewars
        - URL Route builder

## Simple Example

### [With a correctly configured Go toolchain:](https://go.dev/doc/install)

```sh
go get 5tk.dev/braza
```

 _main.go_

```go
package main

import "5tk.dev/braza"

func main() {
 app := braza.NewApp()
 app.Get("/hello", helloWorld)
 app.Get("/hello/{name}", helloUser) // 'name' is any string
 app.Get("/hello/{userID:int}", userByID) // 'userID' is only int

 fmt.Println(app.Listen())
}

func helloWorld(ctx *braza.Ctx) {
 hello := map[string]any{"Hello": "World"}
 ctx.JSON(hello, 200)
}

func helloUser(ctx *braza.Ctx) {
 rq := ctx.Request   // current Request
 name := rq.PathArgs["name"]
 ctx.HTML("<h1>Hello "+name+"</h1>", 200)
}

func userByID(ctx *braza.Ctx) {
 rq := ctx.Request   // current Request
 id := rq.PathArgs["userID"]
 user := AnyQuery(id)
 ctx.JSON(user, 200)
}
```
