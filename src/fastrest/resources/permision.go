package resources

import (
	"fastrest/db"
	"fastrest/exceptions"
	"fastrest/logging"
	"fmt"

	"github.com/valyala/fasthttp"
)

func permision_list(ctx *fasthttp.RequestCtx) {
	userid := ctx.UserValue("userid").(string)
	logging.DEBUG.Printf("Get the user permistion with userid %s\n", userid)
	permisions, err := db.AssgnmentList(userid)
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
	fmt.Fprintf(ctx, "%s\n", permisions)
}

func permision_create(ctx *fasthttp.RequestCtx) {
	userid := ctx.UserValue("userid").(string)
	roleid := ctx.UserValue("roleid").(string)
	logging.DEBUG.Printf("Create permision with userid %s and roleid %s\n",
		userid, roleid)
	err := db.AssgnmentCreate(
		db.Assignment{UserUuId: userid, RoleUuId: roleid})
	if nil != err {
		logging.ERROR.Printf("%v\n", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Intelnal ERROR:%v\n", err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func permision_delete(ctx *fasthttp.RequestCtx) {
	userid := ctx.UserValue("userid").(string)
	roleid := ctx.UserValue("roleid").(string)
	logging.DEBUG.Printf("Delete permision with userid %s and roleid %s\n",
		userid, roleid)
	err := db.AssgnmentDelete(userid, roleid)
	if nil != err {
		logging.ERROR.Printf("%v\n", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Intelnal ERROR:%v\n", err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusNoContent)
}
