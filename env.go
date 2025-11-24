package braza

import (
	"flag"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	allowEnv = map[string]string{
		"":            "development",
		"d":           "development",
		"dev":         "development",
		"development": "development",

		"t":       "test",
		"test":    "test",
		"testing": "test",

		"p":          "production",
		"prod":       "production",
		"production": "production",
	}
	l            *logger
	listenAll    bool
	localAddress = getOutboundIP()
	mapStackApps = map[string]*App{}
)

var (
	env          string
	port         string
	address      string
	listRoutes   bool
	listRouteSch string
)

func init() {
	flag.StringVar(&env, "env", "", "set a address listener")
	flag.StringVar(&port, "port", "", "set a address listener")
	flag.StringVar(&address, "address", "", "set a address listener")
	flag.BoolVar(&listRoutes, "routes", false, "list all routes")
	flag.StringVar(&listRouteSch, "routeSch", "", "show routes schema-> app.route || app.route:GET")
}

/*
ENV funcs
*/

func (app *App) setFlags() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if app.Srv == nil {
		app.Srv = &http.Server{}
	}
	if address != "" {
		app.Srv.Addr = address
	}
	if port != "" {
		if !re.httpPort.MatchString(port) {
			l.err.Panicf("port '%s' is not valid!", port)
		}
		port = strings.TrimPrefix(port, ":")
		if app.Srv.Addr != "" {
			h, p, err := net.SplitHostPort(app.Srv.Addr)
			if h == "" && p == "" && err != nil {
				l.err.Panic(err)
			}
			app.Srv.Addr = net.JoinHostPort(h, port)
		} else {
			app.Srv.Addr = ":" + port
		}
	}

	if listRoutes {
		app.ShowRoutes()
	}
	if listRouteSch != "" {
		showRouteSchema(app, listRouteSch)
	}
}

func (app *App) setEnv() {
	if e := os.Getenv("ENV"); e != "" {
		app.Env = e
	}
	if s := os.Getenv("SERVERNAME"); s != "" {
		app.Servername = s
	}
	if addr := os.Getenv("ADDRESS"); addr != "" {
		if app.Srv == nil {
			app.Srv = &http.Server{}
		}
		app.Srv.Addr = addr
	}
}

func (app *App) logStarterListener() {
	addr, port, err := net.SplitHostPort(app.Srv.Addr)
	if err != nil {
		l.err.Panic(err)
	}
	envDev := app.Env == "development"
	if listenAll {
		app.Srv.Addr = localAddress
		if envDev {
			l.Defaultf("Server is listening on all address in %sdevelopment mode%s", _RED, _RESET)
		} else {
			l.Default("Server is listening on all address")
		}
		l.info.Printf("          listening on: http://%s:%s", getOutboundIP(), port)
		l.info.Printf("          listening on: http://0.0.0.0:%s", port)
	} else {
		if envDev {
			l.Defaultf("Server is listening in %sdevelopment mode%s", _RED, _RESET)
		} else {
			l.Default("Server is linsten")
		}
		if addr == "" {
			addr = "localhost"
		}
		l.info.Printf("          listening on: http://%s:%s", addr, port)
	}
	if app.Servername != "" {
		l.info.Printf("          setting servername: '%s'", app.Servername)
	}
}
