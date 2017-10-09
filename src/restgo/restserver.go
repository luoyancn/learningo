package restgo

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restgo/resources"
	"sync"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

var once sync.Once
var RootMux *xmux.Mux

func init() {
	once.Do(func() {
		RootMux = xmux.New()
		RootMux.GET("/", xhandler.HandlerFuncC(Root))
		resources.InitRouter(RootMux)
	})
}

func Root(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	fmt.Fprintf(respwriter, "Welcome to the rest world of go !!!\n")
}

func Serve() {
	log.Fatal(http.ListenAndServe(
		":8080", xhandler.New(context.Background(), RootMux)))
}
