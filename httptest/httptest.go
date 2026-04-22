// Package httptest provides HTTP testing utilities.
package httptest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Response is a test response wrapper
type Response struct {
	*httptest.ResponseRecorder
}

// NewResponse creates a new Response recorder
func NewResponse() *Response {
	return &Response{httptest.NewRecorder()}
}

// Request creates a request from method, path, and body
func Request(method, path string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// Get creates a GET request
func Get(path string) *http.Request {
	return Request(http.MethodGet, path, nil)
}

// Post creates a POST request with JSON body
func Post(path string, data interface{}) *http.Request {
	var body io.Reader
	if data != nil {
		b, _ := json.Marshal(data)
		body = strings.NewReader(string(b))
	}
	return Request(http.MethodPost, path, body)
}

// Put creates a PUT request with JSON body
func Put(path string, data interface{}) *http.Request {
	var body io.Reader
	if data != nil {
		b, _ := json.Marshal(data)
		body = strings.NewReader(string(b))
	}
	return Request(http.MethodPut, path, body)
}

// Delete creates a DELETE request
func Delete(path string) *http.Request {
	return Request(http.MethodDelete, path, nil)
}

// JSON returns the response body as parsed JSON
func (r *Response) JSON(v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// String returns the response body as string
func (r *Response) String() string {
	return r.Body.String()
}

// OK asserts the response status is 200
func (r *Response) OK(t TestingTB) {
	if r.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", r.Code)
	}
}

// Created asserts the response status is 201
func (r *Response) Created(t TestingTB) {
	if r.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", r.Code)
	}
}

// BadRequest asserts the response status is 400
func (r *Response) BadRequest(t TestingTB) {
	if r.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", r.Code)
	}
}

// NotFound asserts the response status is 404
func (r *Response) NotFound(t TestingTB) {
	if r.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", r.Code)
	}
}

// StatusCode asserts a specific status code
func (r *Response) StatusCode(t TestingTB, code int) {
	if r.Code != code {
		t.Fatalf("expected status %d, got %d", code, r.Code)
	}
}

// TestingTB is the testing interface
type TestingTB interface {
	Fatalf(format string, args ...interface{})
}
