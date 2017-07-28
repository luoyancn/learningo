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
	"errors"
	"fmt"
	"openstack"
	"os"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

// novaCmd represents the nova command
var novaCmd = &cobra.Command{
	Use:   "nova",
	Short: "The command of nova",
	Long: `
The commands of nova to request nova for vm instances.`,
	PreRunE: check_nova,
}

// bootCmd represents the boot command
var bootCmd = &cobra.Command{
	Use:   "boot",
	Short: "The sub command of nova",
	Long: `
The commands of nova to request boot vm instances.`,
	PreRunE: check_boot,
	Run:     exec_boot,
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "The sub command of nova",
	Long: `
The commands of nova to request list vm instances.`,
	Run: exec_list,
}

var (
	image   string
	flavor  string
	network string
	size    int
	region  string
)

func check_nova(cmd *cobra.Command, args []string) error {
	if "" == region {
		return errors.New("\"--region\" requires")
	}
	return nil
}

func init() {
	RootCmd.AddCommand(novaCmd)
	novaCmd.AddCommand(bootCmd)
	novaCmd.AddCommand(listCmd)
	novaCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "The region name")

	bootCmd.Flags().StringVarP(&image, "image", "i", "", "The image id")
	bootCmd.Flags().StringVarP(&flavor, "flavor", "f", "", "The flavor id")
	bootCmd.Flags().StringVarP(&network, "network", "n", "", "The network id")
	bootCmd.Flags().IntVarP(&size, "size", "s", 0,
		"The size of block device, GB")
}

func generate_common() (string, string) {
	token, context_ptr, err := openstack.GenerateContext(
		user, password, project, user_domain, project_domain, auth_url)

	if nil != err {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(2)
	}

	var req_url string
	for _, catalog := range (*context_ptr).Token.Catalog {
		if catalog.Type == "compute" {
			for _, endpoint := range catalog.Endpoints {
				if endpoint.Region == region &&
					endpoint.Interface == "public" {
					req_url = endpoint.Url
				}
			}
		}
	}

	req_url += "/servers"

	return token, req_url
}

func check_boot(cmd *cobra.Command, args []string) error {
	errs := []string{}
	if "" == image {
		errs = append(errs, "\"--image\" requires")
	}
	if "" == flavor {
		errs = append(errs, "\"--flavor\" requires")
	}
	if "" == network {
		errs = append(errs, "\"--network\" requires")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	if 0 < size {
		fmt.Printf(
			"Notice: Try to boot vm from cinder with block device size %d\n",
			size)
	}

	return nil
}

func exec_boot(cmd *cobra.Command, args []string) {
	token, req_url := generate_common()
	create_chan := make(chan string, concurrency)
	instance_uuids := make([]string, 0, concurrency)
	success := 0
	failed := 0
	begin := time.Now()
	for index := 0; index < concurrency; index++ {

		go func(token string, req_url string, flavor string,
			network string, size int) {
			servername := uuid.NewV4().String()
			instance_uuid, err := openstack.CreateInstance(
				token, req_url, servername, image, flavor, network, size)
			if nil != err {
				fmt.Printf("%v\n", err)
				create_chan <- ""
			} else {
				create_chan <- instance_uuid
			}
		}(token, req_url, flavor, network, size)

	}
	end := time.Since(begin)

	for index := 0; index < concurrency; index++ {
		instance_uuid := <-create_chan
		if "" == instance_uuid {
			failed += 1
		} else {
			instance_uuids = append(instance_uuids, instance_uuid)
			success += 1
		}
	}

	fmt.Println("############################################################")
	fmt.Printf("Sucess :%d, and faild %d and send %d requsts used %v\n",
		success, failed, concurrency, end)
	fmt.Printf("And the success instances uuid are %v\n", instance_uuids)

	if 0 == len(instance_uuids) {
		fmt.Printf("%s\n", "Faild to request create instances")
		return
	}

	fmt.Println("\nNow, let`s wait the instances becomes ACTIVE")
	// wait for instance to create success
	wait_chan := make(chan bool, len(instance_uuids))
	for _, instance_uuid := range instance_uuids {
		go func(token string, req_url string, instance_uuid string) {
			req_url += "/" + instance_uuid
			wait_begin := time.Now()
			var wait_retry time.Duration
			for {
				status, err := openstack.WaitInstance(token, req_url)
				if nil != err {
					fmt.Printf("Create Instance %s failed: %v\n",
						instance_uuid, err)
					wait_chan <- false
					break
				} else {
					wait_retry = time.Since(wait_begin)
					if "ACTIVE" == status {
						fmt.Printf("Instance %s created success and used %v\n",
							instance_uuid, wait_retry)
						wait_chan <- true
						break
					} else if "ERROR" == status {
						fmt.Printf(
							"Create Instance %s failed\n", instance_uuid)
						wait_chan <- false
						break
					}
					if 600 < wait_retry.Seconds() {
						fmt.Printf(
							"Create Instance %s failed after waiting %v \n",
							instance_uuid, wait_retry)
						wait_chan <- false
						break
					}
					fmt.Printf("The instance %s status is %s,"+
						" try to get it`s status after 5 seconds\n",
						instance_uuid, status)
					time.Sleep(5 * time.Second)
				}
			}
		}(token, req_url, instance_uuid)
	}

	wait_success := 0
	wait_failed := 0
	for index := 0; index < len(instance_uuids); index++ {
		if <-wait_chan {
			wait_success += 1
		} else {
			wait_failed += 1
		}
	}
	fmt.Println("############################################################")
	fmt.Printf("Sucess :%d, and faild %d\n", wait_success, wait_failed)
}

func exec_list(cmd *cobra.Command, args []string) {
	token, req_url := generate_common()
	get_chan := make(chan bool, concurrency)
	success := 0
	failed := 0
	begin := time.Now()
	for index := 0; index < concurrency; index++ {
		go func(token string, req_url string) {
			instances, err := openstack.GetInstance(token, req_url)
			if nil != err {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				get_chan <- false
			} else {
				fmt.Printf("%v\n", instances)
				get_chan <- true
			}
		}(token, req_url)
	}
	end := time.Since(begin)
	for index := 0; index < concurrency; index++ {
		if <-get_chan {
			success += 1
		} else {
			failed += 1
		}
	}

	fmt.Println("############################################################")
	fmt.Printf("Sucess :%d, and faild %d and send %d requsts used %v\n",
		success, failed, concurrency, end)
}
