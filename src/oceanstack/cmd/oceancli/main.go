package main

import (
	"oceanstack/common"
	"oceanstack/db"
	"oceanstack/logging"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var once sync.Once
var configfile string

var rootcmd = &cobra.Command{
	Short: "Manage tools for ocean stack",
	Long:  ` This is the binary for ocean stack`,
}

var dbcmd = &cobra.Command{
	Use:   "db_sync",
	Short: "DB operations for ocean stack",
	Long:  ` DB operations for ocean stack`,
	Run:   migrate,
	Args:  cobra.NoArgs,
}

var vercmd = &cobra.Command{
	Use:   "version",
	Short: "Get the ocean cli version",
	Long:  ` Get the version of ocean cli binary`,
	Run:   get_version,
	Args:  cobra.NoArgs,
}

func init() {
	once.Do(func() {
		dbcmd.PersistentFlags().StringVarP(
			&configfile, "config-file", "c", "",
			"The full path of config file (Required)")
		dbcmd.MarkPersistentFlagRequired("config-file")
		rootcmd.AddCommand(dbcmd)
		rootcmd.AddCommand(vercmd)
	})
}

func migrate(cmd *cobra.Command, args []string) {
	common.ReadConfig(configfile, "oceancli",
		logging.STD_ENABLED|logging.FILE_ENABLED)
	db.InitDbConnection()
	db.MigrateDB()
}

func get_version(cmd *cobra.Command, args []string) {
	common.Versions("OceanCli")
}

func Execute() {
	if err := rootcmd.Execute(); nil != err {
		os.Exit(1)
	}
}

func main() {
	Execute()
}
