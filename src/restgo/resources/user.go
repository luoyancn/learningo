package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"restgo/db"
	"restgo/exceptions"

	"github.com/rs/xmux"
)

func user_lists(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	users, err := db.UserList()
	if nil != err {
		respwriter.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(respwriter, "Intelnal ERROR:%v\n", err)
		return
	}
	fmt.Fprintf(respwriter, "%s\n", users)
}

func user_get(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	userid := xmux.Param(ctx, "userid")
	user, err := db.UserGet(userid)
	if nil != err {
		switch err.(type) {
		case exceptions.NotFoundException:
			respwriter.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(respwriter, "ERROR:%v\n", err)
			return
		default:
			respwriter.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(respwriter, "Intelnal ERROR:%v\n", err)
			return
		}
	}
	fmt.Fprintf(respwriter, "%s\n", user)
}

func user_create(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(respwriter, "Failed to parse the body in request\n",
			http.StatusBadRequest)
		return
	}
	if !valite_req_body(respwriter, string(body), create_user_loader) {
		return
	}
	var user db.User
	err = json.Unmarshal(body, &user)
	if nil != err {
		http.Error(respwriter, err.Error(), http.StatusInternalServerError)
		return
	}
	err = db.UserCreate(user)
	if nil != err {
		respwriter.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(respwriter, "Intelnal ERROR:%v\n", err)
		return
	}
	respwriter.WriteHeader(http.StatusCreated)
}

func user_update(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(respwriter, "Failed to parse the body in request\n",
			http.StatusBadRequest)
		return
	}
	if !valite_req_body(respwriter, string(body), update_user_loader) {
		return
	}
	var user db.User
	err = json.Unmarshal(body, &user)
	if nil != err {
		http.Error(respwriter, err.Error(), http.StatusInternalServerError)
		return
	}
	userid := xmux.Param(ctx, "userid")
	err = db.UserUpdate(user, userid)
	if nil != err {
		switch err.(type) {
		case exceptions.NotFoundException:
			respwriter.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(respwriter, "ERROR:%v\n", err)
			return
		default:
			respwriter.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(respwriter, "Intelnal ERROR:%v\n", err)
			return
		}
	}
	respwriter.WriteHeader(http.StatusAccepted)
}

func user_delete(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	userid := xmux.Param(ctx, "userid")
	err := db.UserDelete(userid)
	if nil != err {
		switch err.(type) {
		case exceptions.NotFoundException:
			respwriter.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(respwriter, "ERROR:%v\n", err)
			return
		default:
			respwriter.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(respwriter, "Intelnal ERROR:%v\n", err)
			return
		}
	}
	respwriter.WriteHeader(http.StatusNoContent)
}
