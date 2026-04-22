// Package mock provides a simple mocking framework for Go tests.
package mock

import (
	"reflect"
	"sync"
)

// Mock represents a mock object
type Mock struct {
	mu       sync.RWMutex
	calls    map[string][]*Call
	expected map[string]*Call
	t        TB
}

type TB interface {
	Helper()
	Fatalf(format string, args ...interface{})
}

// New creates a new Mock
func New(t TB) *Mock {
	return &Mock{
		calls:    make(map[string][]*Call),
		expected: make(map[string]*Call),
		t:        t,
	}
}

// Call represents a single method call
type Call struct {
	MethodName string
	Args       []interface{}
	ReturnArgs []interface{}
	CallCount  int32
	Called     bool
}

// Expect sets up an expected call with return values
func (m *Mock) Expect(methodName string, returnArgs ...interface{}) *Call {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := &Call{
		MethodName: methodName,
		ReturnArgs: returnArgs,
	}
	m.expected[methodName] = call
	return call
}

// Called records that a method was called
func (m *Mock) Called(methodName string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := m.expected[methodName]
	if call == nil {
		// Record unexpected call
		call = &Call{MethodName: methodName, Args: args}
		m.calls[methodName] = append(m.calls[methodName], call)
		return
	}

	call.CallCount++
	call.Called = true
	call.Args = args

	// Store call for assertion
	m.calls[methodName] = append(m.calls[methodName], call)
}

// AssertCalled verifies a method was called
func (m *Mock) AssertCalled(t TB, methodName string, args ...interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := m.calls[methodName]
	if len(calls) == 0 {
		t.Fatalf("mock: expected %s to be called, but it was not", methodName)
	}

	lastCall := calls[len(calls)-1]
	if len(args) > 0 {
		for i, arg := range args {
			if len(lastCall.Args) <= i {
				t.Fatalf("mock: call %s missing arg %d", methodName, i)
			}
			if !safeEqual(arg, lastCall.Args[i]) {
				t.Fatalf("mock: call %s arg %d mismatch\nexpected: %v\nactual: %v", methodName, i, arg, lastCall.Args[i])
			}
		}
	}
}

// AssertNotCalled verifies a method was never called
func (m *Mock) AssertNotCalled(t TB, methodName string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := m.calls[methodName]
	if len(calls) > 0 {
		t.Fatalf("mock: expected %s to not be called, but it was called %d times", methodName, len(calls))
	}
}

// AssertExpectations verifies all expected calls were made
func (m *Mock) AssertExpectations(t TB) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, call := range m.expected {
		if !call.Called {
			t.Fatalf("mock: expected %s was not called", name)
		}
	}
}

// CallCount returns how many times a method was called
func (m *Mock) CallCount(methodName string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := m.calls[methodName]
	return len(calls)
}

// Method creates a method stub that records calls
func (m *Mock) Method(methodName string, returnArgs ...interface{}) func(...interface{}) interface{} {
	return func(args ...interface{}) interface{} {
		m.Called(methodName, args...)
		if len(returnArgs) > 0 {
			return returnArgs[0]
		}
		return nil
	}
}

// AwaitCall waits for a specific call with timeout
func (m *Mock) AwaitCall(methodName string, args ...interface{}) bool {
	attempts := 0
	maxAttempts := 100
	for attempts < maxAttempts {
		m.mu.RLock()
		calls := m.calls[methodName]
		m.mu.RUnlock()

		for _, call := range calls {
			if len(args) == 0 || safeEqualSlice(args, call.Args) {
				return true
			}
		}
		attempts++
	}
	return false
}

// Reset clears all recorded calls
func (m *Mock) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = make(map[string][]*Call)
	m.expected = make(map[string]*Call)
}

func safeEqual(a, b interface{}) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	// Use reflect to compare
	return reflect.DeepEqual(a, b)
}

func safeEqualSlice(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !safeEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}
