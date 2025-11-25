package posts

import "5tk.dev/braza"

var routes = []*braza.Route{
	// url == "" -> "/"
	// methods == [] -> GET | HEAD
	{
		Name: "get",
		Func: get,
		MapCtrl: braza.MapCtrl{
			"GET":  &braza.Meth{Func: getMany},
			"POST": &braza.Meth{Func: post},
		},
	},
	{
		Name: "putOrDelete",
		Url:  "/{uuid}",
		MapCtrl: braza.MapCtrl{
			"GET": &braza.Meth{Func: get},
			"PUT": &braza.Meth{
				Func:   put,
				Schema: &PutPostSchema{},
			},
			"DELETE": &braza.Meth{Func: del},
		},
	},
}
