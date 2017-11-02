package middleware

import (
	"github.com/valyala/fasthttp"
)

func AuthMidddle(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		auth_token := string(ctx.Request.Header.Peek("X-Auth-Token"))
		if "" == auth_token {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized),
				fasthttp.StatusUnauthorized)
			return
		}
		if "123" != auth_token {
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

type middleware func(next fasthttp.RequestHandler) fasthttp.RequestHandler

func BuildMiddleWareChain(app fasthttp.RequestHandler,
	mid ...middleware) fasthttp.RequestHandler {
	if 0 == len(mid) {
		return app
	}
	return mid[0](BuildMiddleWareChain(app, mid[1:cap(mid)]...))
}
