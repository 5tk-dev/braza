package auth

import "5tk.dev/braza"

func NewRouter() *braza.Router {
	return &braza.Router{
		Name:        "auth",
		Prefix:      "/auth",
		Routes:      routes,
		Middlewares: []braza.Func{middleAuth},
	}
}
