package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSet(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8080/req", strings.NewReader(`{"method":"SET","key":"test_key","value":"test_val"}`))
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(req)

	handler.ServeHTTP(resp, request)

	expected := `{"method":"SET","key":"test_key","value":"test_val"}`
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Body.String(), expected)
	}
}

func TestSetLongKey(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8080/req", strings.NewReader(`{"method":"SET","key":"test_key_aaaaaaaaaaaaaaaaaaa","value":"test_val"}`))
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(req)

	handler.ServeHTTP(resp, request)

	expected := `{"method":"SET","key":"test_key_aaaaaaaaaaaaaaaaaaa","value":"test_val","error":"key too long or empty"}`
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Body.String(), expected)
	}
}

func TestGet(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8080/req", strings.NewReader(`{"method":"GET","key":"test_key"}`))
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(req)

	handler.ServeHTTP(resp, request)

	expected := `{"method":"GET","key":"test_key","value":"test_val"}`
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Body.String(), expected)
	}
}

func TestGetFail(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8080/req", strings.NewReader(`{"method":"GET","key":"no_key"}`))
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(req)

	handler.ServeHTTP(resp, request)

	expected := `{"method":"GET","key":"no_key","error":"not found"}`
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Body.String(), expected)
	}
}

func TestExists(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8080/req", strings.NewReader(`{"method":"EXISTS","key":"test_key"}`))
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(req)

	handler.ServeHTTP(resp, request)

	expected := `{"method":"EXISTS","key":"test_key","result":"exists"}`
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Body.String(), expected)
	}
}

func TestExistsFail(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8080/req", strings.NewReader(`{"method":"EXISTS","key":"tmp_key"}`))
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(req)

	handler.ServeHTTP(resp, request)

	expected := `{"method":"EXISTS","key":"tmp_key","result":"not exists"}`
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Body.String(), expected)
	}
}

func TestRemove(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8080/req", strings.NewReader(`{"method":"DELETE","key":"test_key"}`))
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(req)

	handler.ServeHTTP(resp, request)

	expected := `{"method":"DELETE","key":"test_key","result":"success"}`
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Body.String(), expected)
	}
}

func TestRemoveFail(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8080/req", strings.NewReader(`{"method":"DELETE","key":"secondary_key"}`))
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(req)

	handler.ServeHTTP(resp, request)

	expected := `{"method":"DELETE","key":"secondary_key","error":"not found"}`
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Body.String(), expected)
	}
}

func TestWrongMethod(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8080/req", strings.NewReader(`{"method":"SOME_STRING","key":"test_key"}`))
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(req)

	handler.ServeHTTP(resp, request)

	expected := `{"error":"method not allowed. allowed methods is: GET,SET,DELETE,EXISTS"}`
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Body.String(), expected)
	}
}
