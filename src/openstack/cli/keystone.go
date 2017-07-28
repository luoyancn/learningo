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
	"fmt"
	"net/url"
	"openstack"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// keystoneCmd represents the keystone command
var keystoneCmd = &cobra.Command{
	Use:   "keystone",
	Short: "The commands of keystone",
	Long: `
The commands of keystone to request keystone for auth and identity.`,
	Run: exec_keystone,
}

var (
	validate bool
)

func init() {
	RootCmd.AddCommand(keystoneCmd)
	keystoneCmd.Flags().BoolVarP(&validate, "validate", "v", false,
		"Only validate the token")
}

func exec_keystone(cmd *cobra.Command, args []string) {
	if !validate {
		create_token()
	} else {
		valid_token()
	}
}

func create_token() {
	if 0 > concurrency {
		concurrency = 1
	}
	create_chan := make(chan bool, concurrency)
	begin := time.Now()
	faild := 0
	success := 0
	for i := 0; i < concurrency; i++ {
		go func(user string, password string, project string,
			user_domain string, project_domain string, auth_url string) {
			token, err := openstack.GenerateToken(user, password, project,
				user_domain, project_domain, auth_url)
			if nil != err {
				fmt.Printf("%v\n", err)
				create_chan <- false
			} else {
				fmt.Printf("%s\n", token)
				create_chan <- true
			}
		}(user, password, project, user_domain, project_domain, auth_url)
	}
	end := time.Since(begin)

	for i := 0; i < concurrency; i++ {
		if <-create_chan {
			success += 1
		} else {
			faild += 1
		}
	}
	fmt.Println("############################################################")
	fmt.Printf("Sucess :%d, and faild %d and send %d requsts used %v\n",
		success, faild, concurrency, end)
}

func valid_token() {
	token, context_ptr, err := openstack.GenerateContext(
		user, password, project, user_domain, project_domain, auth_url)
	if nil != err {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(2)
	}

	url_obj, err := url.ParseRequestURI(auth_url)
	if nil != err {
		fmt.Fprintf(os.Stderr, "Fail to parse the url %v\n", err)
		os.Exit(2)
	}

	req_url := "http://" + url_obj.Host + "/v3/projects/" +
		(*context_ptr).Token.Project.Id

	validate_chan := make(chan bool, concurrency)
	begin := time.Now()
	for i := 0; i < concurrency; i++ {
		go func(token string, req_url string) {
			project_info, err := openstack.ValidateToken(token, req_url)
			if nil == err {
				fmt.Printf("%s\n", project_info)
				validate_chan <- true
			} else {
				fmt.Printf("%v\n", err)
				validate_chan <- false
			}
		}(token, req_url)
	}
	elasped := time.Since(begin)

	sucess := 0
	fail := 0
	for i := 0; i < concurrency; i++ {
		if <-validate_chan {
			sucess += 1
		} else {
			fail += 1
		}
	}
	fmt.Println("############################################################")
	fmt.Printf("Sucess :%d, and faild %d and send %d requsts used %v\n",
		sucess, fail, concurrency, elasped)
}
