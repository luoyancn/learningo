package utils

import (
	"encoding/json"
	"io/ioutil"
	"k8sdeploy/logging"
	"os"
	"path"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/spf13/viper"
)

type outputFile struct {
	Filename string
	Contents string
	IsBinary bool
	Perms    os.FileMode
}

func writeFile(filespec, contents string, perms os.FileMode) err {
	err := ioutil.WriteFile(filespec, []byte(contents), perms)
	if nil != err {
		logging.ERROR.Print("Cannot Write to the CA files:%v\n", err)
		return err
	}
	return nil
}

func InitCA() error {
	template_path := viper.GetString("cfs.templates")
	csrJSONFileBytes, err := cli.ReadStdin(
		path.Join(template_path, "admin-csr.json"))
	if err != nil {
		logging.ERROR.Print("Cannot to generate the CA files:%v\n", err)
		return err
	}

	req := csr.CertificateRequest{
		KeyRequest: csr.NewBasicKeyRequest(),
	}
	err = json.Unmarshal(csrJSONFileBytes, &req)
	if err != nil {
		logging.ERROR.Print(
			"Invalid file content:%v, need json-like files\n", err)
		return err
	}
	var key, csrPEM, cert []byte
	cert, csrPEM, key, err = initca.New(&req)
	if nil != err {
		logging.ERROR.Print("Failed to generate the CA files:%v\n", err)
		return err
	}

	base_name := "ca"
	var outs []outputFile
	if cert != nil {
		outs = append(outs, outputFile{
			Filename: path.Join(template_path, base_name+".pem"),
			Contents: string(cert),
			Perms:    0664,
		})
	}

	if key != nil {
		outs = append(outs, outputFile{
			Filename: path.Join(template_path, base_name+"-key.pem"),
			Contents: string(key),
			Perms:    0600,
		})
	}

	if csrPEM != nil {
		outs = append(outs, outputFile{
			Filename: path.Join(template_path, base_name+".csr"),
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
