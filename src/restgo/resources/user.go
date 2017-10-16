package resources

import (
	"context"
	"fmt"
	"net/http"

	"restgo/db"

	"github.com/rs/xmux"
)

func user_lists(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	db.UserList()
	fmt.Fprintf(respwriter, "List all of users !!!\n")
}

func user_get(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	if !check_valid_uuid(xmux.Param(ctx, "userid")) {
		respwriter.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(respwriter, "Invalid format of userid: need uuid format\n")
		return
	}
}

func user_create(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	fmt.Fprintf(respwriter, "Get the user which id equals %s !!!\n",
		xmux.Param(ctx, "userid"))
}

func user_update(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	fmt.Fprintf(respwriter, "Get the user which id equals %s !!!\n",
		xmux.Param(ctx, "userid"))
}

func user_delete(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	fmt.Fprintf(respwriter, "Get the user which id equals %s !!!\n",
		xmux.Param(ctx, "userid"))
}
