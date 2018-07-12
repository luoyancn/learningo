package main

import (
	"context"
	"fmt"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
)

func main() {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()
	config := etcd.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
		Context:     ctx,
	}
	client, err := etcd.New(config)
	if nil != err {
		fmt.Printf("error occured:%v\n", err)
	}
	resp, err := client.Grant(ctx, int64(5*time.Second))
	if nil != err {
		switch t := err.(type) {
		default:
			fmt.Printf("%v\n", t)
		}
		return
	}

	fmt.Printf("%v\n", resp.String())

	time.Sleep(10 * time.Second)
	res, err := client.Get(ctx, "zhangjl")
	if nil != err {
		fmt.Printf("ERROR:%s\n", err)
	}

	fmt.Printf("%v\n", res)
}
