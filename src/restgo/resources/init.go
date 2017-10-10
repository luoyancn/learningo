package resources

import (
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	uuid "github.com/satori/go.uuid"
)

func check_valid_uuid(id string) bool {
	_, err := uuid.FromString(id)
	if nil != err {
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

	permision_mux := user_mux.NewGroup("/:userid/permision")
	permision_mux.GET("/", xhandler.HandlerFuncC(permision_list))
	permision_mux.GET("/:permisionid", xhandler.HandlerFuncC(permision_get))
}
