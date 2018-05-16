package api

import (
	"encoding/json"
	"fmt"
	"oceanstack/db"
	"oceanstack/utils"

	"github.com/valyala/fasthttp"
)

func user_list(ctx *fasthttp.RequestCtx) {
}

func user_get(ctx *fasthttp.RequestCtx) {
}

func user_create(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	err := valite_req_body(string(body), user_create_loader, ctx)
	if nil != err {
		return
	}
	var user_map map[string]db.User
	_ = json.Unmarshal(body, &user_map)
	user := user_map["user"]
	user.Password = utils.Md5Crypto(user.Password)
	uuid, err := db.UserCreate(user)
	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(
			ctx, "Failed to create user %s:%s\n", user.Name, err.Error())
		return
	}
	resp, _ := json.Marshal(
		map[string]string{"uuid": uuid, "name": user.Name})
	fmt.Fprintf(ctx, "%s\n", string(resp))
	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func user_update(ctx *fasthttp.RequestCtx) {
}

func user_delete(ctx *fasthttp.RequestCtx) {
}
