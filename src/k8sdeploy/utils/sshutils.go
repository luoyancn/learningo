package utils

import (
	"io/ioutil"
	"k8sdeploy/logging"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/sftp"
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

func get_ssh_connection(host string,
	ssh_key *ssh.ClientConfig) (*ssh.Client, error) {
	conn, err := ssh.Dial(
		"tcp", strings.Join([]string{host, "22"}, ":"), ssh_key)
	if nil != err {
		return nil, err
	}
	return conn, nil
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

func scp_source_to_dest(ssh_conn *ssh.Client, binarys ...string) bool {
	sftp_client, err := sftp.NewClient(ssh_conn)
	if nil != err {
		logging.ERROR.Printf("Cannot create scp tunnel:%v\n", err)
		return false
	}
	defer sftp_client.Close()
	dest_binary_path := viper.GetString("default.target_path")
	_, err = sftp_client.Stat(dest_binary_path)
	if nil != err {
		logging.WARNING.Printf(
			"The target binary path %s not exsit:%v, Try to create it ...\n",
			dest_binary_path, err)
		if err = sftp_client.Mkdir(dest_binary_path); nil != err {
			logging.TRACE.Printf(
				"Cannot create the path :%v\n", err)
			return false
		}
	}

	//	scp_res := make(chan bool, len(binarys))
	//	runtime.GOMAXPROCS(len(binarys))
	for _, binary := range binarys {
		//		go func(binary string) {
		logging.INFO.Printf("Scp file %s to remote ...\n", binary)
		src_binary, err := os.Open(binary)
		if nil != err {
			logging.ERROR.Printf(
				"Cannot read binary file %s:%v\n", binary, err)
			//				scp_res <- false
			//				return
			return false
		}
		defer src_binary.Close()

		_, dest_binary_name := filepath.Split(binary)
		dest_binary_full_path := path.Join(dest_binary_path, dest_binary_name)
		_, err = sftp_client.Stat(dest_binary_full_path)
		if nil == err {
			logging.WARNING.Printf("The target binary file exsit\n")
		} else {
			dest_binary, err := sftp_client.Create(dest_binary_full_path)
			if err != nil {
				logging.ERROR.Printf(
					"Fail to create remote file %s in dest path %s :%v\n",
					dest_binary_name, dest_binary_path, err)
				//				scp_res <- false
				//				return
				return false
			}
			defer dest_binary.Close()

			buf := make([]byte, 1024)
			for {
				n, _ := src_binary.Read(buf)
				if n == 0 {
					break
				}
				dest_binary.Write(buf[0:n])
			}
			//			scp_res <- true
		}
		sftp_client.Chmod(dest_binary_full_path, 0755)
		logging.INFO.Printf(
			"Success scp file %s to remote\n", dest_binary_name)
		//		}(binary)
	}

	/*
		for index := 0; index < len(binarys); index++ {
			if !<-scp_res {
				return false
			}
		}
	*/
	return true
}

func SSHCheck(k8snodes ...string) bool {
	logging.INFO.Printf("Waiting for ssh check for nodes :%v\n", k8snodes)
	ssh_key := generate_ssh_auth_config()

	check := make(chan bool, len(k8snodes))
	runtime.GOMAXPROCS(len(k8snodes))
	for _, node := range k8snodes {
		go func(node string) {
			logging.INFO.Printf("Connecting to the host %s ...\n", node)
			conn, err := get_ssh_connection(node, ssh_key)
			if nil != err {
				logging.ERROR.Printf(
					"Cannot connect to the host %s: %v\n", node, err)
				check <- false
				return
			}
			check <- true
			defer conn.Close()
			logging.INFO.Printf("Success connect to the host %s\n", node)
		}(node)
	}

	for index := 0; index < len(k8snodes); index++ {
		if !<-check {
			return false
		}
	}
	return true
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

	ssh_res := make(chan bool, len(k8snodes))
	runtime.GOMAXPROCS(len(k8snodes))

	for _, node := range k8snodes {
		go func(node string) {
			logging.INFO.Printf("Scp %s binarys to host %s ...\n",
				binary_type, node)
			ssh_conn, err := get_ssh_connection(
				node, generate_ssh_auth_config())
			if nil != err {
				logging.ERROR.Printf(
					"Cannot scp the %s binarys to %s: %v\n",
					binary_type, node, err)
				ssh_res <- false
				return
			}
			defer ssh_conn.Close()
			if !scp_source_to_dest(ssh_conn, files...) {
				logging.ERROR.Printf(
					"Fail to scp the %s binarys to %s\n", binary_type, node)
				ssh_res <- false
				return
			}
			logging.INFO.Printf(
				"Scp %s binarys to host %s success \n", binary_type, node)
			ssh_res <- true
		}(node)
	}

	for index := 0; index < len(k8snodes); index++ {
		if !<-ssh_res {
			return false
		}
	}
	return true
}
