package api

import (
	"fmt"
	"oceanstack/exceptions"

	"github.com/valyala/fasthttp"
)

func authentication(ctx *fasthttp.RequestCtx) {
	if err := valite_req_body(
		string(ctx.PostBody()), auth_loader); nil != err {
		switch err.(type) {
		case exceptions.JsonMarshallException:
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			fmt.Fprintf(ctx, "ERROR: %v\n", err)
			return
		default:
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			fmt.Fprintf(ctx, "ERROR: %v\n", err)
			return
		}
	}
}
