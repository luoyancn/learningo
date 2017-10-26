package resources

import (
	"fastrest/db"
	"fastrest/exceptions"
	"fmt"

	"github.com/valyala/fasthttp"
)

func permision_list(ctx *fasthttp.RequestCtx) {
	userid := ctx.UserValue("userid").(string)
	permisions, err := db.AssgnmentList(userid)
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
	fmt.Fprintf(ctx, "%s\n", permisions)
}

func permision_create(ctx *fasthttp.RequestCtx) {
	userid := ctx.UserValue("userid").(string)
	roleid := ctx.UserValue("roleid").(string)
	err := db.AssgnmentCreate(db.Assignment{UserUuId: userid, RoleUuId: roleid})
	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Intelnal ERROR:%v\n", err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func permision_delete(ctx *fasthttp.RequestCtx) {
	userid := ctx.UserValue("userid").(string)
	roleid := ctx.UserValue("roleid").(string)
	err := db.AssgnmentDelete(userid, roleid)
	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Intelnal ERROR:%v\n", err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusNoContent)
}
