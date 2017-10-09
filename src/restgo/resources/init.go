package resources

import (
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	uuid "github.com/satori/go.uuid"
)

func InitRouter(root_mux *xmux.Mux) {
	user_mux := root_mux.NewGroup("/users")
	user_mux.GET("/", xhandler.HandlerFuncC(UserLists))
	user_mux.GET("/:userid", xhandler.HandlerFuncC(UserGet))
	user_mux.POST("/", xhandler.HandlerFuncC(UserCreate))
	user_mux.POST("/:userid", xhandler.HandlerFuncC(UserUpdate))
	user_mux.PUT("/:userid", xhandler.HandlerFuncC(UserUpdate))
	user_mux.DELETE("/:userid", xhandler.HandlerFuncC(UserDelete))

	permision_mux := user_mux.NewGroup("/:userid/permision")
	permision_mux.GET("/", xhandler.HandlerFuncC(PermisionList))
	permision_mux.GET("/:permisionid", xhandler.HandlerFuncC(PermisionGet))
}

func check_valid_uuid(id string) bool {
	_, err := uuid.FromString(id)
	if nil != err {
		return false
	}
	return true
}
