package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
)

func main() {
	url_ori := "postgres://user:pass@host.com:5432/path?k=v"
	url_str_ptr, err := url.Parse(url_ori)
	if nil != err {
		fmt.Printf("Cannot parse the url :%v\n", err)
		os.Exit(2)
	}

	fmt.Println(url_str_ptr.Scheme)
	fmt.Println(url_str_ptr.User)

	fmt.Println(url_str_ptr.User.Username())
	fmt.Println(url_str_ptr.User.Password())

	fmt.Println(url_str_ptr.User.String())

	fmt.Println(url_str_ptr.Host)
	host, port, _ := net.SplitHostPort(url_str_ptr.Host)
	fmt.Println(host)
	fmt.Println(port)
	fmt.Println(url_str_ptr.Hostname())

	fmt.Println(url_str_ptr.Query())
	fmt.Println(url_str_ptr.RawQuery)
	m, _ := url.ParseQuery(url_str_ptr.RawQuery)
	fmt.Println(m)

	fmt.Println(url_str_ptr.Path)
}
