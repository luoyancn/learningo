package resources

import (
	"encoding/json"
	"fmt"

	"fastrest/db"
	"fastrest/exceptions"
	"fastrest/logging"

	"github.com/valyala/fasthttp"
)

func user_lists(ctx *fasthttp.RequestCtx) {
	users, err := db.UserList()
	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logging.ERROR.Printf("%v\n", err)
		fmt.Fprintf(ctx, "Intelnal ERROR:%v\n", err)
		return
	}
	fmt.Fprintf(ctx, "%s\n", users)
}

func user_get(ctx *fasthttp.RequestCtx) {
	userid := ctx.UserValue("userid").(string)
	logging.DEBUG.Printf("Query the user with id %s\n", userid)
	user, err := db.UserGet(userid)
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
	fmt.Fprintf(ctx, "%s\n", user)
}

func user_create(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	if !valite_req_body(ctx, string(body), create_user_loader) {
		return
	}
	logging.DEBUG.Printf("Create user with %v\n", string(body))
	var user db.User
	err := json.Unmarshal(body, &user)
	if nil != err {
		logging.ERROR.Printf("%v\n", err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	err = db.UserCreate(user)
	if nil != err {
		logging.ERROR.Printf("%v\n", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Intelnal ERROR:%v\n", err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func user_update(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	if !valite_req_body(ctx, string(body), update_user_loader) {
		return
	}
	logging.DEBUG.Printf("Update user with entity %v\n", string(body))
	var user db.User
	err := json.Unmarshal(body, &user)
	if nil != err {
		logging.ERROR.Printf("%v\n", err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	userid := ctx.UserValue("userid").(string)
	err = db.UserUpdate(user, userid)
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
	ctx.SetStatusCode(fasthttp.StatusAccepted)
}

func user_delete(ctx *fasthttp.RequestCtx) {
	userid := ctx.UserValue("userid").(string)
	logging.DEBUG.Printf("Delete user with id %v\n", userid)
	err := db.UserDelete(userid)
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
