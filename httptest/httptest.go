// Package httptest provides HTTP testing utilities.
package httptest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

// Patch creates a PATCH request with JSON body
func Patch(path string, data interface{}) *http.Request {
	var body io.Reader
	if data != nil {
		b, _ := json.Marshal(data)
		body = strings.NewReader(string(b))
	}
	return Request(http.MethodPatch, path, body)
}

// Head creates a HEAD request
func Head(path string) *http.Request {
	return Request(http.MethodHead, path, nil)
}

// Options creates an OPTIONS request
func Options(path string) *http.Request {
	return Request(http.MethodOptions, path, nil)
}

// Delete creates a DELETE request
func Delete(path string) *http.Request {
	return Request(http.MethodDelete, path, nil)
}

// WithHeader sets a header on the request
func WithHeader(req *http.Request, key, value string) *http.Request {
	req.Header.Set(key, value)
	return req
}

// WithHeaders sets multiple headers on the request
func WithHeaders(req *http.Request, headers map[string]string) *http.Request {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return req
}

// WithAuth sets Authorization header
func WithAuth(req *http.Request, username, password string) *http.Request {
	req.SetBasicAuth(username, password)
	return req
}

// WithBearer sets Authorization header with Bearer token
func WithBearer(req *http.Request, token string) *http.Request {
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// WithQuery sets query parameters on the request
func WithQuery(req *http.Request, params map[string]string) *http.Request {
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	return req
}

// WithForm creates a request with form data
func WithForm(req *http.Request, data map[string]string) *http.Request {
	form := url.Values{}
	for key, value := range data {
		form.Add(key, value)
	}
	req.Body = io.NopCloser(strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

// JSON returns the response body as parsed JSON
func (r *Response) JSON(v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// String returns the response body as string
func (r *Response) String() string {
	return r.Body.String()
}

// Bytes returns the response body as bytes
func (r *Response) Bytes() []byte {
	return r.Body.Bytes()
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

// Accepted asserts the response status is 202
func (r *Response) Accepted(t TestingTB) {
	if r.Code != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d", r.Code)
	}
}

// NoContent asserts the response status is 204
func (r *Response) NoContent(t TestingTB) {
	if r.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", r.Code)
	}
}

// BadRequest asserts the response status is 400
func (r *Response) BadRequest(t TestingTB) {
	if r.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", r.Code)
	}
}

// Unauthorized asserts the response status is 401
func (r *Response) Unauthorized(t TestingTB) {
	if r.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", r.Code)
	}
}

// Forbidden asserts the response status is 403
func (r *Response) Forbidden(t TestingTB) {
	if r.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", r.Code)
	}
}

// NotFound asserts the response status is 404
func (r *Response) NotFound(t TestingTB) {
	if r.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", r.Code)
	}
}

// Conflict asserts the response status is 409
func (r *Response) Conflict(t TestingTB) {
	if r.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", r.Code)
	}
}

// InternalServerError asserts the response status is 500
func (r *Response) InternalServerError(t TestingTB) {
	if r.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", r.Code)
	}
}

// StatusCode asserts a specific status code
func (r *Response) StatusCode(t TestingTB, code int) {
	if r.Code != code {
		t.Fatalf("expected status %d, got %d", code, r.Code)
	}
}

// GetHeader returns the value of a response header
func (r *Response) GetHeader(key string) string {
	return r.ResponseRecorder.Header().Get(key)
}

// ContentType asserts the Content-Type header
func (r *Response) ContentType(t TestingTB, expected string) {
	contentType := r.ResponseRecorder.Header().Get("Content-Type")
	if !strings.Contains(contentType, expected) {
		t.Fatalf("expected Content-Type containing %q, got %q", expected, contentType)
	}
}

// JSONEq asserts the response body equals the expected JSON
func (r *Response) JSONEq(t TestingTB, expected interface{}) {
	var actual interface{}
	if err := json.Unmarshal(r.Body.Bytes(), &actual); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	expectedJSON, _ := json.Marshal(expected)
	var expectedInterface interface{}
	if err := json.Unmarshal(expectedJSON, &expectedInterface); err != nil {
		t.Fatalf("expected value is not valid JSON: %v", err)
	}
	if !jsonEq(actual, expectedInterface) {
		t.Fatalf("expected JSON:\n%s\ngot:\n%s", expectedJSON, r.Body.String())
	}
}

// JSONContains asserts the response JSON contains the expected value
func (r *Response) JSONContains(t TestingTB, expected interface{}) {
	var actual interface{}
	if err := json.Unmarshal(r.Body.Bytes(), &actual); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	expectedJSON, _ := json.Marshal(expected)
	var expectedInterface interface{}
	if err := json.Unmarshal(expectedJSON, &expectedInterface); err != nil {
		t.Fatalf("expected value is not valid JSON: %v", err)
	}
	if !jsonContains(actual, expectedInterface) {
		t.Fatalf("expected JSON to contain:\n%s\ngot:\n%s", expectedJSON, r.Body.String())
	}
}

// BodyEquals asserts the response body equals the expected string
func (r *Response) BodyEquals(t TestingTB, expected string) {
	if r.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, r.Body.String())
	}
}

// BodyContains asserts the response body contains the expected string
func (r *Response) BodyContains(t TestingTB, expected string) {
	if !strings.Contains(r.Body.String(), expected) {
		t.Fatalf("expected body to contain %q, got %q", expected, r.Body.String())
	}
}

// XML parses the response body as XML into v
func (r *Response) XML(v interface{}) error {
	return xml.NewDecoder(r.Body).Decode(v)
}

// JSONEq asserts the response body equals the expected XML
func (r *Response) XMLEq(t TestingTB, expected interface{}) {
	var actual interface{}
	if err := xml.Unmarshal(r.Body.Bytes(), &actual); err != nil {
		t.Fatalf("response body is not valid XML: %v", err)
	}
	expectedXML, _ := xml.Marshal(expected)
	var expectedInterface interface{}
	if err := xml.Unmarshal(expectedXML, &expectedInterface); err != nil {
		t.Fatalf("expected value is not valid XML: %v", err)
	}
	if !jsonEq(actual, expectedInterface) {
		t.Fatalf("expected XML:\n%s\ngot:\n%s", expectedXML, r.Body.String())
	}
}

// ServeHTTP is a helper to test http.Handler
func ServeHTTP(h http.Handler, method, path string, body interface{}) *Response {
	var bodyReader io.Reader
	if body != nil {
		if s, ok := body.(string); ok {
			bodyReader = strings.NewReader(s)
		} else {
			b, _ := json.Marshal(body)
			bodyReader = bytes.NewReader(b)
		}
	}
	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	w := NewResponse()
	h.ServeHTTP(w, req)
	return w
}

// TestingTB is the testing interface
type TestingTB interface {
	Fatalf(format string, args ...interface{})
}

// jsonEq compares two JSON representations
func jsonEq(a, b interface{}) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}

// jsonContains checks if a contains b
func jsonContains(a, b interface{}) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return strings.Contains(string(aJSON), string(bJSON))
}
