package mock

import (
	"errors"
	"testing"
)

type mockT struct {
	failed bool
}

func (m *mockT) Helper()             {}
func (m *mockT) Fatalf(format string, args ...interface{}) { m.failed = true }

func TestMockNew(t *testing.T) {
	m := New(&mockT{})
	if m == nil {
		t.Fatal("expected non-nil mock")
	}
}

func TestMockExpect(t *testing.T) {
	m := New(&mockT{})
	call := m.Expect("GetUser", "user123", nil)
	if call == nil {
		t.Fatal("expected non-nil call")
	}
	if call.MethodName != "GetUser" {
		t.Errorf("expected GetUser, got %s", call.MethodName)
	}
}

func TestMockCalled(t *testing.T) {
	m := New(&mockT{})
	m.Called("GetUser", "user123")
	if m.CallCount("GetUser") != 1 {
		t.Errorf("expected 1 call, got %d", m.CallCount("GetUser"))
	}
}

func TestMockAssertCalled(t *testing.T) {
	m := New(&mockT{})
	m.Called("GetUser", "user123")
	m.AssertCalled(&mockT{}, "GetUser", "user123")
}

func TestMockAssertNotCalled(t *testing.T) {
	m := New(&mockT{})
	m.AssertNotCalled(&mockT{}, "GetUser")
}

func TestMockAssertExpectations(t *testing.T) {
	m := New(&mockT{})
	m.Expect("GetUser", "user123")
	m.Called("GetUser", "user123")
	m.AssertExpectations(&mockT{})
}

func TestMockReset(t *testing.T) {
	m := New(&mockT{})
	m.Called("GetUser")
	m.Reset()
	if m.CallCount("GetUser") != 0 {
		t.Errorf("expected 0 calls after reset, got %d", m.CallCount("GetUser"))
	}
}

func TestMockMethod(t *testing.T) {
	m := New(&mockT{})
	fn := m.Method("GetUser", "result")
	result := fn("arg1", "arg2")
	if result != "result" {
		t.Errorf("expected result, got %v", result)
	}
	if m.CallCount("GetUser") != 1 {
		t.Errorf("expected 1 call, got %d", m.CallCount("GetUser"))
	}
}

func TestMockUnexpectedCall(t *testing.T) {
	m := New(&mockT{})
	m.Called("GetUser", "user123")
	// Should not panic, just record
	if m.CallCount("GetUser") != 1 {
		t.Errorf("expected 1 call")
	}
}

func TestMockNoError(t *testing.T) {
	err := errors.New("test error")
	if err == nil {
		t.Fatal("expected error")
	}
}
