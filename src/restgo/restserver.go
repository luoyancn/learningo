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
var rootMux *xmux.Mux

func init() {
	once.Do(func() {
		rootMux = xmux.New()
		rootMux.GET("/", xhandler.HandlerFuncC(root))
		resources.InitRouter(rootMux)
	})
}

func root(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	fmt.Fprintf(respwriter, "Welcome to the rest world of go !!!\n")
}

func Serve() {
	log.Fatal(http.ListenAndServe(
		":8080", xhandler.New(context.Background(), rootMux)))
}
