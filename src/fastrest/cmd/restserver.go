package restgo

import (
	_ "fastrest/conf"
	"fastrest/db"
	"fastrest/logging"
	"fastrest/resources"
	"fmt"
	"os"
	"sync"

	"github.com/buaazp/fasthttprouter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

var once sync.Once
var router *fasthttprouter.Router

var RootCmd = &cobra.Command{
	Short: "The commands of rest application of golang",
	Long:  ` The commands to start rest server or init the database`,
}

var migrateCmd = &cobra.Command{
	Use:   "db_sync",
	Short: "Init the project`s database",
	Long:  ` Create database or update the schema of database.`,
	Run:   Sync,
}

var serveCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the rest server",
	Long:  ` Start the rest server for handle the http request`,
	Run:   Serve,
}

var configfile string

func init() {
	once.Do(func() {
		RootCmd.AddCommand(migrateCmd)
		RootCmd.AddCommand(serveCmd)
		RootCmd.PersistentFlags().StringVarP(
			&configfile, "config", "c", "", "The full path of config file")

		router = fasthttprouter.New()
		router.GET("/", root)
		resources.InitRouter(router)
	})
}

func root(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Welcome to the rest world of go !!!\n")
}

func read_config() {
	viper.SetConfigFile(configfile)
	if err := viper.ReadInConfig(); nil != err {
		logging.WARNING.Printf("ERROR:%v\n", err)
		logging.WARNING.Printf(
			"And we will use default values of all config params!\n")
	}
	logging.GetLogger()
}

func Serve(cmd *cobra.Command, args []string) {
	read_config()
	db.InitDbConnection()
	fasthttp.ListenAndServe(
		viper.GetString("default.listen"), router.Handler)
}

func Sync(cmd *cobra.Command, args []string) {
	read_config()
	db.MigrateDB()
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
