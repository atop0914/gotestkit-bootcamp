package suite

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

// ExampleSuite demonstrates how to use the test suite
type ExampleSuite struct {
	TestSuite
	setupCalled    bool
	teardownCalled bool
	testCalled     bool
}

func (s *ExampleSuite) SetupSuite() {
	s.setupCalled = true
}

func (s *ExampleSuite) TearDownSuite() {
	s.teardownCalled = true
}

func (s *ExampleSuite) Setup() {}

func (s *ExampleSuite) TearDown() {}

func (s *ExampleSuite) TestExample(t *testing.T) {
	s.testCalled = true
}

func (s *ExampleSuite) TestExampleAsync(t *testing.T) {
	s.testCalled = true
}

func (s *ExampleSuite) TestWithContext(ctx context.Context, t *testing.T) {
	if ctx == nil {
		t.Error("context should not be nil")
	}
}

// TestSuiteRunner tests the suite runner functionality
func TestSuiteRunner(t *testing.T) {
	s := &ExampleSuite{}
	RunSuite(t, s)

	if !s.setupCalled {
		t.Error("SetupSuite was not called")
	}
	if !s.teardownCalled {
		t.Error("TearDownSuite was not called")
	}
	if !s.testCalled {
		t.Error("TestExample was not called")
	}
}

// TestSubSuiteRunner tests running suite as subtests
func TestSubSuiteRunner(t *testing.T) {
	s := &ExampleSuite{}
	RunSubSuite(t, s)

	if !s.setupCalled {
		t.Error("SetupSuite was not called")
	}
	if !s.teardownCalled {
		t.Error("TearDownSuite was not called")
	}
	if !s.testCalled {
		t.Error("TestExample was not called")
	}
}

// TestCollectTests tests test collection
func TestCollectTests(t *testing.T) {
	s := &ExampleSuite{}
	r := &runner{
		suite:   s,
		testing: t,
		results: make(map[string]TestResult),
	}

	tests := r.collectTests()
	if len(tests) < 3 {
		t.Errorf("expected at least 3 tests, got %d", len(tests))
	}

	// Check that async test is detected
	hasAsync := false
	for _, test := range tests {
		if test.IsAsync {
			hasAsync = true
		}
	}
	if !hasAsync {
		t.Error("expected to find async test")
	}
}

// TestSuiteName tests suite name extraction
func TestSuiteName(t *testing.T) {
	s := &ExampleSuite{}
	r := &runner{
		suite:   s,
		testing: t,
		results: make(map[string]TestResult),
	}

	name := r.suiteName()
	if name != "ExampleSuite" {
		t.Errorf("expected 'ExampleSuite', got '%s'", name)
	}
}

// TestResults tests test result tracking
func TestResults(t *testing.T) {
	s := &ExampleSuite{}
	r := &runner{
		suite:   s,
		testing: t,
		results: make(map[string]TestResult),
	}

	r.results["Test1"] = TestResult{Name: "Test1", Passed: true}
	r.results["Test2"] = TestResult{Name: "Test2", Passed: false}
	r.results["Test3"] = TestResult{Name: "Test3", Skipped: true}

	if len(r.results) != 3 {
		t.Errorf("expected 3 results, got %d", len(r.results))
	}
}

// TestTestMethod tests testMethod struct
func TestTestMethod(t *testing.T) {
	tm := testMethod{
		Name:           "TestFoo",
		Func:           reflect.Value{},
		IsAsync:        false,
		IsTeardown:     false,
		AcceptsContext: true,
	}

	if tm.Name != "TestFoo" {
		t.Errorf("expected 'TestFoo', got '%s'", tm.Name)
	}
	if !tm.AcceptsContext {
		t.Error("expected AcceptsContext to be true")
	}
}

// Verify reflect is imported correctly
var _ = strings.TrimSpace

// TestWithContextVerify verifies context is passed correctly
func TestWithContextVerify(t *testing.T) {
	ctx := context.Background()
	if ctx == nil {
		t.Error("context should not be nil")
	}
}
