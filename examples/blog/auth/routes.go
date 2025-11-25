package auth

import (
	"5tk.dev/braza"
	"5tk.dev/braza/examples/blog/models"
)

func middleAuth(ctx *braza.Ctx) {
	_, ok := ctx.Global["user"]
	if ok {
		ctx.Redirect("/")
	}
}

var routes = []*braza.Route{
	braza.Post("/register", register),
	braza.Post("/login", login),
}

func register(ctx *braza.Ctx) {
	email, pass, ok := ctx.Request.BasicAuth()
	if !ok {
		ctx.BadRequest()
	}
	dbUser := models.GetDBUser()
	if _, ok := dbUser[email]; ok {
		ctx.Unauthorized() // ctx.JSON -> email unavaiable, 400
	}
	name := ctx.Request.Form["name"].(string)
	dbUser[email] = &models.User{
		Name:       name,
		Email:      email,
		HashedPass: pass,
	}
	ctx.Session.Set("user", email)
	ctx.Created()
}

func login(ctx *braza.Ctx) {
	email, pass, ok := ctx.Request.BasicAuth()
	if !ok {
		ctx.BadRequest()
	}
	dbUser := models.GetDBUser()

	user, ok := dbUser[email]
	if !ok {
		ctx.Unauthorized()
	}
	if pass != user.HashedPass {
		ctx.Unauthorized()
	}
	ctx.Session.Set("user", email)
	ctx.Ok()
}
