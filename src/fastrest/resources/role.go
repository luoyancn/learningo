package resources

import (
	"encoding/json"
	"fastrest/db"
	"fastrest/exceptions"
	"fastrest/logging"
	"fmt"

	"github.com/valyala/fasthttp"
)

func role_lists(ctx *fasthttp.RequestCtx) {
	roles, err := db.RoleList()
	if nil != err {
		logging.ERROR.Printf("%v\n", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Intelnal ERROR:%v\n", err)
		return
	}
	fmt.Fprintf(ctx, "%s\n", roles)
}

func role_get(ctx *fasthttp.RequestCtx) {
	roleid := ctx.UserValue("roleid").(string)
	logging.DEBUG.Printf("Get the role with id %s\n", roleid)
	role, err := db.RoleGet(roleid)
	if nil != err {
		logging.ERROR.Printf("%v\n", err)
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
	logging.DEBUG.Printf("Create the role with entity %v\n", string(body))
	var role db.Role
	err := json.Unmarshal(body, &role)
	if nil != err {
		logging.ERROR.Printf("%v\n", err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	err = db.RoleCreate(role)
	if nil != err {
		logging.ERROR.Printf("%v\n", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Intelnal ERROR:%v\n", err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func role_delete(ctx *fasthttp.RequestCtx) {
	roleid := ctx.UserValue("roleid").(string)
	logging.DEBUG.Printf("Delete the role with id %s\n", roleid)
	err := db.RoleDelete(roleid)
	if nil != err {
		logging.ERROR.Printf("%v\n", err)
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
