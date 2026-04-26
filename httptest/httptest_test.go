package httptest

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestNewResponse(t *testing.T) {
	r := NewResponse()
	if r == nil {
		t.Fatal("expected non-nil response")
	}
	if r.ResponseRecorder == nil {
		t.Fatal("expected non-nil ResponseRecorder")
	}
}

func TestGet(t *testing.T) {
	req := Get("/api/users")
	if req.Method != http.MethodGet {
		t.Errorf("expected GET, got %s", req.Method)
	}
	if req.URL.Path != "/api/users" {
		t.Errorf("expected path /api/users, got %s", req.URL.Path)
	}
}

func TestPost(t *testing.T) {
	req := Post("/api/users", map[string]string{"name": "test"})
	if req.Method != http.MethodPost {
		t.Errorf("expected POST, got %s", req.Method)
	}
}

func TestPut(t *testing.T) {
	req := Put("/api/users/1", map[string]string{"name": "updated"})
	if req.Method != http.MethodPut {
		t.Errorf("expected PUT, got %s", req.Method)
	}
}

func TestPatch(t *testing.T) {
	req := Patch("/api/users/1", map[string]string{"name": "patched"})
	if req.Method != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", req.Method)
	}
}

func TestHead(t *testing.T) {
	req := Head("/api/users")
	if req.Method != http.MethodHead {
		t.Errorf("expected HEAD, got %s", req.Method)
	}
}

func TestOptions(t *testing.T) {
	req := Options("/api/users")
	if req.Method != http.MethodOptions {
		t.Errorf("expected OPTIONS, got %s", req.Method)
	}
}

func TestDelete(t *testing.T) {
	req := Delete("/api/users/1")
	if req.Method != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", req.Method)
	}
}

func TestWithHeader(t *testing.T) {
	req := Get("/api/users")
	req = WithHeader(req, "X-Custom-Header", "custom-value")
	if req.Header.Get("X-Custom-Header") != "custom-value" {
		t.Errorf("expected header value custom-value, got %s", req.Header.Get("X-Custom-Header"))
	}
}

func TestWithHeaders(t *testing.T) {
	req := Get("/api/users")
	req = WithHeaders(req, map[string]string{
		"X-Header-1": "value1",
		"X-Header-2": "value2",
	})
	if req.Header.Get("X-Header-1") != "value1" {
		t.Errorf("expected header X-Header-1=value1, got %s", req.Header.Get("X-Header-1"))
	}
	if req.Header.Get("X-Header-2") != "value2" {
		t.Errorf("expected header X-Header-2=value2, got %s", req.Header.Get("X-Header-2"))
	}
}

func TestWithAuth(t *testing.T) {
	req := Get("/api/users")
	req = WithAuth(req, "user", "pass")
	username, password, ok := req.BasicAuth()
	if !ok {
		t.Fatal("expected basic auth to be set")
	}
	if username != "user" {
		t.Errorf("expected username user, got %s", username)
	}
	if password != "pass" {
		t.Errorf("expected password pass, got %s", password)
	}
}

func TestWithBearer(t *testing.T) {
	req := Get("/api/users")
	req = WithBearer(req, "my-token")
	if req.Header.Get("Authorization") != "Bearer my-token" {
		t.Errorf("expected Authorization header Bearer my-token, got %s", req.Header.Get("Authorization"))
	}
}

func TestWithQuery(t *testing.T) {
	req := Get("/api/users")
	req = WithQuery(req, map[string]string{"page": "1", "limit": "10"})
	if req.URL.RawQuery != "limit=10&page=1" {
		t.Errorf("expected query limit=10&page=1, got %s", req.URL.RawQuery)
	}
}

func TestWithForm(t *testing.T) {
	req := Post("/api/users", nil)
	req = WithForm(req, map[string]string{"name": "test", "email": "test@example.com"})
	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Errorf("expected Content-Type application/x-www-form-urlencoded, got %s", req.Header.Get("Content-Type"))
	}
}

func TestResponseString(t *testing.T) {
	r := NewResponse()
	r.Write([]byte("hello world"))
	if r.String() != "hello world" {
		t.Errorf("expected hello world, got %s", r.String())
	}
}

