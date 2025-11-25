package braza

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

/*

SERVER funcs

*/

type ServeOptions struct {
	Listener net.Listener
	PubKey   string
	PrivKey  string
	Host     string
}

func (app *App) startListener(opt *ServeOptions, c chan error) {
	if opt.Host == "" {
		opt.Host = ":5000"
	}
	if app.Srv.Addr != "" {
		app.Srv.Addr = opt.Host
	}

	if opt.Listener != nil {
		c <- app.Srv.Serve(opt.Listener)
	} else if opt.PrivKey != "" && opt.PubKey != "" {
		c <- app.Srv.ListenAndServeTLS(opt.PrivKey, opt.PubKey)
	} else {
		c <- app.Srv.ListenAndServe()
	}
}

func (app *App) parseSrvApp(addr string) {
	app.Srv.Handler = app
	if app.Srv.MaxHeaderBytes == 0 {
		app.Srv.MaxHeaderBytes = 1 << 20
	}
	if app.Srv.Addr == "" && addr != "" {
		app.Srv.Addr = addr
	} else if app.Srv.Addr == "" {
		app.Srv.Addr = "localhost:5000"
	}

}

func runSrv(app *App, opt *ServeOptions) (err error) {
	if opt == nil {
		opt = &ServeOptions{}
	}
	app.Build(opt.Host)
	if listRouteSch != "" {
		showRouteSchema(app, listRouteSch)
	}
	var reboot = make(chan bool)
	var srvErr = make(chan error)

	if app.Env == "development" && !app.DisableFileWatcher {
		go fileWatcher(reboot)
	}

	if !app.Silent {
		app.logStarterListener()
	}

	go app.startListener(opt, srvErr)

	if !app.DisableFileWatcher {
		for {
			select {
			case <-reboot:
				app.Srv.Close()
				selfReboot()
			case err = <-srvErr:
				if !errors.Is(err, http.ErrServerClosed) {
					l.err.Println(err)
					return err
				}
			}
		}
	} else {
		e := <-srvErr
		l.err.Println(e)
		return e
	}
}

// Start HTTP Server
func (app *App) Listen(opt *ServeOptions) (err error) { return runSrv(app, opt) }

/*

HTTP funcs

*/

func (app *App) match(ctx *Ctx) {
	rq := ctx.Request
	if app.hostname != "" {
		if net.ParseIP(rq.Host) != nil {
			ctx.NotFound()
		}

		if !strings.Contains(rq.Host, app.hostname) {
			ctx.NotFound()
		}
	}
	for _, router := range app.routers {
		if router.match(ctx) {
			if ctx.MatchInfo.MethodNotAllowed != nil {
				l.warn.Printf("url match err : Method Not Allowed")
			}
			if router.StrictSlash && !strings.HasSuffix(rq.URL.Path, "/") {
				args := []string{}
				for k, v := range rq.PathArgs {
					args = append(args, k, v)
				}
				ctx.Response.Redirect(ctx.UrlFor(ctx.MatchInfo.Route.Name, true, args...))
			}
			return
		}
	}

	mi := ctx.MatchInfo
	if mi.MethodNotAllowed != nil {
		ctx.MethodNotAllowed()
	}
	ctx.NotFound()
}

// exec route and handle errors of application
func (app *App) execRoute(ctx *Ctx) {
	app.match(ctx)
	rq := ctx.Request
	mi := ctx.MatchInfo

	rq.parse()
	ctx.parseMids()

	if app.BeforeRequest != nil {
		app.BeforeRequest(ctx)
	}

	if mi.Func == nil && rq.Method == "OPTIONS" {
		optionsHandler(ctx)
		return
	}
	ctx.next()
}

func (app *App) execHandlerError(ctx *Ctx, code int) {
	ctx.Reset()
	if h, ok := app.errHandlers[code]; ok {
		ctx.StatusCode = code
		h(ctx)
	} else {
		ctx.StatusCode = code
		statusText := http.StatusText(code)
		body := fmt.Sprintf("%d %s", code, statusText)
		ctx.header.Set("Content-Type", "text/plain")
		ctx.WriteString(body)
	}
}

func (app *App) closeConn(ctx *Ctx) {
	defer req500(ctx)
	defer execTeardown(ctx)

	err := recover()
	if err == nil {
		reqOK(ctx)
		return
	}

	if e, ok := err.(error); ok && errors.Is(ErrHttpAbort, e) {
		code := ctx.backCtx.Value(abortCode(1))
		if c, ok := code.(int); ok {
			app.execHandlerError(ctx, c)
		}
		reqOK(ctx)
	}

}

// http.Handler
func (app *App) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	ctx := NewCtx(app, wr, req)
	defer app.closeConn(ctx)
	app.execRoute(ctx)
}
