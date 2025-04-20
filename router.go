package braza

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

func NewRouter(name string) *Router {
	return &Router{
		Name:         name,
		Routes:       []*Route{},
		routesByName: map[string]*Route{},
	}
}

type Router struct {
	Cors        *Cors
	Name        string
	Routes      []*Route
	Prefix      string
	Subdomain   string // only {sub} or {sub:int}
	WsUpgrader  *websocket.Upgrader
	Middlewares []Func
	StrictSlash bool

	main           bool
	routesByName   map[string]*Route
	subdomainRegex []*regexp.Regexp
	errHandlers    map[int]Func
}

func (r *Router) compileSub() {
	r.subdomainRegex = make([]*regexp.Regexp, 0) // reset
	sub := r.Subdomain
	subSplit := strings.Split(sub, ".")
	for _, str := range subSplit {
		if str == "" {
			continue
		}
		if re.isVar.MatchString(str) {
			str = re.str.ReplaceAllString(str, `([a-zA-Z0-9\-\_]+)`)
			str = re.digit.ReplaceAllString(str, `(\d+)`)
			if re.isVar.MatchString(str) {
				panic(fmt.Errorf("only 'str' and 'int' are allowed in dynamic subdomains - Router:'%s', Subdomain:'%s'", r.Name, r.Subdomain))
			}
		}
		r.subdomainRegex = append(r.subdomainRegex, regexp.MustCompile(str))
	}
}

func (r *Router) parseRoute(route *Route) {
	if route.Name == "" {
		if route.Func == nil {
			panic("the route needs to be named or have a 'Route.Func'")
		}
		route.Name = getFunctionName(route.Func)
	}
	route.simpleName = r.Name
	if r.Name != "" {
		route.Name = r.Name + "." + route.Name
	}
	if _, ok := r.routesByName[route.Name]; ok {
		if route.isStatic {
			return
		}
		panic(fmt.Errorf("Route with name '%s' already registered", route.Name))
	}
	if r.Prefix != "" && !strings.HasPrefix(r.Prefix, "/") {
		panic(fmt.Errorf("Router '%v' Prefix must start with slash or be a null string ", r.Name))
	} else if route.Url != "" && (!strings.HasPrefix(route.Url, "/") && !strings.HasSuffix(r.Prefix, "/")) {
		panic(fmt.Errorf("Route '%v' Prefix must start with slash or be a null String", route.Name))
	}

	route.simpleUrl = route.Url
	route.Url = filepath.Join(r.Prefix, route.Url)
	route.parse()
	r.routesByName[route.Name] = route
	route.router = r
	route.parsed = true
}

func (r *Router) parse(servername string) {
	if r.routesByName == nil {
		r.routesByName = map[string]*Route{}
	}
	r.subdomainRegex = make([]*regexp.Regexp, 0) // reset

	if r.Name == "" && !r.main {
		panic(fmt.Errorf("the routers must be named"))
	}
	if r.Subdomain != "" {
		if servername == "" {
			panic(fmt.Errorf("to use subdomains you need to first add a ServerName in the app. Router:'%s'", r.Name))
		}
		r.compileSub()
	}

	if servername != "" {
		srvSplit := strings.Split(servername, ".")
		for _, s := range srvSplit {
			r.subdomainRegex = append(r.subdomainRegex, regexp.MustCompile(s)) // if nao tiver subdomain, ainda precisa usar o servername
		}
	}

	for _, route := range r.Routes {
		if !route.parsed {
			r.parseRoute(route)
		}
		r.routesByName[route.Name] = route
	}
}

/*
 */

func (r *Router) match(ctx *Ctx) bool {
	rq := ctx.Request
	if len(r.subdomainRegex) > 0 {
		subSplit := strings.Split(rq.Host, ".")
		if len(subSplit) != len(r.subdomainRegex) {
			return false
		}
		for i, s := range r.subdomainRegex {
			if !s.MatchString(subSplit[i]) { // regex do not work => create a new e replace in url too
				return false
			}
		}
	}

	for _, route := range r.Routes {
		if route.match(ctx) {
			ctx.MatchInfo.Router = r
			return true
		}
	}
	return false
}

/*
 */
func (r *Router) AddRoute(routes ...*Route) {
	r.Routes = append(r.Routes, routes...)
}

func (r *Router) Add(url, name string, f Func, meths []string) {
	r.AddRoute(
		&Route{
			Name:    name,
			Url:     url,
			Func:    f,
			Methods: meths,
		})
}

func (r *Router) GET(url string, f Func) { r.AddRoute(GET(url, f)) }

func (r *Router) PUT(url string, f Func) { r.AddRoute(PUT(url, f)) }

func (r *Router) HEAD(url string, f Func) { r.AddRoute(HEAD(url, f)) }

func (r *Router) POST(url string, f Func) { r.AddRoute(POST(url, f)) }

func (r *Router) TRACE(url string, f Func) { r.AddRoute(TRACE(url, f)) }

func (r *Router) PATCH(url string, f Func) { r.AddRoute(PATCH(url, f)) }

func (r *Router) DELETE(url string, f Func) { r.AddRoute(DELETE(url, f)) }

func (r *Router) CONNECT(url string, f Func) { r.AddRoute(CONNECT(url, f)) }

func (r *Router) OPTIONS(url string, f Func) { r.AddRoute(OPTIONS(url, f)) }
