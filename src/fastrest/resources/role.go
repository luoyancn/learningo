package resources

import (
	"encoding/json"
	"fastrest/db"
	"fastrest/exceptions"
	"fmt"

	"github.com/valyala/fasthttp"
)

func role_lists(ctx *fasthttp.RequestCtx) {
	roles, err := db.RoleList()
	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Intelnal ERROR:%v\n", err)
		return
	}
	fmt.Fprintf(ctx, "%s\n", roles)
}

func role_get(ctx *fasthttp.RequestCtx) {
	roleid := ctx.UserValue("roleid").(string)
	role, err := db.RoleGet(roleid)
	if nil != err {
		switch err.(type) {
		case exceptions.NotFoundException:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			fmt.Fprintf(ctx, "ERROR:%v\n", err)
			return
		default:
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			fmt.Fprintf(ctx, "Intelnal ERROR:%v\n", err)
			return
		}
	}
	fmt.Fprintf(ctx, "%s\n", role)
}

func role_create(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	if !valite_req_body(ctx, string(body), create_role_loader) {
		return
	}
	var role db.Role
	err := json.Unmarshal(body, &role)
	if nil != err {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	err = db.RoleCreate(role)
	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Intelnal ERROR:%v\n", err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func role_delete(ctx *fasthttp.RequestCtx) {
	roleid := ctx.UserValue("roleid").(string)
	err := db.RoleDelete(roleid)
	if nil != err {
		switch err.(type) {
		case exceptions.NotFoundException:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			fmt.Fprintf(ctx, "ERROR:%v\n", err)
			return
		default:
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			fmt.Fprintf(ctx, "Intelnal ERROR:%v\n", err)
			return
		}
	}
	ctx.SetStatusCode(fasthttp.StatusNoContent)
}
