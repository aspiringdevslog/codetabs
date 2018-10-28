package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type headersTestOutput struct {
	StatusCode int    // `json:"statusCode"`
	Error      string `json:"Error,omitempty"`
}

var headersTests = []struct {
	it         string
	endpoint   string
	errorText  string
	statusCode int
}{
	{"not domain (empty)", "/v1/headers/", "ERROR Bad Request, not enough parameters", 400},
	{"not valid hostname", "/v1/headers/no-valid-hostname.com", "ERROR Couldn't resolve host no-valid-hostname.com", 400},
	{"not registered domain", "/v1/headers/sure-this-is-not-registered.com", "ERROR Couldn't resolve host sure-this-is-not-registered.com", 400},
	{"domain without redirection", "/v1/headers/codetabs.com", "", 200},
	{"domain with redirection", "/v1/headers/www.codetabs.com", "", 200},
}

func TestHeadersApi(t *testing.T) {
	c.App.Mode = "test"

	for _, test := range headersTests {
		var to = headersTestOutput{}
		pass := true
		//fmt.Println(`Test...`, test.endpoint, "...", test.it)
		req, err := http.NewRequest("GET", test.endpoint, nil)
		if err != nil {
			//fmt.Println(`------------------------------`)
			//fmt.Println(err.Error())
			//fmt.Println(test.errorText)
			//fmt.Println(`------------------------------`)
			if err.Error() != test.errorText {
				t.Fatalf("Error Request %s\n", err.Error())
			} else {
				pass = false
			}
		}
		if pass {
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(router)
			handler.ServeHTTP(rr, req)
			if rr.Code != test.statusCode {
				t.Errorf("%s got %v want %v\n", test.endpoint, rr.Code, test.statusCode)
			}
			_ = json.Unmarshal(rr.Body.Bytes(), &to)
			if to.Error != test.errorText {
				t.Errorf("%s got %v want %v\n", test.endpoint, to.Error, test.errorText)
			}
		}
		fmt.Printf("Test OK...%s\n", test.endpoint)
	}
}
