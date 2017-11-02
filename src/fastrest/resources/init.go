package resources

import (
	"fastrest/logging"
	"fastrest/middleware"
	"fmt"
	"strings"
	"sync"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/xeipuuv/gojsonschema"
)

var once_load sync.Once
var create_user_loader gojsonschema.JSONLoader
var update_user_loader gojsonschema.JSONLoader
var create_role_loader gojsonschema.JSONLoader

func init() {
	once_load.Do(func() {
		create_user_loader = gojsonschema.NewStringLoader(
			user_create_json_schema)
		update_user_loader = gojsonschema.NewStringLoader(
			user_update_json_schema)
		create_role_loader = gojsonschema.NewStringLoader(
			role_create_json_schema)
	})
}

func valite_req_body(ctx *fasthttp.RequestCtx, str_body string,
	loader gojsonschema.JSONLoader) bool {
	req_body := gojsonschema.NewStringLoader(str_body)
	res, err := gojsonschema.Validate(loader, req_body)
	if err != nil {
		logging.ERROR.Printf("%v\n", err)
		ctx.Error("Only json-liked body accepted\n", fasthttp.StatusBadRequest)
		return false
	}

	if !res.Valid() {
		var errs []string
		errs = append(errs, fmt.Sprintln(
			"Input is invalid, following errors found:"))
		for _, desc := range res.Errors() {
			errs = append(errs, fmt.Sprintf("- %s", desc))
		}
		logging.ERROR.Printf("%v\n", strings.Join(errs, "\n"))
		ctx.Error(strings.Join(errs, "\n"), fasthttp.StatusBadRequest)
		return false
	}
	return true
}

func InitRouter(router *fasthttprouter.Router) {
	router.GET("/users", middleware.AuthMidddle(user_lists))
	router.GET("/users/:userid", middleware.AuthMidddle(user_get))
	router.POST("/users", middleware.BuildMiddleWareChain(
		user_create, middleware.JsonMiddleware, middleware.AuthMidddle))
	router.POST("/users/:userid", middleware.AuthMidddle(user_update))
	router.PUT("/users/:userid", middleware.AuthMidddle(user_update))
	router.DELETE("/users/:userid", middleware.AuthMidddle(user_delete))

	router.GET("/roles", middleware.AuthMidddle(role_lists))
	router.GET("/roles/:roleid", middleware.AuthMidddle(role_get))
	router.POST("/roles", middleware.AuthMidddle(role_create))
	router.DELETE("/roles/:roleid", middleware.AuthMidddle(role_delete))

	router.GET("/users/:userid/permisions",
		middleware.AuthMidddle(permision_list))
	router.PUT("/permisions/:userid/roles/:roleid",
		middleware.AuthMidddle(permision_create))
	router.DELETE("/permisions/:userid/roles/:roleid",
		middleware.AuthMidddle(permision_delete))
}
