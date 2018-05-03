package api

import (
	"github.com/valyala/fasthttp"
)

func authentication(ctx *fasthttp.RequestCtx) {
	body := string(ctx.PostBody())
	if err := valite_req_body(body, auth_loader, ctx); nil != err {
		return
	}
}
