// Package mock provides a simple mocking framework for Go tests.
package mock

import (
	"reflect"
	"sync"
	"sync/atomic"
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
	MethodName  string
	Args        []interface{}
	ReturnArgs  []interface{}
	CallCount   int32
	Called      bool
	Repeatable  bool // can be called multiple times
	MinCalls    int32
	MaxCalls    int32 // -1 means unlimited
	PanicValue  interface{}
	WaitForArgs []Matcher
}

// Expect sets up an expected call with return values
func (m *Mock) Expect(methodName string, returnArgs ...interface{}) *Call {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := &Call{
		MethodName: methodName,
		ReturnArgs: returnArgs,
		Repeatable: true,
		MaxCalls:   -1,
	}
	m.expected[methodName] = call
	return call
}

// Once returns a call that must be called exactly once
func (m *Mock) Once(methodName string, returnArgs ...interface{}) *Call {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := &Call{
		MethodName: methodName,
		ReturnArgs: returnArgs,
		Repeatable: false,
		MaxCalls:   1,
	}
	m.expected[methodName] = call
	return call
}

// Maybe marks an expected call as optional (may be called 0 or 1 times)
func (m *Mock) Maybe(methodName string, returnArgs ...interface{}) *Call {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := &Call{
		MethodName: methodName,
		ReturnArgs: returnArgs,
		Repeatable: true,
		MinCalls:   0,
		MaxCalls:   1,
	}
	m.expected[methodName] = call
	return call
}

// Times restricts how many times a method can be called
func (m *Mock) Times(methodName string, n int, returnArgs ...interface{}) *Call {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := &Call{
		MethodName: methodName,
		ReturnArgs: returnArgs,
		Repeatable: false,
		MinCalls:   int32(n),
		MaxCalls:   int32(n),
	}
	m.expected[methodName] = call
	return call
}

// WaitFor sets argument matchers that must be satisfied before this call proceeds
func (c *Call) WaitFor(matchers ...Matcher) *Call {
	c.WaitForArgs = matchers
	return c
}

// Panic configures this call to panic with the given value
func (c *Call) Panic(value interface{}) *Call {
	c.PanicValue = value
	return c
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

	// Check if call satisfies argument matchers
	if len(call.WaitForArgs) > 0 {
		if !m.matchArgs(call.WaitForArgs, args) {
			// Re-record for retry later
			call = &Call{MethodName: methodName, Args: args}
			m.calls[methodName] = append(m.calls[methodName], call)
			return
		}
	}

	// Check MaxCalls constraint
	if call.MaxCalls >= 0 {
		if atomic.LoadInt32(&call.CallCount) >= call.MaxCalls {
			// Already satisfied max calls, just record the unexpected call info but don't add to calls list
			return
		}
	}

	// Handle panic
	if call.PanicValue != nil {
		panic(call.PanicValue)
	}

	call.CallCount++
	call.Called = true
	call.Args = args

	// Store call for assertion
	m.calls[methodName] = append(m.calls[methodName], call)
}

func (m *Mock) matchArgs(matchers []Matcher, args []interface{}) bool {
	if len(matchers) != len(args) {
		return false
	}
	for i, matcher := range matchers {
		if i >= len(args) || !matcher.Matches(args[i]) {
			return false
		}
	}
	return true
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
		if call.MinCalls > 0 && atomic.LoadInt32(&call.CallCount) < call.MinCalls {
			t.Fatalf("mock: expected %s to be called at least %d times, got %d", name, call.MinCalls, call.CallCount)
		}
		if !call.Repeatable && !call.Called {
			t.Fatalf("mock: expected %s was not called", name)
		}
	}
}

// AssertCalledTimes verifies a method was called exactly n times
func (m *Mock) AssertCalledTimes(t TB, methodName string, n int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := m.callCountUnlocked(methodName)
	if count != n {
		t.Fatalf("mock: expected %s to be called %d times, got %d", methodName, n, count)
	}
}

