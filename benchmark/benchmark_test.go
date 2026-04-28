package benchmark

import (
	"math"
	"testing"
	"time"
)

func TestResult(t *testing.T) {
	rb := &testing.B{}
	rb.N = 1000
	rb.ResetTimer()
	rb.StartTimer()
	time.Sleep(100 * time.Microsecond)
	rb.StopTimer()

	result := Result{
		Name:    "TestBenchmark",
		Ops:     rb.N,
		NsPerOp: float64(100*1000) / float64(rb.N),
		Measured: 100 * time.Microsecond,
		Runs:     1,
	}

	if result.Name != "TestBenchmark" {
		t.Errorf("expected Name 'TestBenchmark', got '%s'", result.Name)
	}
	if result.Ops != rb.N {
		t.Errorf("expected Ops %d, got %d", rb.N, result.Ops)
	}
	if result.NsPerOp <= 0 {
		t.Errorf("expected positive NsPerOp, got %f", result.NsPerOp)
	}
}

func TestStats(t *testing.T) {
	// Test with known values
	stats := Stats{
		Mean:   100.0,
		Median: 100.0,
		StdDev: 10.0,
		Min:    90.0,
		Max:    110.0,
	}

	if stats.Mean != 100.0 {
		t.Errorf("expected Mean 100.0, got %f", stats.Mean)
	}
	if stats.Median != 100.0 {
		t.Errorf("expected Median 100.0, got %f", stats.Median)
	}
	if stats.StdDev != 10.0 {
		t.Errorf("expected StdDev 10.0, got %f", stats.StdDev)
	}
}

func TestCompare(t *testing.T) {
	// r1 is baseline (1000ns), r2 is faster (800ns) -> 20% improvement
	result1 := Result{NsPerOp: 1000}
	result2 := Result{NsPerOp: 800}

	improvement := Compare(result1, result2)
	expected := 20.0

	if math.Abs(improvement-expected) > 0.01 {
		t.Errorf("expected improvement %.2f%%, got %.2f%%", expected, improvement)
	}

	// r1 is 800ns, r2 is 1000ns -> r2 is 25% slower (negative improvement)
	result1 = Result{NsPerOp: 800}
	result2 = Result{NsPerOp: 1000}

	improvement = Compare(result1, result2)
	expected = -25.0

	if math.Abs(improvement-expected) > 0.01 {
		t.Errorf("expected improvement %.2f%%, got %.2f%%", expected, improvement)
	}

	// Test with zero baseline
	result1 = Result{NsPerOp: 0}
	result2 = Result{NsPerOp: 800}

	improvement = Compare(result1, result2)
	if improvement != 0 {
		t.Errorf("expected 0 for zero NsPerOp, got %f", improvement)
	}
}

func TestBenchmarkTimer(t *testing.T) {
	timer := NewTimer()
	time.Sleep(10 * time.Millisecond)
	elapsed := timer.Elapsed()

	if elapsed < 10*time.Millisecond {
		t.Errorf("expected at least 10ms elapsed, got %v", elapsed)
	}

	duration := timer.Duration()
	if duration == "" {
		t.Error("expected non-empty duration string")
	}
}

func TestPrecisionTimer(t *testing.T) {
	timer := NewPrecisionTimer()
	time.Sleep(10 * time.Millisecond)

	nanos := timer.ElapsedNanos()
	if nanos < 10*1000*1000 {
		t.Errorf("expected at least 10ms in nanoseconds, got %d", nanos)
	}

	ms := timer.ElapsedMs()
	if ms < 10 {
		t.Errorf("expected at least 10ms, got %f", ms)
	}
}



func TestRunMultiple(t *testing.T) {
	b := &testing.B{}
	counter := 0

	stats := RunMultiple(b, func() {
		time.Sleep(1 * time.Millisecond)
		counter++
	}, 5)

	if counter != 5 {
		t.Errorf("expected 5 runs, got %d", counter)
	}

	if stats.Mean <= 0 {
		t.Errorf("expected positive mean, got %f", stats.Mean)
	}

	if stats.Median <= 0 {
		t.Errorf("expected positive median, got %f", stats.Median)
	}

	if stats.Min <= 0 {
		t.Errorf("expected positive min, got %f", stats.Min)
	}

	if stats.Max <= 0 {
		t.Errorf("expected positive max, got %f", stats.Max)
	}
}

func BenchmarkSimple(b *testing.B) {
	f := func() {
		sum := 0
		for i := 0; i < 100; i++ {
			sum += i
		}
	}
	Run(b, f)
}

func BenchmarkMultipleRuns(b *testing.B) {
	f := func() {
		sum := 0
		for i := 0; i < 1000; i++ {
			sum += i
		}
	}
	RunMultiple(b, f, 3)
}
