package resources

import (
	"fastrest/logging"
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
	router.GET("/users", user_lists)
	router.GET("/users/:userid", user_get)
	router.POST("/users", user_create)
	router.POST("/users/:userid", user_update)
	router.PUT("/users/:userid", user_update)
	router.DELETE("/users/:userid", user_delete)

	router.GET("/roles", role_lists)
	router.GET("/roles/:roleid", role_get)
	router.POST("/roles", role_create)
	router.DELETE("/roles/:roleid", role_delete)

	router.GET("/users/:userid/permisions", permision_list)
	router.PUT("/permisions/:userid/roles/:roleid", permision_create)
	router.DELETE("/permisions/:userid/roles/:roleid", permision_delete)
}
