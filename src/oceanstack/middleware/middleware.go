package middleware

import (
	"oceanstack/conf"
	"oceanstack/db/redisdb"
	"oceanstack/logging"

	"github.com/valyala/fasthttp"
)

type middleware func(next fasthttp.RequestHandler) fasthttp.RequestHandler

func AuthMidddle(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		auth_token := string(ctx.Request.Header.Peek("X-Auth-Token"))
		if "" == auth_token {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized),
				fasthttp.StatusUnauthorized)
			return
		}
		logging.LOG.Debugf("The token is :%s\n", auth_token)
		if conf.ADMIN_TOKEN != auth_token && !redisdb.ValidToken(auth_token) {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusForbidden),
				fasthttp.StatusForbidden)
			return
		}
		next(ctx)
	})
}

func JsonMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		if "application/json" != string(
			ctx.Request.Header.Peek("Content-Type")) {
			ctx.Error(
				"Only application/json accepted\n",
				fasthttp.StatusUnsupportedMediaType)
			return
		}
		next(ctx)
	})
}

func JsonResponseMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		next(ctx)
		ctx.SetContentType("application/json")
	})
}

func BuildPipeLine(app fasthttp.RequestHandler,
	mid ...middleware) fasthttp.RequestHandler {
	if 0 == len(mid) {
		return app
	}
	return mid[0](BuildPipeLine(app, mid[1:cap(mid)]...))
}
