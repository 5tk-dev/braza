package braza

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func serveFileHandler(ctx *Ctx) {
	rq := ctx.Request

	urlFilePath := rq.PathArgs["filepath"]
	pathToFile := filepath.Join(ctx.App.StaticFolder, urlFilePath)
	if f, err := os.Open(pathToFile); err == nil {
		_, file := filepath.Split(pathToFile)
		defer f.Close()
		if fStat, err := f.Stat(); err != nil || fStat.IsDir() {
			ctx.NotFound()
		}
		ctx.ReadFrom(f)
		ctype := mime.TypeByExtension(filepath.Ext(file))
		if ctype == "application/octet-stream" {
			ctype = http.DetectContentType(ctx.Bytes())
		}
		ctx.header.Set("Content-Type", ctype)
		ctx.Close()
	} else {
		if ctx.App.Env == "development" {
			ctx.TEXT(err, 404)
		}
		ctx.Response.NotFound()
	}
}

func optionsHandler(ctx *Ctx) {
	mi := ctx.MatchInfo
	rsp := ctx.Response
	rsp.StatusCode = 200
	strMeths := strings.Join(mi.Route.Cors.AllowMethods, ", ")
	if rsp.header.Get("Access-Control-Allow-Methods") == "" {
		rsp.header.Set("Access-Control-Allow-Methods", strMeths)
	}
	rsp.parseHeaders()
	SetHeader(rsp.raw, rsp.header)
}

func execTeardown(ctx *Ctx) {
	if ctx.App.TearDownRequest != nil {
		go ctx.App.TearDownRequest(ctx)
	}
}

func req500(ctx *Ctx) {
	defer l.LogRequest(ctx)
	if err := recover(); err != nil {
		statusText := "500 Internal Server Error"
		l.Error(err)
		ctx.raw.WriteHeader(500)
		fmt.Fprint(ctx.raw, statusText)
	}

}

func reqOK(ctx *Ctx) {
	mi := ctx.MatchInfo
	if mi.Match {
		if ctx.Session.changed {
			ctx.SetCookie(ctx.Session.save(ctx))
		}
		ctx.parseHeaders()
		SetHeader(ctx.raw, ctx.header)
	}
	ctx.raw.WriteHeader(ctx.StatusCode)
	fmt.Fprint(ctx.raw, ctx.String())
}
