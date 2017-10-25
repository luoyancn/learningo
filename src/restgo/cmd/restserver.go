package restgo

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "restgo/conf"
	"restgo/db"
	"restgo/resources"
	"sync"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var once sync.Once
var rootMux *xmux.Mux

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
		rootMux = xmux.New()
		rootMux.GET("/", xhandler.HandlerFuncC(root))
		resources.InitRouter(rootMux)
	})
}

func root(ctx context.Context, respwriter http.ResponseWriter,
	req *http.Request) {
	fmt.Fprintf(respwriter, "Welcome to the rest world of go !!!\n")
}

func read_config() {
	viper.SetConfigFile(configfile)
	if err := viper.ReadInConfig(); nil != err {
		fmt.Printf("%v\n", err)
		fmt.Println("And we will use default values of all config params!")
	}
}

func Serve(cmd *cobra.Command, args []string) {
	read_config()
	db.InitDbConnection()
	log.Fatal(http.ListenAndServe(viper.GetString("default.listen"),
		xhandler.New(context.Background(), rootMux)))
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
