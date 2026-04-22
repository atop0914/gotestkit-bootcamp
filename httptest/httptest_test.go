package httptest

import (
	"net/http"
	"testing"
)

func TestNewResponse(t *testing.T) {
	r := NewResponse()
	if r == nil {
		t.Fatal("expected non-nil response")
	}
}

func TestGet(t *testing.T) {
	req := Get("/api/users")
	if req.Method != http.MethodGet {
		t.Errorf("expected GET, got %s", req.Method)
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

func TestDelete(t *testing.T) {
	req := Delete("/api/users/1")
	if req.Method != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", req.Method)
	}
}

func TestResponseStatusCode(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusOK)
	r.StatusCode(t, http.StatusOK)
}

type failingTB struct{}

func (f *failingTB) Fatalf(format string, args ...interface{}) {}

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

func TestResponseBadRequest(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusBadRequest)
	r.BadRequest(&failingTB{})
}

func TestResponseNotFound(t *testing.T) {
	r := NewResponse()
	r.WriteHeader(http.StatusNotFound)
	r.NotFound(&failingTB{})
}
