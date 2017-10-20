package resources

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	uuid "github.com/satori/go.uuid"
	"github.com/xeipuuv/gojsonschema"
)

func check_valid_uuid(id string) bool {
	_, err := uuid.FromString(id)
	if nil != err {
		return false
	}
	return true
}

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

func valite_req_body(writer http.ResponseWriter, str_body string,
	loader gojsonschema.JSONLoader) bool {
	req_body := gojsonschema.NewStringLoader(str_body)
	res, err := gojsonschema.Validate(loader, req_body)
	if err != nil {
		http.Error(writer, "Only json-liked body accepted\n",
			http.StatusBadRequest)
		return false
	}

	if !res.Valid() {
		var errs []string
		errs = append(errs, fmt.Sprintln(
			"Input is invalid, following errors found:"))
		for _, desc := range res.Errors() {
			errs = append(errs, fmt.Sprintf("- %s", desc))
		}
		http.Error(writer, strings.Join(errs, "\n"),
			http.StatusBadRequest)
		return false
	}
	return true
}

func InitRouter(root_mux *xmux.Mux) {
	user_mux := root_mux.NewGroup("/users")
	user_mux.GET("/", xhandler.HandlerFuncC(user_lists))
	user_mux.GET("/:userid", xhandler.HandlerFuncC(user_get))
	user_mux.POST("/", xhandler.HandlerFuncC(user_create))
	user_mux.POST("/:userid", xhandler.HandlerFuncC(user_update))
	user_mux.PUT("/:userid", xhandler.HandlerFuncC(user_update))
	user_mux.DELETE("/:userid", xhandler.HandlerFuncC(user_delete))

	role_mux := root_mux.NewGroup("/roles")
	role_mux.GET("/", xhandler.HandlerFuncC(role_lists))
	role_mux.GET("/:roleid", xhandler.HandlerFuncC(role_get))
	role_mux.POST("/", xhandler.HandlerFuncC(role_create))
	role_mux.DELETE("/:roleid", xhandler.HandlerFuncC(role_delete))

	root_mux.GET("/users/:userid/permisions",
		xhandler.HandlerFuncC(permision_list))
	root_mux.PUT("/permisions/:userid/roles/:roleid",
		xhandler.HandlerFuncC(permision_create))
	root_mux.DELETE("/permisions/:userid/roles/:roleid",
		xhandler.HandlerFuncC(permision_delete))
}
