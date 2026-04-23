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

// Greater asserts that actual is greater than expected
func Greater(t T, expected, actual int, msg ...string) {
	if actual <= expected {
		fail(t, "Greater", message(msg...), expected, actual)
	}
}

// GreaterOrEqual asserts that actual is greater than or equal to expected
func GreaterOrEqual(t T, expected, actual int, msg ...string) {
	if actual < expected {
		fail(t, "GreaterOrEqual", message(msg...), expected, actual)
	}
}

// Less asserts that actual is less than expected
func Less(t T, expected, actual int, msg ...string) {
	if actual >= expected {
		fail(t, "Less", message(msg...), expected, actual)
	}
}

// LessOrEqual asserts that actual is less than or equal to expected
func LessOrEqual(t T, expected, actual int, msg ...string) {
	if actual > expected {
		fail(t, "LessOrEqual", message(msg...), expected, actual)
	}
}

// Empty asserts that a string, slice, map, etc. is empty
func Empty(t T, actual interface{}, msg ...string) {
	if !isEmpty(actual) {
		fail(t, "Empty", message(msg...), "empty", actual)
	}
}

// NotEmpty asserts that a string, slice, map, etc. is not empty
func NotEmpty(t T, actual interface{}, msg ...string) {
	if isEmpty(actual) {
		fail(t, "NotEmpty", message(msg...), "not empty", actual)
	}
}

// TypeOf asserts that the actual value is of the expected type
func TypeOf(t T, expected interface{}, actual interface{}, msg ...string) {
	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)
	if expectedType != actualType {
		fail(t, "TypeOf", message(msg...), expectedType, actualType)
	}
}

// Implements asserts that actual implements the expected interface
func Implements(t T, expectedInterface interface{}, actual interface{}, msg ...string) {
	var expectedType reflect.Type
	if et, ok := expectedInterface.(reflect.Type); ok {
		expectedType = et
	} else {
		expectedType = reflect.TypeOf(expectedInterface)
	}
	if expectedType == nil {
		return
	}
	actualType := reflect.TypeOf(actual)
	if actualType == nil {
		fail(t, "Implements", message(msg...), expectedType, nil)
		return
	}
	if !actualType.Implements(expectedType) {
		fail(t, "Implements", message(msg...), expectedType, actualType)
	}
}

// PanicsWithValue asserts that the function panics with the expected value
func PanicsWithValue(t T, expected interface{}, fn func(), msg ...string) {
	defer func() {
		r := recover()
		if r == nil {
			fail(t, "PanicsWithValue", message(msg...), expected, "no panic")
			return
		}
		if !reflect.DeepEqual(expected, r) {
			fail(t, "PanicsWithValue", message(msg...), expected, r)
		}
	}()
	fn()
}

// ElementsMatch asserts that two slices contain the same elements (order independent)
func ElementsMatch(t T, expected, actual interface{}, msg ...string) {
	expVal := reflect.ValueOf(expected)
	actVal := reflect.ValueOf(actual)

	if expVal.Kind() != reflect.Slice || actVal.Kind() != reflect.Slice {
		fail(t, "ElementsMatch", message(msg...), "slice", reflect.ValueOf(actual).Kind())
		return
	}

	expElems := expandAndNormalize(expVal)
	actElems := expandAndNormalize(actVal)

	if len(expElems) != len(actElems) {
		fail(t, "ElementsMatch", message(msg...), expElems, actElems)
		return
	}

	for _, e := range expElems {
		found := false
		for _, a := range actElems {
			if reflect.DeepEqual(e, a) {
				found = true
				break
			}
		}
		if !found {
			fail(t, "ElementsMatch", message(msg...), expElems, actElems)
			return
		}
	}
}

// JSONContains asserts that the JSON string contains the expected JSON fragment
func JSONContains(t T, jsonStr, expectedFragment string, msg ...string) {
	var expected, actual interface{}
	if err := json.Unmarshal([]byte(expectedFragment), &expected); err != nil {
		fail(t, "JSONContains", message(msg...), expectedFragment, "invalid JSON fragment")
		return
	}
	if err := json.Unmarshal([]byte(jsonStr), &actual); err != nil {
		fail(t, "JSONContains", message(msg...), expectedFragment, "invalid JSON string")
		return
	}
	if !containsJSON(expected, actual) {
		fail(t, "JSONContains", message(msg...), expectedFragment, jsonStr)
	}
}

// Subset asserts that the actual slice contains all expected elements
func Subset(t T, expected, actual interface{}, msg ...string) {
	expVal := reflect.ValueOf(expected)
	actVal := reflect.ValueOf(actual)

	if expVal.Kind() != reflect.Slice || actVal.Kind() != reflect.Slice {
		fail(t, "Subset", message(msg...), "slice", reflect.ValueOf(actual).Kind())
		return
	}

	for i := 0; i < expVal.Len(); i++ {
		found := false
		for j := 0; j < actVal.Len(); j++ {
			if reflect.DeepEqual(expVal.Index(i).Interface(), actVal.Index(j).Interface()) {
				found = true
				break
			}
		}
		if !found {
			fail(t, "Subset", message(msg...), expected, actual)
			return
		}
	}
}

// Zero asserts that value is zero value
func Zero(t T, actual interface{}, msg ...string) {
	if !isZero(actual) {
		fail(t, "Zero", message(msg...), "zero value", actual)
	}
}

// NotZero asserts that value is not zero value
func NotZero(t T, actual interface{}, msg ...string) {
	if isZero(actual) {
		fail(t, "NotZero", message(msg...), "non-zero value", actual)
	}
}

// ---- Helper functions ----

func isEmpty(v interface{}) bool {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return rv.Len() == 0
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return true
		}
		return isEmpty(rv.Elem().Interface())
	default:
		return false
	}
}

func expandAndNormalize(v reflect.Value) []interface{} {
	var elems []interface{}
	for i := 0; i < v.Len(); i++ {
		elems = append(elems, v.Index(i).Interface())
	}
	return elems
}

func containsJSON(expected, actual interface{}) bool {
	if reflect.DeepEqual(expected, actual) {
		return true
	}
	// Simple containment check for nested objects
	switch exp := expected.(type) {
	case map[string]interface{}:
		if act, ok := actual.(map[string]interface{}); ok {
			for k, v := range exp {
				if av, exists := act[k]; exists {
					if !containsJSON(v, av) {
						return false
					}
				} else {
					return false
				}
			}
			return true
		}
	}
	return false
}

func isZero(v interface{}) bool {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return true
		}
		return isZero(rv.Elem().Interface())
	case reflect.Struct:
		// Check if all fields are zero
		for i := 0; i < rv.NumField(); i++ {
			if !isZero(rv.Field(i).Interface()) {
				return false
			}
		}
		return true
	case reflect.Array, reflect.Slice:
		return rv.Len() == 0
	case reflect.Map:
		return rv.Len() == 0
	case reflect.String:
		return rv.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Bool:
		return !rv.Bool()
	default:
		return false
	}
}
