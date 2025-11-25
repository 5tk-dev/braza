package posts

import (
	"time"

	"5tk.dev/braza"
	"5tk.dev/braza/examples/blog/models"
	"github.com/google/uuid"
)

func middleAuth(ctx *braza.Ctx) {
	if _, ok := ctx.Global["user"]; !ok {
		ctx.Unauthorized()
	}
}

func getMany(ctx *braza.Ctx) {
	dbPost := models.GetDBPost()
	ctx.JSON(dbPost, 200)
}

func get(ctx *braza.Ctx) {
	dbPost := models.GetDBPost()
	post_uuid := ctx.Request.PathArgs["uuid"]
	if post, ok := dbPost[post_uuid]; ok {
		ctx.JSON(post, 200)
	}
	ctx.NotFound()
}

func post(ctx *braza.Ctx) {
	dbPost := models.GetDBPost()
	user := ctx.Global["user"].(*models.User)
	text, ok := ctx.Request.Form["text"].(string)
	if !ok {
		ctx.BadRequest()
	}
	new_uuid := uuid.NewString()
	dbPost[new_uuid] = &models.Post{
		UUID:    new_uuid,
		User:    user.Email,
		Text:    text,
		Created: time.Now(),
	}

	ctx.JSON(dbPost[new_uuid], 201)
}

func put(ctx *braza.Ctx) {
	dbPost := models.GetDBPost()
	user := ctx.Global["user"].(*models.User)

	sch := ctx.Schema.(*PutPostSchema)
	post, ok := dbPost[sch.UUID]
	if !ok {
		ctx.NotFound()
	}
	if post.User != user.Email {
		ctx.Forbidden()
	}
	post.Text = sch.Text
	// db.save(...)
	ctx.JSON(post, 200)
}

func del(ctx *braza.Ctx) {
	dbPost := models.GetDBPost()
	user := ctx.Global["user"].(*models.User)
	post_uuid := ctx.Request.PathArgs["uuid"] // same name from url path var ex.:-> /foo/{uuid}

	post, ok := dbPost[post_uuid]
	if !ok {
		ctx.NotFound()
	}
	if post.User != user.Email {
		ctx.Forbidden()
	}
	delete(dbPost, post_uuid)
	ctx.NoContent()
}
