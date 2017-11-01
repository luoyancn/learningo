// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "openstack",
	Short: "A brief description of your application",
	Long: `
The commands to request openstack for resources.`,
	// Uncomment the following line if your bare  application
	// has an action associated with it:
	PersistentPreRunE: check_required,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var (
	user           string
	password       string
	project        string
	auth_url       string
	concurrency    int
	user_domain    string
	project_domain string
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.openstack.yaml)")

	RootCmd.PersistentFlags().StringVar(&user, "user", "",
		"The name of openstack user")
	RootCmd.PersistentFlags().StringVar(&password, "password", "",
		"The password of openstack user")
	RootCmd.PersistentFlags().StringVar(&project, "project", "",
		"The project of openstack user")
	RootCmd.PersistentFlags().StringVar(&auth_url, "auth-url", "",
		"The auth url of openstack keystone")
	RootCmd.PersistentFlags().IntVar(&concurrency, "concurrency", 10,
		"The concurrency number of requests")
	RootCmd.PersistentFlags().StringVar(&user_domain, "user-domain",
		"default", "The domain name of user")
	RootCmd.PersistentFlags().StringVar(&project_domain, "project-domain",
		"default", "The domain name of project")
	RootCmd.MarkPersistentFlagRequired("user")
	RootCmd.MarkPersistentFlagRequired("password")
	RootCmd.MarkPersistentFlagRequired("project")
	RootCmd.MarkPersistentFlagRequired("auth-url")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func check_required(cmd *cobra.Command, args []string) error {
	/*
		errs := []string{}
		if "" == user {
			errs = append(errs, "\"--user\" requires")
		}
		if "" == password {
			errs = append(errs, "\"--password\" requires")
		}
		if "" == project {
			errs = append(errs, "\"--project\" requires")
		}
		if "" == auth_url {
			errs = append(errs, "\"--auth-url\" requires")
		}
		if len(errs) > 0 {
			return errors.New(strings.Join(errs, "\n"))
		}
		/*/

	_, err := url.ParseRequestURI(auth_url)
	if nil != err {
		return err
	}
	return nil
}
