package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func main() {
	/*
		ssh_config := &ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.Password("luoyan"),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		}
	*/

	ssh_config_with_key := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			PublicKeyFile("D:\\ssh-key\\luoyan_root"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//Timeout:         5 * time.Second,
	}

	// conn, err := ssh.Dial("tcp", "10.1.1.10:22", ssh_config)
	conn, err := ssh.Dial("tcp", "10.1.1.10:22", ssh_config_with_key)

	if nil != err {
		fmt.Printf("ERROR:%v\n", err)
		os.Exit(2)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if nil != err {
		fmt.Printf("ERROR:%v\n", err)
		os.Exit(2)
	}
	defer session.Close()

	ssh_out, err := session.Output("ls -l /opt")
	if nil != err {
		fmt.Printf("ERROR:%v\n", err)
		os.Exit(2)
	}

	new_session, err := conn.NewSession()
	if nil != err {
		fmt.Printf("ERROR:%v\n", err)
		os.Exit(2)
	}
	fmt.Printf("%v\n", string(ssh_out))
	defer new_session.Close()

	success := make(chan struct{})
	go func() {
		err = new_session.Run("dnf update -y")
		success <- struct{}{}
	}()

dnf_loop:
	for {
		select {
		case <-success:
			fmt.Println("excute end")
			break dnf_loop
		case <-time.Tick(1 * time.Second):
			fmt.Println("Waiting the executing end")
		}
	}
}
