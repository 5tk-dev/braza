package main

import (
	"fmt"
	"log"

	"5tk.dev/braza"
)

type SchemaPost struct {
	Item string
}

type SchemaPutDel struct {
	Id   int `braza:"in=path"`
	Item string
}

var db = map[string]string{}

func main() {
	app := braza.NewApp(&braza.Config{
		Servername:    "localhost", // need for subdomains - 'port' is not required
		DisableStatic: true,        // disable static route - its is a api...
	})

	app.Prefix = "/v1"    // url prefix for this router
	app.Subdomain = "api" // subdomain for this router

	app.Cors = &braza.Cors{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Authorization", "*"},
	}

	app.Get("/todo", get)

	app.AddRoute(
		&braza.Route{
			Url: "/todo",
			// Name:    "post", // if omit, name set function name. ex: post
			Func:    post,
			Methods: []string{"POST"},
			Schema:  &SchemaPost{},
		},
	)

	app.AddRoute(&braza.Route{
		Url:  "/todo/{id:int}",
		Name: "todos",
		MapCtrl: braza.MapCtrl{
			"PUT": &braza.Meth{
				Func:   put,
				Schema: &SchemaPutDel{},
			},
			"DELETE": &braza.Meth{
				Func:   del,
				Schema: &SchemaPutDel{},
			},
		},
	})

	app.ShowRoutes() // can also be accessed by "go run . -routes"
	log.Fatal(app.Listen(nil))
}

func get(ctx *braza.Ctx) {
	ctx.JSON(db, 200)
}

func post(ctx *braza.Ctx) {
	item, ok := ctx.Request.Form["item"].(string)
	if !ok && item == "" {
		ctx.BadRequest()
	}
	db[fmt.Sprint(len(db))] = item
	ctx.JSON(db, 201)
}

func put(ctx *braza.Ctx) {
	id := ctx.Request.PathArgs["id"]
	item, ok := ctx.Request.Form["item"].(string)
	if !ok && item == "" {
		ctx.BadRequest()
	}

	if _, ok := db[id]; !ok {
		ctx.NotFound()
	}
	db[id] = item
	ctx.JSON(db, 200)
}

func del(ctx *braza.Ctx) {
	id := ctx.Request.PathArgs["id"]
	if _, ok := db[id]; ok {
		delete(db, id)
		ctx.NoContent()
	}
	ctx.NotFound()
}
