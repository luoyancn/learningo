package utils

import (
	"io/ioutil"
	"k8sdeploy/logging"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

func parsePublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		logging.ERROR.Printf("Cannot found the id_rsa files:%v\n", err)
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		logging.ERROR.Printf("Cannot parse the private key file:%v\n", err)
		return nil
	}
	return ssh.PublicKeys(key)
}

func generate_ssh_auth_config() *ssh.ClientConfig {
	ssh_config_with_key := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			parsePublicKeyFile("/root/.ssh/id_rsa"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	return ssh_config_with_key
}

func SSHCheck(k8snodes ...string) bool {
	logging.INFO.Printf("Waiting for ssh check for nodes :%v\n", k8snodes)
	ssh_key := generate_ssh_auth_config()

	check := make(chan bool, len(k8snodes))
	runtime.GOMAXPROCS(len(k8snodes))
	for _, node := range k8snodes {
		go func(node string) {
			logging.INFO.Printf("Connecting to the host %s ...\n", node)
			conn, err := ssh.Dial(
				"tcp", strings.Join([]string{node, "22"}, ":"), ssh_key)
			if nil != err {
				logging.ERROR.Printf(
					"Cannot connect to the host %s: %v\n", node, err)
				check <- false
			}
			defer conn.Close()
			logging.INFO.Printf("Success connect to the host %s \n", node)
			check <- true
		}(node)
	}

	for index := 0; index < len(k8snodes); index++ {
		if !<-check {
			return false
		}
	}
	return true
}

func get_binary_files(path string) []string {
	files := []string{}
	err := filepath.Walk(
		path, func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			files = append(files, path)
			return nil
		})
	if err != nil {
		logging.ERROR.Printf("%v\n", err)
	}
	return files
}

func SCPBinary(binary_type string, k8snodes ...string) bool {
	binary_files_path := viper.GetString(
		strings.Join([]string{binary_type, "binary_path"}, "."))
	files := get_binary_files(binary_files_path)
	if 0 == len(files) {
		logging.ERROR.Printf(
			"Please ensure the %s binary files in your path %s\n",
			binary_type, binary_files_path)
		return false
	}
	logging.DEBUG.Printf("The %s binary files were follows: %v\n",
		binary_type, files)
	return true
}
