package braza

import (
	"fmt"
	"log"
	"maps"
	"net"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

/*
Create a new app with a default settings

	app := NewApp(nil) // or NewApp(*braza.Config{})
	...
	app.Listen()
*/
func NewApp(cfg *Config) *App {

	router := NewRouter("")
	router.main = true

	if cfg == nil {
		cfg = &Config{}
	}

	return &App{
		Config:       cfg,
		Router:       router,
		routers:      []*Router{router},
		routerByName: map[string]*Router{"": router},
	}
}

type App struct {
	// Main Router
	*Router

	// Main Config
	*Config

	/*
		custom auth parser
			app.BasicAuth = (ctx *braza.Ctx) {
				a := ctx.Request.Header.Get("Authorization")
				if a != ""{
					...
					return user,pass, true
				}
				return "","",false
			}
	*/
	BasicAuth func(*Ctx) (user string, pass string, ok bool) // custom func to parse Authorization Header
	/*
		exec after each request (if the application dont crash)
			app.AfterRequest = (ctx *braza.Ctx) {
				h := ctx.Response.Header()
				h.Set("X-Foo","Bar")
				ctx.Response.SetHeader(h)
				...
			}
	*/
	AfterRequest Func

	/*
		exec before each request
			app.BeforeRequest = (ctx *braza.Ctx) {
				db := database.Open()
				ctx.Global["db"] = db
				user,pass,ok := ctx.Request.BasicAuth()
				if ok {
					ctx.Global["user"] = user
				}
			}
	*/
	BeforeRequest Func

	/*
		exec after close each conn ( this dont has effect in response)
			app.TearDownRequest = (ctx *braza.Ctx) {
				database.Close()
				log.Print(...)
				...
			}
	*/
	TearDownRequest Func

	// The Http.Server
	Srv *http.Server

	/*
		func main(){
			sockPath := "/temp/my_http.sock"
			if err := os.RemoveAll(sockPath); err != nil {
				log.Fatal(err)
			}

			l, err := net.Listen("unix", sockPath)
			if err != nil {
				log.Fatal(err)

			}
			defer l.Close()

			app := braza.NewApp(nil)
			app.Listener = l
			app.Listen()
		}

	*/
	Listener net.Listener

	routers      []*Router
	routerByName map[string]*Router
	built        bool
}

/*
APP methods
*/

func (app *App) parseServername() {
	if app.Servername != "" {
		srv := strings.TrimPrefix(
			strings.TrimPrefix(
				app.Servername, "https//",
			),
			"https://",
		) // so pra evitar erros

		srv = strings.TrimPrefix(
			strings.TrimSuffix(
				srv, "/",
			), ".") // so pra evitar erros²

		h, p, err := net.SplitHostPort(srv)
		if err != nil && p != "" && h != "" {
			log.Fatal(err)
		}
		if p != "" {
			app.serverport = p
		}
		if h != "" {
			app.hostname = h
		}
	}
}

func (app *App) parseStaticRoute() {
	if !app.DisableStatic {
		staticUrl := "/assets"
		fp := "/{filepath:path}"
		if app.StaticUrlPath != "" {
			staticUrl = app.StaticUrlPath
		}
		path := filepath.Join(staticUrl, fp)
		app.AddRoute(&Route{
			Url:      path,
			Func:     serveFileHandler,
			Name:     "assets",
			isStatic: true,
		})
	}
}

func (app *App) parseAppEnv() {
	if env, ok := allowEnv[app.Env]; ok {
		app.Env = env
	} else {
		l.err.Fatalf("environment '%s' is not valid", app.Env)
	}
}

// se o usuario mudar o router principal, isso evita alguns erro
func (app *App) parseMainRouter() {
	if !app.main {
		app.main = true

		if app.Router.Routes == nil {
			app.Router.Routes = []*Route{}
		}
		if app.Router.routesByName == nil {
			app.Router.routesByName = map[string]*Route{}
		}
		if app.Router.Cors == nil {
			app.Router.Cors = &Cors{}
		}
		if app.Router.Middlewares == nil {
			app.Router.Middlewares = []Func{}
		}
	}
}

// Parse the router and your routes
func (app *App) parseApp() {
	app.checkConfig()
	app.parseServername()
	app.parseStaticRoute()
	app.parseAppEnv()
	app.parseMainRouter()

	if !slices.Contains(app.routers, app.Router) {
		app.routers = append(app.routers, app.Router)
	}
	for _, router := range app.routers {
		router.parse(app.Servername)
		maps.Copy(app.routesByName, router.routesByName)
	}
}

/*
Register Router in app

	func main() {
		api := braza.NewRouter("api")
		api.post("/products")
		api.get("/products/{productID:int}")

		app := braza.NewApp(nil)
		app.Mount(api)
		// do anything ...
		app.Listen()
	}
*/
func (app *App) Mount(routers ...*Router) {
	for _, router := range routers {
		if router.Name == "" {
			panic(fmt.Errorf("the routers must be named"))
		} else if _, ok := app.routerByName[router.Name]; ok {
			panic(fmt.Errorf("router '%s' already regitered", router.Name))
		}
		if !slices.Contains(app.routers, router) {
			app.routers = append(app.routers, router)
			app.routerByName[router.Name] = router
		}
	}
}

/*
Build the App, but not start serve

example:

	func index(ctx braza.Ctx){}

	// it's work
	func main() {
		app := braza.NewApp()
		app.GET("/",index)
		app.Build()
		app.UrlFor("index",true)
	}
	// it's don't work
	func main() {
		app := braza.NewApp()
		app.GET("/",index)
		app.UrlFor("index",true)
	}
*/
func (app *App) Build(addr ...string) {
	app.setEnv()
	app.parseAppEnv() // redundant
	l = newLogger(app.LogFile)
	app.parseApp()
	app.built = true

	var address = ":5000"
	if len(addr) > 0 {
		a_ := addr[0]
		if a_ != "" {
			_, _, err := net.SplitHostPort(a_)
			if err == nil {
				address = a_
			}
		}
	}

	if strings.Contains(address, "0.0.0.0") {
		listenAll = true
	}
	app.setFlags()
	app.parseSrvApp(address)
}

/*
	 Url Builder

		app.GET("/users/{userID:int}", index)

		app.UrlFor("index", false, "userID", "1"}) //  /users/1
		app.UrlFor("index", true, "userID", "1"}) // http://servername/users/1
*/
func (app *App) UrlFor(name string, external bool, args ...string) string {
	var (
		host   = ""
		route  *Route
		router *Router
	)
	if !app.built {
		l.err.Fatalf("you are trying to use this function outside of a context")
	}
	if len(args)%2 != 0 {
		l.err.Fatalf("numer of args of build url, is invalid: UrlFor only accept pairs of args ")
	}

	// check route name
	if r, ok := app.routesByName[name]; ok {
		route = r
	} else {
		panic(fmt.Sprintf("Route '%s' is undefined \n", name))
	}
	router = route.router

	params := map[string]string{}
	for i := range len(args) {
		if i%2 == 0 {
			params[args[i]] = args[i+1]
		}
	}

	// Build Host
	if external {
		schema := "http://"
		if app.ListeningInTLS {
			schema = "https://"
		}
		srvname := app.Servername
		if app.serverport != "" && (app.serverport != "80" && app.serverport != "443") {
			srvname = net.JoinHostPort(app.Servername, app.serverport)
		}
		if router.Subdomain != "" {
			host = schema + router.Subdomain + "." + srvname
		} else {
			if app.Servername == "" {
				_, p, _ := net.SplitHostPort(app.Srv.Addr)
				h := net.JoinHostPort(localAddress, p)
				host = schema + h
			} else {
				host = schema + srvname
			}
		}
	}
	url := route.mountURI(args...)
	return host + url
}

/*
Custom Http Error Handler

	app.ErrorHandler(401,(ctx *braza.Ctx) {
		ctx.HTML("Access denied",401)
	})
	app.ErrorHandler(404,(ctx *braza.Ctx) {
		ctx.HTML("Hey boy, you're a little lost",404)
	})
*/
func (app *App) ErrorHandler(statusCode int, f Func) {
	if app.errHandlers == nil {
		app.errHandlers = map[int]Func{}
	}
	app.errHandlers[statusCode] = f
}

func (app *App) GetRouters() []*Router { return app.routers }

func (app *App) GetMapRouters() map[string]*Router { return app.routerByName }

func (app *App) GetRouterByName(name string) *Router { return app.routerByName[name] }
