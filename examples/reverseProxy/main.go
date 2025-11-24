package main

import (
	"log"
	"net/http/httputil"
	"net/url"

	"5tk.dev/braza"
)

var rvp = httputil.ReverseProxy{
	Rewrite: func(pr *httputil.ProxyRequest) {
		host := pr.In.Header.Get("X-Host")
		u, _ := url.Parse("http://" + host)
		pr.SetURL(u)
	},
}

func main() {
	app := braza.NewApp(&braza.Config{
		DisableStatic:        true, // disble static handler
		DisableParseFormBody: true, //disables form parser
	})

	app.AddRoute(&braza.Route{
		Name: "proxy",
		Url:  "/{path:*}",
		Func: func(ctx *braza.Ctx) {
			// check if user has valid credentials
			// do anything ....
			ctx.Request.Header.Set("X-Host", "example.com")
			rvp.ServeHTTP(ctx, ctx.Request.HttpRequest())
		},
	})

	log.Fatal(app.Listen(nil))
}
