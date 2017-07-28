package openstack

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	CONTENT_TYPE     = "Content-Type"
	APPLICATION_JSON = "application/json"
	X_AUTH_TOKEN     = "X-Auth-Token"
)

func request(token *string, reqdata *string,
	method string, url string) (string, error) {

	var reqbuf io.Reader

	if nil != reqdata && "" != *reqdata {
		reqbuf = bytes.NewBufferString(*reqdata)
	} else {
		reqbuf = nil
	}

	req, err := http.NewRequest(strings.ToUpper(method), url, reqbuf)
	req.Header.Set(CONTENT_TYPE, APPLICATION_JSON)

	if nil != token && "" != *token {
		req.Header.Set(X_AUTH_TOKEN, *token)
	}

	transport := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if nil != err {
		return "", err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if 1 != resp.StatusCode/200 {
		return "", errors.New(string(body))
	}

	return string(body), nil
}
