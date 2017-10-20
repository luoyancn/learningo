package resources

import (
	"context"
	"fmt"
	"net/http"
	"restgo/db"
	"restgo/exceptions"

	"github.com/rs/xmux"
)

func permision_list(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	userid := xmux.Param(ctx, "userid")
	permisions, err := db.AssgnmentList(userid)
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
	fmt.Fprintf(respwriter, "%s\n", permisions)
}

func permision_create(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	userid := xmux.Param(ctx, "userid")
	roleid := xmux.Param(ctx, "roleid")
	err := db.AssgnmentCreate(db.Assignment{UserUuId: userid, RoleUuId: roleid})
	if nil != err {
		respwriter.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(respwriter, "Intelnal ERROR:%v\n", err)
		return
	}
	respwriter.WriteHeader(http.StatusCreated)
}

func permision_delete(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	userid := xmux.Param(ctx, "userid")
	roleid := xmux.Param(ctx, "roleid")
	err := db.AssgnmentDelete(userid, roleid)
	if nil != err {
		respwriter.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(respwriter, "Intelnal ERROR:%v\n", err)
		return
	}
	respwriter.WriteHeader(http.StatusNoContent)
}
