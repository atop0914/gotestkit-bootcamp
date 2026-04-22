package assert

import (
	"errors"
	"testing"
)

func TestEqual(t *testing.T) {
	Equal(t, 1, 1)
	Equal(t, "hello", "hello")
	Equal(t, []int{1, 2}, []int{1, 2})
}

func TestNotEqual(t *testing.T) {
	NotEqual(t, 1, 2)
	NotEqual(t, "hello", "world")
}

func TestNil(t *testing.T) {
	var ptr *string
	Nil(t, ptr)
	Nil(t, nil)
}

func TestNotNil(t *testing.T) {
	s := "hello"
	NotNil(t, s)
	NotNil(t, &s)
}

func TestContains(t *testing.T) {
	Contains(t, "hello world", "world")
	Contains(t, "foo bar baz", "bar")
}

func TestNotContains(t *testing.T) {
	NotContains(t, "hello", "world")
}

func TestTrue(t *testing.T) {
	True(t, true)
}

func TestFalse(t *testing.T) {
	False(t, false)
}

func TestError(t *testing.T) {
	err := errors.New("something went wrong")
	Error(t, err)
}

func TestNoError(t *testing.T) {
	NoError(t, nil)
}

func TestEqualError(t *testing.T) {
	err := errors.New("not found")
	EqualError(t, err, "not found")
}

func TestPanics(t *testing.T) {
	Panics(t, func() { panic("oops") })
}

func TestNotPanics(t *testing.T) {
	NotPanics(t, func() { /* no panic */ })
}

func TestLen(t *testing.T) {
	Len(t, []int{1, 2, 3}, 3)
	Len(t, map[string]int{"a": 1, "b": 2}, 2)
	Len(t, "hello", 5)
}

func TestJSONEq(t *testing.T) {
	JSONEq(t, `{"a":1}`, `{"a":1}`)
	JSONEq(t, `{"a":1}`, `{"a": 1}`)
}
