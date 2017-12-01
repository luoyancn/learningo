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
	"golang.org/x/crypto/ssh"
)

func parsePublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		logging.LOG.Errorf("Cannot found the id_rsa files:%v\n", err)
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		logging.LOG.Errorf("Cannot parse the private key file:%v\n", err)
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

func get_files(path string) []string {
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
		logging.LOG.Errorf("%v\n", err)
	}
	return files
}

func scp_source_to_dest(ssh_conn *ssh.Client, file_type string,
	dest_path string, binarys ...string) bool {
	sftp_client, err := sftp.NewClient(ssh_conn)
	if nil != err {
		logging.LOG.Errorf("Cannot create scp tunnel:%v\n", err)
		return false
	}
	defer sftp_client.Close()
	_, err = sftp_client.Stat(dest_path)
	if nil != err {
		logging.LOG.Warningf(
			"The target path %s not exsit:%v, Try to create it ...\n",
			dest_path, err)
		ssh_session, err := ssh_conn.NewSession()
		if nil != err {
			logging.LOG.Errorf("Cannot create ssh tunnel:%v\n", err)
			return false
		}
		cmd := "mkdir -p " + dest_path
		if err = ssh_session.Run(cmd); nil != err {
			logging.LOG.Errorf("Cannot create the dest path %s :%v\n",
				dest_path, err)
			return false
		}
	}

	//	scp_res := make(chan bool, len(binarys))
	//	runtime.GOMAXPROCS(len(binarys))
	for _, binary := range binarys {
		//		go func(binary string) {
		logging.LOG.Infof("Scp file %s to remote ...\n", binary)
		src_binary, err := os.Open(binary)
		if nil != err {
			logging.LOG.Errorf(
				"Cannot read binary file %s:%v\n", binary, err)
			//				scp_res <- false
			//				return
			return false
		}
		defer src_binary.Close()

		_, dest_binary_name := filepath.Split(binary)
		dest_binary_full_path := path.Join(dest_path, dest_binary_name)
		_, err = sftp_client.Stat(dest_binary_full_path)
		if nil == err {
			logging.LOG.Warningf("The target file exsit\n")
		} else {
			dest_binary, err := sftp_client.Create(dest_binary_full_path)
			if err != nil {
				logging.LOG.Errorf(
					"Fail to create remote file %s in dest path %s :%v\n",
					dest_binary_name, dest_path, err)
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
		if file_type == "binary" {
			sftp_client.Chmod(dest_binary_full_path, 0755)
		}
		logging.LOG.Infof(
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
	logging.LOG.Infof("Waiting for ssh check for nodes :%v\n", k8snodes)
	ssh_key := generate_ssh_auth_config()

	check := make(chan bool, len(k8snodes))
	runtime.GOMAXPROCS(len(k8snodes))
	for _, node := range k8snodes {
		go func(node string) {
			logging.LOG.Infof("Connecting to the host %s ...\n", node)
			conn, err := get_ssh_connection(node, ssh_key)
			if nil != err {
				logging.LOG.Errorf(
					"Cannot connect to the host %s: %v\n", node, err)
				check <- false
				return
			}
			check <- true
			defer conn.Close()
			logging.LOG.Infof("Success connect to the host %s\n", node)
		}(node)
	}

	for index := 0; index < len(k8snodes); index++ {
		if !<-check {
			return false
		}
	}
	return true
}

func SCPFiles(source_path []string, dest_path string,
	file_type string, k8snodes ...string) bool {
	files := []string{}
	for _, pat := range source_path {
		files = append(files, get_files(pat)...)
	}
	if 0 == len(files) {
		logging.LOG.Errorf(
			"There is no files in your source path %s\n", source_path)
		return false
	}
	logging.LOG.Debugf("The files were follows: %v\n", files)

	ssh_res := make(chan bool, len(k8snodes))
	runtime.GOMAXPROCS(len(k8snodes))

	for _, node := range k8snodes {
		go func(node string) {
			logging.LOG.Infof("Scp files to host %s ...\n", node)
			ssh_conn, err := get_ssh_connection(
				node, generate_ssh_auth_config())
			if nil != err {
				logging.LOG.Errorf(
					"Cannot scp the files to %s: %v\n", node, err)
				ssh_res <- false
				return
			}
			defer ssh_conn.Close()
			if !scp_source_to_dest(ssh_conn, file_type, dest_path, files...) {
				logging.LOG.Errorf(
					"Fail to scp the files to %s\n", node)
				ssh_res <- false
				return
			}
			logging.LOG.Infof(
				"Scp files to host %s success \n", node)
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