// callCountUnlocked returns the call count (must be called with lock held)
func (m *Mock) callCountUnlocked(methodName string) int {
	calls := m.calls[methodName]
	return len(calls)
}

// CallCount returns how many times a method was called
func (m *Mock) CallCount(methodName string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.callCountUnlocked(methodName)
}

// Method creates a method stub that records calls and returns configured values
func (m *Mock) Method(methodName string, returnArgs ...interface{}) func(...interface{}) interface{} {
	return func(args ...interface{}) interface{} {
		m.Called(methodName, args...)
		if len(returnArgs) > 0 {
			return returnArgs[0]
		}
		return nil
	}
}

// MethodWithPanic creates a method stub that panics when called
func (m *Mock) MethodWithPanic(methodName string, panicValue interface{}) func(...interface{}) {
	return func(args ...interface{}) {
		m.Called(methodName, args...)
		if panicValue != nil {
			panic(panicValue)
		}
		panic("mock method called")
	}
}

// AwaitCall waits for a specific call with timeout (in iterations)
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

// AwaitCallWithMatcher waits for a call matching the given matchers
func (m *Mock) AwaitCallWithMatcher(methodName string, matchers ...Matcher) bool {
	attempts := 0
	maxAttempts := 100
	for attempts < maxAttempts {
		m.mu.RLock()
		calls := m.calls[methodName]
		m.mu.RUnlock()

		for _, call := range calls {
			if m.matchArgs(matchers, call.Args) {
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

// ---- Argument Matchers ----

// Matcher interface for argument matching
type Matcher interface {
	Matches(arg interface{}) bool
	String() string
}

// Any returns a matcher that matches any value
func Any() Matcher {
	return &anyMatcher{}
}

type anyMatcher struct{}

func (a *anyMatcher) Matches(arg interface{}) bool {
	return true
}

func (a *anyMatcher) String() string {
	return "<any>"
}

// Eq returns a matcher that matches if arg equals the expected value
func Eq(expected interface{}) Matcher {
	return &eqMatcher{expected: expected}
}

type eqMatcher struct {
	expected interface{}
}

func (e *eqMatcher) Matches(arg interface{}) bool {
	return safeEqual(e.expected, arg)
}

func (e *eqMatcher) String() string {
	return "eq(" + toString(e.expected) + ")"
}

// Contains returns a matcher that matches if arg contains the expected value (for strings/slices)
func Contains(expected interface{}) Matcher {
	return &containsMatcher{expected: expected}
}

type containsMatcher struct {
	expected interface{}
}

func (c *containsMatcher) Matches(arg interface{}) bool {
	if argStr, ok := arg.(string); ok {
		if expStr, ok := c.expected.(string); ok {
			return containsString(argStr, expStr)
		}
	}
	return safeEqual(arg, c.expected)
}

func (c *containsMatcher) String() string {
	return "contains(" + toString(c.expected) + ")"
}

// Matches returns a matcher that matches using a function
func Matches(fn func(interface{}) bool) Matcher {
	return &fnMatcher{fn: fn}
}

type fnMatcher struct {
	fn func(interface{}) bool
}

func (f *fnMatcher) Matches(arg interface{}) bool {
	return f.fn(arg)
}

func (f *fnMatcher) String() string {
	return "<fn>"
}

// Not returns a matcher that negates the given matcher
func Not(m Matcher) Matcher {
	return &notMatcher{inner: m}
}

type notMatcher struct {
	inner Matcher
}

func (n *notMatcher) Matches(arg interface{}) bool {
	return !n.inner.Matches(arg)
}

func (n *notMatcher) String() string {
	return "not(" + n.inner.String() + ")"
}

// ---- Helpers ----

func safeEqual(a, b interface{}) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
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

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toString(v interface{}) string {
	if v == nil {
		return "nil"
	}
	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return "int"
	case uint, uint8, uint16, uint32, uint64:
		return "uint"
	case float32, float64:
		return "float"
	default:
		return "value"
	}
}
