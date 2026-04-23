package assert

import (
	"errors"
	"reflect"
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

func TestGreater(t *testing.T) {
	Greater(t, 1, 5)
	Greater(t, -1, 0)
}

func TestGreaterOrEqual(t *testing.T) {
	GreaterOrEqual(t, 1, 1)
	GreaterOrEqual(t, 1, 5)
}

func TestLess(t *testing.T) {
	Less(t, 5, 1)
	Less(t, 0, -1)
}

func TestLessOrEqual(t *testing.T) {
	LessOrEqual(t, 1, 1)
	LessOrEqual(t, 5, 1)
}

func TestEmpty(t *testing.T) {
	Empty(t, "")
	Empty(t, []int{})
	Empty(t, map[string]int{})
}

func TestNotEmpty(t *testing.T) {
	NotEmpty(t, "hello")
	NotEmpty(t, []int{1})
	NotEmpty(t, map[string]int{"a": 1})
}

func TestZero(t *testing.T) {
	Zero(t, 0)
	Zero(t, "")
	Zero(t, false)
}

func TestNotZero(t *testing.T) {
	NotZero(t, 1)
	NotZero(t, "hello")
	NotZero(t, true)
}

func TestPanicsWithValue(t *testing.T) {
	PanicsWithValue(t, "oops", func() { panic("oops") })
	PanicsWithValue(t, 42, func() { panic(42) })
}

func TestElementsMatch(t *testing.T) {
	ElementsMatch(t, []int{1, 2, 3}, []int{3, 2, 1})
	ElementsMatch(t, []string{"a", "b"}, []string{"b", "a"})
}

func TestSubset(t *testing.T) {
	Subset(t, []int{1, 2}, []int{1, 2, 3})
	Subset(t, []string{"a"}, []string{"a", "b", "c"})
}

type myError struct{}

func (myError) Error() string { return "error" }

func TestTypeOf(t *testing.T) {
	TypeOf(t, 1, 1)
	TypeOf(t, "hello", "world")
}

func TestImplements(t *testing.T) {
	var e error = myError{}
	Implements(t, reflect.TypeOf((*error)(nil)).Elem(), e)
}
