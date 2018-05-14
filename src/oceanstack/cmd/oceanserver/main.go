package main

import (
	"oceanstack/api"
	"oceanstack/common"
	"oceanstack/conf"
	"oceanstack/db"
	"oceanstack/db/redisdb"
	"oceanstack/logging"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
)

var once sync.Once
var configfile string

var rootcmd = &cobra.Command{
	Short: "Server of ocean stack",
	Long:  ` This is the binary of ocean server`,
}

var startcmd = &cobra.Command{
	Use:   "start",
	Short: "Start Ocean server",
	Long:  ` Start Ocean server`,
	Run:   serve,
	Args:  cobra.NoArgs,
}

var vercmd = &cobra.Command{
	Use:   "version",
	Short: "Get the ocean server version",
	Long:  ` Get the version of ocean server binary`,
	Run:   get_version,
	Args:  cobra.NoArgs,
}

func init() {
	once.Do(func() {
		startcmd.PersistentFlags().StringVarP(
			&configfile, "config-file", "c", "",
			"The full path of config file (Required, Only yaml, json and toml supported)")
		startcmd.MarkPersistentFlagRequired("config-file")
		rootcmd.AddCommand(startcmd)
		rootcmd.AddCommand(vercmd)
	})
}

func serve(cmd *cobra.Command, args []string) {
	common.ReadConfig(configfile, "oceanserver", logging.FILE_ENABLED)
	db.InitDbConnection()
	redisdb.InitRedisConnection()
	logging.LOG.Infof("Ocean Server started, and listen on %s\n", conf.LISTEN)
	go stop()
	api.Serve()
}

func stop() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case s := <-sig:
		logging.LOG.Infof("Exit Ocean Server: Recived signal %s", s)
		os.Exit(0)
	}
}

func get_version(cmd *cobra.Command, args []string) {
	common.Versions("OceanServer")
}

func execute() {
	if err := rootcmd.Execute(); nil != err {
		os.Exit(1)
	}
}

func main() {
	execute()
}
