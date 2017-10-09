package httpgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

type reqbody struct {
	Name    string `json:"name"`
	Id      string `json:"id,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

const post_body_schema = `
{
	"title": "POSTBODY",
	"type": "object",
	"properties": {
		"name": {
			"type": "string",
			"minLength": 6,
			"maxLength": 36
		},
		"id": {
			"type": "string",
			"minLength": 32,
			"maxLength": 64,
			"format": "uuid"
		},
		"enabled": {
			"type": "boolean"
		}
	},
	"required" : ["name"],
	"additionalProperties": false
}
`

const put_body_schema = `
{
	"title": "POSTBODY",
	"type": "object",
	"properties": {
		"name": {
			"type": "string",
			"minLength": 6,
			"maxLength": 36
		},
		"enabled": {
			"type": "boolean"
		}
	},
	"additionalProperties": false,
	"anyOf": [
		{
			"required": ["name"]
		},
		{
			"required": ["enabled"]
		}
	]
}
`

var once_post_loader sync.Once
var post_validate_loader gojsonschema.JSONLoader
var put_validate_loader gojsonschema.JSONLoader
var loader_map map[string]gojsonschema.JSONLoader

func get_json_validate_loader() map[string]gojsonschema.JSONLoader {
	once_post_loader.Do(func() {
		loader_map = make(map[string]gojsonschema.JSONLoader, 2)
		post_validate_loader = gojsonschema.NewStringLoader(post_body_schema)
		put_validate_loader = gojsonschema.NewStringLoader(put_body_schema)
		loader_map["POST"] = post_validate_loader
		loader_map["PUT"] = put_validate_loader
	})
	return loader_map
}

func valite_req_body(writer http.ResponseWriter, str_body string,
	loader gojsonschema.JSONLoader) bool {
	req_body := gojsonschema.NewStringLoader(str_body)
	res, err := gojsonschema.Validate(loader, req_body)
	if err != nil {
		http.Error(writer, "Only json-liked body accepted\n",
			http.StatusBadRequest)
		return false
	}

	if !res.Valid() {
		var errs []string
		errs = append(errs, fmt.Sprintln(
			"Input is invalid, following errors found:"))
		for _, desc := range res.Errors() {
			errs = append(errs, fmt.Sprintf("- %s", desc))
		}
		http.Error(writer, strings.Join(errs, "\n"),
			http.StatusBadRequest)
		return false
	}

	return true
}

func HttpGoHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet, http.MethodHead:
		go_get_handler(writer, request)
	case http.MethodPost:
		go_post_handler(writer, request)
	case http.MethodPut:
		go_put_handler(writer, request)
	case http.MethodDelete:
		go_delete_handler(writer, request)
	default:
		go_unsupported_handler(writer, request)
	}
}

func go_get_handler(writer http.ResponseWriter, request *http.Request) {
	poem_name := request.Form.Get("name")
	if "" == poem_name {
		fmt.Fprint(writer, "Visit the poem list\n")
		return
	}
	fmt.Fprintf(writer, "Poem %s coming soon\n", poem_name)
}

func go_post_handler(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Failed to parse the body in request\n",
			http.StatusBadRequest)
		return
	}
	if !valite_req_body(
		writer, string(body), get_json_validate_loader()["POST"]) {
		return
	}
	var req_entity reqbody
	err = json.Unmarshal(body, &req_entity)
	if nil != err {
		http.Error(writer, "Failed to parse the body in request\n",
			http.StatusBadRequest)
		return
	}
	fmt.Printf("%v\n", req_entity)
	writer.WriteHeader(http.StatusCreated)
	fmt.Fprintf(writer, "Recevied request body %v\n", string(body))
}

func go_put_handler(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Failed to parse the body in request\n",
			http.StatusBadRequest)
		return
	}
	if !valite_req_body(
		writer, string(body), get_json_validate_loader()["PUT"]) {
		return
	}
	var req_entity reqbody
	err = json.Unmarshal(body, &req_entity)
	if nil != err {
		http.Error(writer, "Failed to parse the body in request\n",
			http.StatusBadRequest)
		return
	}
	fmt.Printf("%v\n", req_entity)
	writer.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(writer, "Recevied request body %v\n", string(body))
}

func go_delete_handler(writer http.ResponseWriter, request *http.Request) {
}

func go_unsupported_handler(writer http.ResponseWriter,
	request *http.Request) {
	http.Error(writer, "Not supported http request method\n",
		http.StatusMethodNotAllowed)
}

func RunServer() {
	http.HandleFunc("/poem", HttpGoHandler)
	http.ListenAndServe(":8080", nil)
}