func TestResponseBytes(t *testing.T) {
	r := NewResponse()
	r.Write([]byte("hello"))
	bytes := r.Bytes()
	if string(bytes) != "hello" {
		t.Errorf("expected hello, got %s", string(bytes))
	}
}

func TestResponseOK(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusOK)
	r.OK(&failingTB{})
}

func TestResponseCreated(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusCreated)
	r.Created(&failingTB{})
}

func TestResponseAccepted(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusAccepted)
	r.Accepted(&failingTB{})
}

func TestResponseNoContent(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusNoContent)
	r.NoContent(&failingTB{})
}

func TestResponseBadRequest(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusBadRequest)
	r.BadRequest(&failingTB{})
}

func TestResponseUnauthorized(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusUnauthorized)
	r.Unauthorized(&failingTB{})
}

func TestResponseForbidden(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusForbidden)
	r.Forbidden(&failingTB{})
}

func TestResponseNotFound(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusNotFound)
	r.NotFound(&failingTB{})
}

func TestResponseConflict(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusConflict)
	r.Conflict(&failingTB{})
}

func TestResponseInternalServerError(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusInternalServerError)
	r.InternalServerError(&failingTB{})
}

func TestResponseStatusCode(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusOK)
	r.StatusCode(&failingTB{}, http.StatusOK)
}

func TestResponseHeader(t *testing.T) {
	r := NewResponse()
	r.ResponseRecorder.Header().Set("X-Custom", "value")
	if r.GetHeader("X-Custom") != "value" {
		t.Errorf("expected header value, got %s", r.GetHeader("X-Custom"))
	}
}

func TestResponseContentType(t *testing.T) {
	r := NewResponse()
	r.ResponseRecorder.Header().Set("Content-Type", "application/json")
	r.ContentType(&failingTB{}, "application/json")
}

func TestResponseJSON(t *testing.T) {
	r := NewResponse()
	json.NewEncoder(r.Body).Encode(map[string]string{"name": "test"})
	var result map[string]string
	err := r.JSON(&result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["name"] != "test" {
		t.Errorf("expected name=test, got %s", result["name"])
	}
}

func TestResponseJSONEq(t *testing.T) {
	r := NewResponse()
	json.NewEncoder(r.Body).Encode(map[string]string{"name": "test", "age": "25"})
	r.JSONEq(&failingTB{}, map[string]string{"name": "test", "age": "25"})
}

func TestResponseJSONContains(t *testing.T) {
	r := NewResponse()
	json.NewEncoder(r.Body).Encode(map[string]interface{}{
		"user": map[string]string{"name": "test"},
	})
	r.JSONContains(&failingTB{}, map[string]string{"name": "test"})
}

func TestResponseBodyEquals(t *testing.T) {
	r := NewResponse()
	r.Write([]byte("hello world"))
	r.BodyEquals(&failingTB{}, "hello world")
}

func TestResponseBodyContains(t *testing.T) {
	r := NewResponse()
	r.Write([]byte("hello world"))
	r.BodyContains(&failingTB{}, "world")
}

func TestResponseXMLEq(t *testing.T) {
	r := NewResponse()
	r.Write([]byte(`<user><name>test</name></user>`))
	type User struct {
		Name string `xml:"name"`
	}
	r.XMLEq(&failingTB{}, User{Name: "test"})
}

func TestServeHTTP(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	resp := ServeHTTP(handler, "GET", "/api/test", nil)
	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]string
	err := resp.JSON(&result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("expected status=ok, got %s", result["status"])
	}
}

func TestServeHTTPWithBody(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		var data map[string]string
		json.NewDecoder(r.Body).Decode(&data)
		json.NewEncoder(w).Encode(data)
	})

	resp := ServeHTTP(handler, "POST", "/api/users", map[string]string{"name": "test"})
	if resp.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.Code)
	}

	var result map[string]string
	err := resp.JSON(&result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["name"] != "test" {
		t.Errorf("expected name=test, got %s", result["name"])
	}
}

type failingTB struct{}

func (f *failingTB) Fatalf(format string, args ...interface{}) {}
