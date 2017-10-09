package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/xmux"
)

func UserLists(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	fmt.Fprintf(respwriter, "List all of users !!!\n")
}
func UserGet(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	if !check_valid_uuid(xmux.Param(ctx, "userid")) {
		respwriter.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(respwriter, "Invalid format of userid: need uuid format\n")
		return
	}
}

func UserCreate(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	fmt.Fprintf(respwriter, "Get the user which id equals %s !!!\n",
		xmux.Param(ctx, "userid"))
}

func UserUpdate(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	fmt.Fprintf(respwriter, "Get the user which id equals %s !!!\n",
		xmux.Param(ctx, "userid"))
}

func UserDelete(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	fmt.Fprintf(respwriter, "Get the user which id equals %s !!!\n",
		xmux.Param(ctx, "userid"))
}
