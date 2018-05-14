package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"oceanstack/db"
	"oceanstack/exceptions"

	"github.com/valyala/fasthttp"
)

func authentication(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	err := valite_req_body(string(body), auth_loader, ctx)
	if nil != err {
		return
	}
	var auth_map map[string]db.User
	_ = json.Unmarshal(body, &auth_map)

	md5_writer := md5.New()
	io.WriteString(md5_writer, auth_map["auth"].Password)
	_, err = db.UserGet(
		auth_map["auth"].Name, hex.EncodeToString(md5_writer.Sum(nil)))
	if nil != err {
		switch err.(type) {
		case exceptions.NotFoundException:
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString(
				"Unauthorized:The request requires authentication\n")
			return
		case exceptions.SQLException:
		case exceptions.JsonMarshallException:
		case exceptions.Exception:
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			fmt.Fprintf(ctx, "ERROR: %v\n", err)
			return
		}
	}
	ctx.SetStatusCode(fasthttp.StatusAccepted)
}
