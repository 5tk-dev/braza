package main

import (
	"5tk.dev/braza"
	"5tk.dev/braza/examples/blog/api/posts"
	"5tk.dev/braza/examples/blog/auth"
	"5tk.dev/braza/examples/blog/models"
)

func main() {
	app := braza.NewApp(&braza.Config{
		SecretKey:     "batata",
		DisableStatic: true,
	})
	app.BeforeRequest = beforeRequest

	app.Mount(auth.NewRouter())
	app.Mount(posts.NewRouter())

	app.ShowRoutes() // same -> go run . -route
	app.Listen(nil)
}

func beforeRequest(ctx *braza.Ctx) {
	email := ctx.Session.Get("user")
	if email != "" {
		dbUser := models.GetDBUser()
		if user, ok := dbUser[email]; ok {
			ctx.Global["user"] = user
		}
	}
}
