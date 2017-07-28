package openstack

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func createRequestCred(username string, password string, userdomain string,
	project string, projectdomain string) (string, error) {
	cred := _credAuth{
		Auth: _auth{
			Identity: _identity{
				Methods: []string{"password"},
				Password: _password{
					User: _user{
						Name:     username,
						Password: password,
						Domain: _domain{
							Name: userdomain,
						},
					},
				},
			},
			Scope: _scope{
				Project: _project{
					Name: project,
					Domain: _domain{
						Name: projectdomain,
					},
				},
			},
		},
	}

	cred_json_str, err := json.Marshal(&cred)
	if nil != err {
		return "", err
	}
	return string(cred_json_str), nil
}

func GenerateToken(user string, password string, project string,
	userdomain string, projectdomain string, auth_url string) (string, error) {
	req_str, err := createRequestCred(user, password,
		userdomain, project, projectdomain)
	if nil != err {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, auth_url,
		bytes.NewBufferString(req_str))
	req.Header.Set(CONTENT_TYPE, APPLICATION_JSON)

	transport := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if nil != err {
		return "", err
	}

	if http.StatusCreated != resp.StatusCode {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New(string(body))
	}
	return resp.Header.Get("X-Subject-Token"), nil
}

func GenerateContext(user string, password string,
	project string, userdomain string, projectdomain string,
	auth_url string) (string, *RespToken, error) {
	req_str, err := createRequestCred(user, password,
		userdomain, project, projectdomain)
	if nil != err {
		return "", nil, err
	}

	req, err := http.NewRequest(http.MethodPost, auth_url,
		bytes.NewBufferString(req_str))
	req.Header.Set(CONTENT_TYPE, APPLICATION_JSON)

	transport := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if nil != err {
		return "", nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if http.StatusCreated != resp.StatusCode {
		return "", nil, errors.New(string(body))
	}

	var token RespToken
	err = json.Unmarshal(body, &token)
	if nil != err {
		return "", nil, err
	}

	return resp.Header.Get("X-Subject-Token"), &token, nil
}

func ValidateToken(token string, url string) (string, error) {
	return request(&token, nil, "GET", url)
}
