package utils

import (
	"encoding/json"
	"io/ioutil"
	"k8sdeploy/logging"
	"os"
	"path"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/cli/sign"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/signer"
	"github.com/spf13/viper"
)

type outputFile struct {
	Filename string
	Contents string
	IsBinary bool
	Perms    os.FileMode
}

func writeFile(filespec, contents string, perms os.FileMode) error {
	err := ioutil.WriteFile(filespec, []byte(contents), perms)
	if nil != err {
		logging.LOG.Errorf("Cannot Write to the CA files:%v\n", err)
		return err
	}
	return nil
}

func generate_ca_files(ca_name string, template_path string,
	config *cli.Config) error {
	csrJSONFileBytes, err := cli.ReadStdin(
		path.Join(template_path, ca_name+"-csr.json"))
	if err != nil {
		return err
	}

	req := csr.CertificateRequest{
		KeyRequest: csr.NewBasicKeyRequest(),
	}

	err = json.Unmarshal(csrJSONFileBytes, &req)
	if err != nil {
		return err
	}

	var key, csrPEM, cert []byte
	if nil == config {
		cert, csrPEM, key, err = initca.New(&req)
		if nil != err {
			return err
		}
	} else {
		g := &csr.Generator{Validator: genkey.Validator}
		csrPEM, key, err = g.ProcessRequest(&req)
		if err != nil {
			key = nil
			return err
		}

		s, err := sign.SignerFromConfig(*config)
		if err != nil {
			return err
		}
		signReq := signer.SignRequest{
			Request: string(csrPEM),
			Hosts:   signer.SplitHosts((*config).Hostname),
			Profile: (*config).Profile,
			Label:   (*config).Label,
		}

		cert, err = s.Sign(signReq)
		if err != nil {
			return err
		}
	}

	var outs []outputFile
	if cert != nil {
		outs = append(outs, outputFile{
			Filename: path.Join(template_path, ca_name+".pem"),
			Contents: string(cert),
			Perms:    0664,
		})
	}

	if key != nil {
		outs = append(outs, outputFile{
			Filename: path.Join(template_path, ca_name+"-key.pem"),
			Contents: string(key),
			Perms:    0600,
		})
	}

	if csrPEM != nil {
		outs = append(outs, outputFile{
			Filename: path.Join(template_path, ca_name+".csr"),
			Contents: string(csrPEM),
			Perms:    0600,
		})
	}

	for _, e := range outs {
		if err := writeFile(e.Filename, e.Contents, e.Perms); nil != err {
			return err
		}
	}
	return nil
}

func CreateCert() error {
	template_path := viper.GetString("cfs.templates")
	err := generate_ca_files("ca", template_path, nil)
	if nil != err {
		logging.LOG.Errorf("Fail to generate the CA files:%v\n", err)
		return err
	}
	var config cli.Config
	config.Address = "127.0.0.1"
	config.Hostname = ""
	config.Label = ""
	config.Port = 8888
	config.CAFile = path.Join(template_path, "ca.pem")
	config.CAKeyFile = path.Join(template_path, "ca-key.pem")
	config.ConfigFile = path.Join(template_path, "ca-config.json")
	config.Profile = "kubernetes"

	err = generate_ca_files("admin", template_path, &config)
	if nil != err {
		logging.LOG.Errorf("Fail to generate the admin ca files:%v\n", err)
		return err
	}

	err = generate_ca_files("kube-proxy", template_path, &config)
	if nil != err {
		logging.LOG.Errorf(
			"Fail to generate the kube-proxy ca files:%v\n", err)
		return err
	}

	err = generate_ca_files("kubernetes", template_path, &config)
	if nil != err {
		logging.LOG.Errorf(
			"Fail to generate the kubernetes ca files:%v\n", err)
		return err
	}

	return nil
}
