package posts

import "5tk.dev/braza"

func NewRouter() *braza.Router {
	return &braza.Router{
		Name:        "posts",
		Prefix:      "/api/posts",
		Routes:      routes,
		Middlewares: []braza.Func{middleAuth},
	}
}
