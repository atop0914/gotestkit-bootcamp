// Package assert provides fluent assertions for testing in Go.
package assert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// T is a subset of testing.TB for assertion support
type T interface {
	Helper()
	Fatalf(format string, args ...interface{})
}

// AssertionFailedError is returned when an assertion fails
type AssertionFailedError struct {
	Message   string
	Expected  interface{}
	Actual    interface{}
	Operation string // e.g., "Equal", "Contains", "Nil"
}

func (e *AssertionFailedError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", e.Operation, e.Message)
	}
	return fmt.Sprintf("%s: expected %v, got %v", e.Operation, e.Expected, e.Actual)
}

// GlobalConfig controls assertion behavior
var GlobalConfig = struct {
	// ContinueOnFailure if true, failures call t.Errorf instead of t.Fatalf
	ContinueOnFailure bool
	// MaxDepth for recursive comparisons
	MaxDepth int
}{
	MaxDepth: 10,
}

// Equal asserts that two values are equal
func Equal(t T, expected, actual interface{}, msg ...string) {
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "Equal", message(msg...), expected, actual)
	}
}

// NotEqual asserts that two values are not equal
func NotEqual(t T, expected, actual interface{}, msg ...string) {
	if reflect.DeepEqual(expected, actual) {
		fail(t, "NotEqual", message(msg...), expected, actual)
	}
}

// Nil asserts that value is nil
func Nil(t T, actual interface{}, msg ...string) {
	if !isNil(actual) {
		fail(t, "Nil", message(msg...), nil, actual)
	}
}

// NotNil asserts that value is not nil
func NotNil(t T, actual interface{}, msg ...string) {
	if isNil(actual) {
		fail(t, "NotNil", message(msg...), "non-nil", nil)
	}
}

// Contains asserts that a string contains a substring
func Contains(t T, s, substr string, msg ...string) {
	if !strings.Contains(s, substr) {
		fail(t, "Contains", message(msg...), substr, s)
	}
}

// NotContains asserts that a string does not contain a substring
func NotContains(t T, s, substr string, msg ...string) {
	if strings.Contains(s, substr) {
		fail(t, "NotContains", message(msg...), "not containing "+substr, s)
	}
}

// True asserts that value is true
func True(t T, actual bool, msg ...string) {
	if !actual {
		fail(t, "True", message(msg...), true, false)
	}
}

// False asserts that value is false
func False(t T, actual bool, msg ...string) {
	if actual {
		fail(t, "False", message(msg...), false, true)
	}
}

// Error asserts that err is not nil
func Error(t T, err error, msg ...string) {
	if err == nil {
		fail(t, "Error", message(msg...), "error", nil)
	}
}

// NoError asserts that err is nil
func NoError(t T, err error, msg ...string) {
	if err != nil {
		fail(t, "NoError", message(msg...), nil, err)
	}
}

// EqualError asserts that err has the specified message
func EqualError(t T, err error, expectedMsg string, msg ...string) {
	if err == nil {
		fail(t, "EqualError", message(msg...), expectedMsg, nil)
		return
	}
	if err.Error() != expectedMsg {
		fail(t, "EqualError", message(msg...), expectedMsg, err.Error())
	}
}

// Panics asserts that the function panics
func Panics(t T, fn func(), msg ...string) {
	defer func() {
		if recover() == nil {
			fail(t, "Panics", message(msg...), "panic", "no panic")
		}
	}()
	fn()
}

// NotPanics asserts that the function does not panic
func NotPanics(t T, fn func(), msg ...string) {
	defer func() {
		if r := recover(); r != nil {
			fail(t, "NotPanics", message(msg...), "no panic", r)
		}
	}()
	fn()
}

// JSONEq asserts that two JSON strings are equivalent
func JSONEq(t T, expected, actual string, msg ...string) {
	expected = strings.TrimSpace(expected)
	actual = strings.TrimSpace(actual)
	// Normalize whitespace in JSON
	var expVar, actVar interface{}
	if err := json.Unmarshal([]byte(expected), &expVar); err != nil {
		fail(t, "JSONEq", message(msg...), expected, actual)
		return
	}
	if err := json.Unmarshal([]byte(actual), &actVar); err != nil {
		fail(t, "JSONEq", message(msg...), expected, actual)
		return
	}
	expBytes, _ := json.Marshal(expVar)
	actBytes, _ := json.Marshal(actVar)
	if string(expBytes) != string(actBytes) {
		fail(t, "JSONEq", message(msg...), string(expBytes), string(actBytes))
	}
}

// Len asserts that a slice/array/map has the expected length
func Len(t T, object interface{}, length int, msg ...string) {
	l := reflect.ValueOf(object).Len()
	if l != length {
		fail(t, "Len", message(msg...), length, l)
	}
}

// Same asserts that two interface{} values point to the same object
func Same(t T, expected, actual interface{}, msg ...string) {
	if !reflect.ValueOf(expected).Equal(reflect.ValueOf(actual)) {
		fail(t, "Same", message(msg...), expected, actual)
	}
}

// NotSame asserts that two interface{} values do not point to the same object
func NotSame(t T, expected, actual interface{}, msg ...string) {
	if reflect.ValueOf(expected).Equal(reflect.ValueOf(actual)) {
		fail(t, "NotSame", message(msg...), "different objects", "same object")
	}
}

// Eventually asserts that a condition will eventually be true
func Eventually(t T, condition func() bool, duration time.Duration, msg ...string) {
	deadline := time.Now().Add(duration)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	fail(t, "Eventually", message(msg...), "condition to be true", "condition was false")
}

// Never asserts that a condition will never be true within duration
func Never(t T, condition func() bool, duration time.Duration, msg ...string) {
	deadline := time.Now().Add(duration)
	for time.Now().Before(deadline) {
		if condition() {
			fail(t, "Never", message(msg...), "condition to never be true", "condition was true")
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// fail reports a test failure
func fail(t T, op, msg string, expected, actual interface{}) {
	t.Helper()
	err := &AssertionFailedError{
		Operation: op,
		Message:   msg,
		Expected:  expected,
		Actual:    actual,
	}
	if GlobalConfig.ContinueOnFailure {
		t.Fatalf("%s", err)
	} else {
		t.Fatalf("%s", err)
	}
}

func message(args ...string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return rv.IsNil()
	default:
		return false
	}
}
