package httpgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type reqbody struct {
	Name    string `json:"name"`
	Id      int    `json:"id"`
	Enabled bool   `json:"enabled,omitempty"`
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
	var req_body reqbody
	err = json.Unmarshal([]byte(body), &req_body)
	if err != nil {
		http.Error(writer, "Only json-liked body accepted\n",
			http.StatusBadRequest)
		return
	}
	if "" == req_body.Name || 0 == req_body.Id {
		http.Error(writer, "Missing id or name in request body\n",
			http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(writer, "Recevied request body %v\n", string(body))
}

func go_put_handler(writer http.ResponseWriter, request *http.Request) {
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
