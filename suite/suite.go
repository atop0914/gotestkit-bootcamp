// Package suite provides testing suite functionality for organizing and running related tests.
package suite

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

// Suite is the interface that test suites must implement
type Suite interface {
	// SetupSuite is called once before any tests in the suite run
	SetupSuite()
	// TearDownSuite is called once after all tests in the suite finish
	TearDownSuite()
	// Setup is called before each test
	Setup()
	// TearDown is called after each test
	TearDown()
}

// TestSuite is a base implementation of Suite with default no-op methods
type TestSuite struct{}

func (s *TestSuite) SetupSuite()   {}
func (s *TestSuite) TearDownSuite() {}
func (s *TestSuite) Setup()         {}
func (s *TestSuite) TearDown()      {}

// runner manages the test suite execution
type runner struct {
	suite   Suite
	testing *testing.T
	results map[string]TestResult
}

// TestResult holds the result of a test run
type TestResult struct {
	Name      string
	Passed    bool
	Duration  time.Duration
	Error     error
	Skipped   bool
	SkipReason string
}

// RunSuite runs a test suite using Go's testing.T
func RunSuite(t *testing.T, suite Suite) {
	r := &runner{
		suite:   suite,
		testing: t,
		results: make(map[string]TestResult),
	}
	r.run()
}

// RunSubSuite runs a test suite as a sub-test
func RunSubSuite(t *testing.T, suite Suite) {
	r := &runner{
		suite:   suite,
		testing: t,
		results: make(map[string]TestResult),
	}
	r.runAsSubTest()
}

func (r *runner) run() {
	// Setup suite once
	defer func() {
		if r.panicked() {
			r.testing.Log("Panic during suite teardown, tests may be incomplete")
		}
	}()

	r.testing.Log("Setting up suite:", r.suiteName())
	r.suite.SetupSuite()
	r.testing.Log("Suite setup complete")

	// Run all test methods
	tests := r.collectTests()
	for _, test := range tests {
		r.runTest(test)
	}

	r.testing.Log("Tearing down suite")
	r.suite.TearDownSuite()
	r.testing.Log("Suite complete")

	// Log summary
	r.logSummary()
}

func (r *runner) runAsSubTest() {
	r.testing.Helper()

	// Use top-level helpers for setup/teardown
	r.suite.SetupSuite()

	tests := r.collectTests()
	for _, test := range tests {
		testName := fmt.Sprintf("%s/%s", r.suiteName(), test.Name)
		r.testing.Run(testName, func(st *testing.T) {
			r.runTestInSubTest(st, test)
		})
	}

	r.suite.TearDownSuite()
}

func (r *runner) runTest(method testMethod) {
	name := fmt.Sprintf("%s/%s", r.suiteName(), method.Name)
	r.testing.Run(name, func(st *testing.T) {
		r.runTestInSubTest(st, method)
	})
}

func (r *runner) runTestInSubTest(st *testing.T, method testMethod) {
	st.Helper()

	// Check if method is a teardown
	if method.IsTeardown {
		return
	}

	// Setup for this test
	r.suite.Setup()

	// Ensure teardown is called
	panicked := false
	defer func() {
		if r.panicked() {
			panicked = true
			st.Log("Panic detected during test or teardown")
		}
		if !panicked {
			r.suite.TearDown()
		}
	}()

	// Run the test
	start := time.Now()
	if method.IsAsync {
		r.runAsyncTest(st, method)
	} else {
		r.runSyncTest(st, method)
	}
	r.results[method.Name] = TestResult{
		Name:     method.Name,
		Passed:   !st.Failed() && !panicked,
		Duration: time.Since(start),
		Skipped:  st.Skipped(),
	}
}

func (r *runner) runSyncTest(st *testing.T, method testMethod) {
	st.Helper()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get suite instance as first argument (receiver)
	receiver := reflect.ValueOf(r.suite)

	// Call the method with or without context
	if method.AcceptsContext {
		method.Func.Call([]reflect.Value{
			receiver,
			reflect.ValueOf(ctx),
			reflect.ValueOf(st),
		})
	} else {
		method.Func.Call([]reflect.Value{
			receiver,
			reflect.ValueOf(st),
		})
	}
}

func (r *runner) runAsyncTest(st *testing.T, method testMethod) {
	done := make(chan struct{})
	receiver := reflect.ValueOf(r.suite)
	go func() {
		defer close(done)
		if method.AcceptsContext {
			ctx := context.Background()
			method.Func.Call([]reflect.Value{
				receiver,
				reflect.ValueOf(ctx),
				reflect.ValueOf(st),
			})
		} else {
			method.Func.Call([]reflect.Value{
				receiver,
				reflect.ValueOf(st),
			})
		}
	}()

	select {
	case <-done:
		// Test completed
	default:
		// Continue without blocking
	}
}

func (r *runner) collectTests() []testMethod {
	t := reflect.TypeOf(r.suite)
	methods := make([]testMethod, 0)

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		name := method.Name

		// Skip non-test methods and suite interface methods
		if !strings.HasPrefix(name, "Test") && !strings.HasPrefix(name, "Benchmark") {
			continue
		}
		if name == "SetupSuite" || name == "TearDownSuite" || name == "Setup" || name == "TearDown" {
			continue
		}

		tm := testMethod{
			Name:       name,
			Func:       method.Func,
			IsAsync:    strings.HasSuffix(name, "Async"),
			IsTeardown: strings.HasPrefix(name, "TearDown"),
		}

		// Check if method accepts context (check index 1 since index 0 is receiver)
		if method.Type.NumIn() >= 3 {
			firstArg := method.Type.In(1)
			tm.AcceptsContext = firstArg == reflect.TypeOf((*context.Context)(nil)).Elem()
		}

		methods = append(methods, tm)
	}

	return methods
}

func (r *runner) suiteName() string {
	t := reflect.TypeOf(r.suite)
	return t.Elem().Name()
}

func (r *runner) panicked() bool {
	return runtime.NumGoroutine() > 1
}

func (r *runner) logSummary() {
	var passed, failed, skipped int
	for _, result := range r.results {
		if result.Skipped {
			skipped++
		} else if result.Passed {
			passed++
		} else {
			failed++
		}
	}
	r.testing.Logf("Suite summary: %d passed, %d failed, %d skipped", passed, failed, skipped)
}

type testMethod struct {
	Name           string
	Func           reflect.Value
	IsAsync        bool
	IsTeardown     bool
	AcceptsContext bool
}

// Skippable is implemented by tests that can be skipped
type Skippable interface {
	Skip(reason string)
}

// WithTimeout sets a custom timeout for a test
func WithTimeout(timeout time.Duration) func(*testing.T) {
	return func(t *testing.T) {
		// This is a placeholder for timeout functionality
		_ = timeout
	}
}
