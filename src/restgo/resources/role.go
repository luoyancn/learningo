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

func role_lists(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	roles, err := db.RoleList()
	if nil != err {
		respwriter.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(respwriter, "Intelnal ERROR:%v\n", err)
		return
	}
	fmt.Fprintf(respwriter, "%s\n", roles)
}

func role_get(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	roleid := xmux.Param(ctx, "roleid")
	role, err := db.RoleGet(roleid)
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
	fmt.Fprintf(respwriter, "%s\n", role)
}

func role_create(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(respwriter, "Failed to parse the body in request\n",
			http.StatusBadRequest)
		return
	}
	if !valite_req_body(respwriter, string(body), create_role_loader) {
		return
	}
	var role db.Role
	err = json.Unmarshal(body, &role)
	if nil != err {
		http.Error(respwriter, err.Error(), http.StatusInternalServerError)
		return
	}
	err = db.RoleCreate(role)
	if nil != err {
		respwriter.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(respwriter, "Intelnal ERROR:%v\n", err)
		return
	}
	respwriter.WriteHeader(http.StatusCreated)
}

func role_delete(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	roleid := xmux.Param(ctx, "roleid")
	err := db.RoleDelete(roleid)
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
