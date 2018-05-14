package api

import (
	"fmt"
	"oceanstack/api/schemas"
	"oceanstack/conf"
	"oceanstack/exceptions"
	"oceanstack/logging"
	"oceanstack/middleware"
	"strings"
	"sync"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/xeipuuv/gojsonschema"
)

var router *fasthttprouter.Router

var once_load sync.Once
var auth_loader gojsonschema.JSONLoader
var user_create_loader gojsonschema.JSONLoader

func init() {
	once_load.Do(func() {
		auth_loader = gojsonschema.NewStringLoader(schemas.AUTH_JSON_SCHEMA)
		user_create_loader = gojsonschema.NewStringLoader(
			schemas.USER_CREATE_JSON_SCHEMA)
	})
}

func valite_req_body(str_body string, loader gojsonschema.JSONLoader,
	ctx *fasthttp.RequestCtx) error {
	req_body := gojsonschema.NewStringLoader(str_body)
	res, err := gojsonschema.Validate(loader, req_body)
	if err != nil {
		logging.LOG.Errorf("Cannot convert request body to json:%v\n", err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprintf(ctx, "EEROR: %v\n", err)
		return exceptions.NewJsonMarshallException(err.Error())
	}

	if !res.Valid() {
		var errs []string
		errs = append(errs, fmt.Sprintln(
			"Input is invalid, following errors found:"))
		for _, desc := range res.Errors() {
			errs = append(errs, fmt.Sprintf("- %s\n", desc))
		}
		var err_msg string
		err_msg = strings.Join(errs, "")
		logging.LOG.Errorf("Invalid Json Request body :%v\n", err_msg)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprintf(ctx, "Internal ERROR: %v\n", err_msg)
		return exceptions.NewInvalidJsonException(err_msg)
	}
	return nil
}

func init() {
	router = fasthttprouter.New()
	router.GET("/", root)
	router.POST("/auth", middleware.BuildPipeLine(
		authentication, middleware.JsonMiddleware))

	router.GET("/users", middleware.BuildPipeLine(
		user_list, middleware.AuthMidddle))
	router.GET("/users/:userid", middleware.BuildPipeLine(
		user_get, middleware.AuthMidddle))
	router.POST("/users", middleware.BuildPipeLine(
		user_create, middleware.JsonMiddleware, middleware.AuthMidddle))
	router.PUT("/users/:userid", middleware.BuildPipeLine(
		user_update, middleware.JsonMiddleware, middleware.AuthMidddle))
	router.DELETE("/users/:userid", middleware.BuildPipeLine(
		user_delete, middleware.AuthMidddle))
}

func root(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Welcome to the world of Ocean Stack !!!\n")
}

func Serve() {
	fasthttp.ListenAndServe(conf.LISTEN, router.Handler)
}
