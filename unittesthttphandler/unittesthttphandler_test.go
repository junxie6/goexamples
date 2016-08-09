package hello

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// Reference:
// http://www.markjberger.com/testing-web-apps-in-golang/
// https://elithrar.github.io/article/testing-http-handlers-go/
func TestSrvHello(t *testing.T) {
	// create a new http request
	var jsonStr = []byte(`{"Data":{"test":{"name":"bot"}}}`)

	req, err := http.NewRequest("GET", "/hello?act=say", bytes.NewBuffer(jsonStr))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	// create a new ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(SrvHello)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"Status":true,"ErrArr":null,"ErrCount":0,"ObjArr":null,"Data":{"name":"bot"}}`
	out := strings.TrimRight(rr.Body.String(), "\n")

	if ok, err := JSONDeepEqual(expected, out); err != nil {
		t.Error(err.Error())
	} else if !ok {
		t.Errorf("handler returned unexpected body: got %v want %v", out, expected)
	}
}

// JSONDeepEqual ...
func JSONDeepEqual(s1 string, s2 string) (bool, error) {
	var m1, m2 map[string]interface{}

	if err := json.Unmarshal([]byte(s1), &m1); err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(s2), &m2); err != nil {
		return false, err
	}

	return reflect.DeepEqual(m1, m2), nil
}
