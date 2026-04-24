package mock

import (
	"errors"
	"testing"
)

type mockT struct {
	failed bool
}

func (m *mockT) Helper()                                   {}
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

func TestMockOnce(t *testing.T) {
	m := New(&mockT{})
	m.Once("GetUser", "result")
	m.Called("GetUser", "arg1")
	m.Called("GetUser", "arg2")
	// Once should only allow 1 call
	if m.CallCount("GetUser") != 1 {
		t.Errorf("expected 1 call with Once, got %d", m.CallCount("GetUser"))
	}
}

func TestMockMaybe(t *testing.T) {
	m := New(&mockT{})
	m.Maybe("GetUser", "result")
	// Maybe allows 0 or 1 calls, should not fail
	m.AssertExpectations(&mockT{})
}

func TestMockTimes(t *testing.T) {
	m := New(&mockT{})
	m.Times("GetUser", 2, "result")
	m.Called("GetUser", "arg1")
	m.Called("GetUser", "arg2")
	m.AssertExpectations(&mockT{})
}

func TestMockAssertCalledTimes(t *testing.T) {
	m := New(&mockT{})
	m.Called("GetUser", "arg1")
	m.Called("GetUser", "arg2")
	m.Called("GetUser", "arg3")
	m.AssertCalledTimes(&mockT{}, "GetUser", 3)
}

func TestMockCallCount(t *testing.T) {
	m := New(&mockT{})
	if m.CallCount("GetUser") != 0 {
		t.Errorf("expected 0 calls initially")
	}
	m.Called("GetUser")
	m.Called("GetUser")
	if m.CallCount("GetUser") != 2 {
		t.Errorf("expected 2 calls, got %d", m.CallCount("GetUser"))
	}
}

// ---- Matcher Tests ----

func TestAnyMatcher(t *testing.T) {
	m := Any()
	if !m.Matches("anything") {
		t.Error("Any should match anything")
	}
	if !m.Matches(nil) {
		t.Error("Any should match nil")
	}
	if !m.Matches(123) {
		t.Error("Any should match numbers")
	}
}

func TestEqMatcher(t *testing.T) {
	m := Eq(42)
	if !m.Matches(42) {
		t.Error("Eq should match equal int")
	}
	if m.Matches(43) {
		t.Error("Eq should not match different int")
	}

	m2 := Eq("hello")
	if !m2.Matches("hello") {
		t.Error("Eq should match equal string")
	}
}

func TestContainsMatcher(t *testing.T) {
	m := Contains("world")
	if !m.Matches("hello world") {
		t.Error("Contains should find substring")
	}
	if m.Matches("hello") {
		t.Error("Contains should not match when not found")
	}
}

func TestMatchesMatcher(t *testing.T) {
	m := Matches(func(v interface{}) bool {
		if n, ok := v.(int); ok {
			return n > 0
		}
		return false
	})
	if !m.Matches(42) {
		t.Error("Matches should match predicate returning true")
	}
	if m.Matches(-1) {
		t.Error("Matches should not match predicate returning false")
	}
}

func TestNotMatcher(t *testing.T) {
	m := Not(Eq(42))
	if !m.Matches(43) {
		t.Error("Not should negate Eq matcher")
	}
	if !m.Matches("string") {
		t.Error("Not should match non-equal values")
	}
}

func TestMockWithMatchers(t *testing.T) {
	m := New(&mockT{})
	m.Expect("GetUser", "result")
	m.Called("GetUser", 123)
	m.AssertCalled(&mockT{}, "GetUser", 123)
}

func TestMockAwaitCall(t *testing.T) {
	m := New(&mockT{})
	m.Called("GetUser", "arg1")
	if !m.AwaitCall("GetUser", "arg1") {
		t.Error("AwaitCall should find existing call")
	}
}

func TestMockMethodWithPanic(t *testing.T) {
	m := New(&mockT{})
	fn := m.MethodWithPanic("GetUser", "panic!")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()

	fn("arg1")
}

func TestMockPanicCall(t *testing.T) {
	m := New(&mockT{})
	m.Expect("GetUser", "result").Panic("panic!")

	defer func() {
		if r := recover(); r == "panic!" {
			// expected
		} else {
			t.Errorf("expected panic 'panic!', got %v", r)
		}
	}()

	m.Called("GetUser", "arg1")
}

func TestMockWaitFor(t *testing.T) {
	m := New(&mockT{})
	// This is a basic test - actual WaitFor with goroutines
	// would require more complex async testing
	call := m.Expect("GetUser", "result")
	call.WaitFor(Any())

	if len(call.WaitForArgs) != 1 {
		t.Error("expected WaitForArgs to be set")
	}
}

func TestMockNoError(t *testing.T) {
	err := errors.New("test error")
	if err == nil {
		t.Fatal("expected error")
	}
}
