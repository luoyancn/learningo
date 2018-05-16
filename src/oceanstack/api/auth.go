package api

import (
	"encoding/json"
	"fmt"
	"oceanstack/db"
	"oceanstack/db/redisdb"
	"oceanstack/exceptions"
	"oceanstack/logging"
	"oceanstack/utils"

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

	user_pointer, err := db.UserGet(
		auth_map["auth"].Name, utils.Md5Crypto(auth_map["auth"].Password))
	if nil != err {
		switch err.(type) {
		case exceptions.NotFoundException:
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString(
				"Unauthorized:The request requires authentication\n")
			return
		case exceptions.SQLException:
		case exceptions.ConnectionException:
		case exceptions.Exception:
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			fmt.Fprintf(ctx, "ERROR: %v\n", err)
			return
		}
	}
	token, err := redisdb.TokenIssue(user_pointer.Uuid)
	if nil != err {
		logging.LOG.Errorf("Cannot generate auth token:%v\n", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "ERROR: Failed to generate auth tokens\n")
		return
	}
	resp, _ := json.Marshal(
		map[string]string{"token": token, "name": user_pointer.Name})
	fmt.Fprintf(ctx, "%s\n", string(resp))
	ctx.SetStatusCode(fasthttp.StatusAccepted)
}
