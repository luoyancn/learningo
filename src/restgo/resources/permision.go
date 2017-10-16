package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/xmux"
)

func permision_list(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	fmt.Fprintf(
		respwriter, "Get the user`s permision which id equals %s !!!\n",
		xmux.Param(ctx, "userid"))
}

func permision_get(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	fmt.Fprintf(
		respwriter, "The user`s permision which id equals %s is %s !!!\n",
		xmux.Param(ctx, "userid"), xmux.Param(ctx, "permisionid"))
}