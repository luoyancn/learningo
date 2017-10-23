package resources

import (
	"net/http/httptest"
	"testing"
)

func Testvalite_req_body_create_user_sucess(t *testing.T) {
	user_json := `{"name": "fake_name", "age": 12, "sex": "men"}`
	resp := httptest.NewRecorder()
	valid_res := valite_req_body(resp, user_json, create_user_loader)
	if !valid_res {
		t.Errorf("error")
	}
}
