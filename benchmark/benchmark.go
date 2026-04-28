package benchmark

import (
	"math"
	"testing"
	"time"
)

// Result holds benchmark execution results
type Result struct {
	Name      string
	Ops       int
	NsPerOp   float64
	MBPerSec  float64
	AllocBytesPerOp uint64
	AllocsPerOp     int
	Measured time.Duration
	Runs     int
}

// Stats calculates statistical metrics from multiple benchmark runs
type Stats struct {
	Mean   float64
	Median float64
	StdDev float64
	Min    float64
	Max    float64
}

// Run executes a benchmark function and collects metrics
func Run(b *testing.B, f func()) Result {
	b.ResetTimer()
	b.StartTimer()

	start := time.Now()
	f()
	elapsed := time.Since(start)

	b.StopTimer()

	return Result{
		Name:      b.Name(),
		Ops:       b.N,
		NsPerOp:   float64(elapsed.Nanoseconds()) / float64(b.N),
		Measured:  elapsed,
		Runs:      1,
	}
}

// RunMultiple executes a benchmark multiple times and returns statistics
func RunMultiple(b *testing.B, f func(), runs int) Stats {
	results := make([]float64, runs)

	for i := 0; i < runs; i++ {
		rb := &testing.B{}
		rb.ResetTimer()
		rb.StartTimer()

		start := time.Now()
		f()
		elapsed := time.Since(start)

		rb.StopTimer()
		results[i] = float64(elapsed.Nanoseconds()) / float64(b.N)
	}

	sum := 0.0
	min := results[0]
	max := results[0]
	for _, r := range results {
		sum += r
		if r < min {
			min = r
		}
		if r > max {
			max = r
		}
	}
	mean := sum / float64(runs)

	// Calculate standard deviation
	variance := 0.0
	for _, r := range results {
		variance += (r - mean) * (r - mean)
	}
	variance /= float64(runs)
	stdDev := math.Sqrt(variance)

	// Calculate median
	sorted := make([]float64, runs)
	copy(sorted, results)
	for i := 0; i < runs-1; i++ {
		for j := i + 1; j < runs; j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	median := sorted[runs/2]
	if runs%2 == 0 {
		median = (sorted[runs/2-1] + sorted[runs/2]) / 2
	}

	return Stats{
		Mean:   mean,
		Median: median,
		StdDev: stdDev,
		Min:    min,
		Max:    max,
	}
}

// Compare compares two benchmark results
func Compare(r1, r2 Result) float64 {
	// Returns percentage improvement: positive means r2 is faster
	if r1.NsPerOp == 0 {
		return 0
	}
	return ((r1.NsPerOp - r2.NsPerOp) / r1.NsPerOp) * 100
}

// BenchmarkTimer is a helper for manual timing
type BenchmarkTimer struct {
	start time.Time
}

// NewTimer creates a new benchmark timer
func NewTimer() *BenchmarkTimer {
	return &BenchmarkTimer{start: time.Now()}
}

// Elapsed returns the elapsed time
func (t *BenchmarkTimer) Elapsed() time.Duration {
	return time.Since(t.start)
}

// Duration returns a formatted duration string
func (t *BenchmarkTimer) Duration() string {
	return t.Elapsed().String()
}

// PrecisionTimer provides high-precision timing using time.Since
type PrecisionTimer struct {
	start int64
}

// NewPrecisionTimer creates a high-precision timer
func NewPrecisionTimer() *PrecisionTimer {
	return &PrecisionTimer{start: time.Now().UnixNano()}
}

// ElapsedNanos returns elapsed time in nanoseconds
func (t *PrecisionTimer) ElapsedNanos() int64 {
	return time.Now().UnixNano() - t.start
}

// ElapsedMs returns elapsed time in milliseconds
func (t *PrecisionTimer) ElapsedMs() float64 {
	return float64(t.ElapsedNanos()) / 1e6
}

// AutoRun automatically runs a benchmark with varying N values
func AutoRun(b *testing.B, f func(), minN, maxN int64) {
	if maxN < minN {
		maxN = minN * 100
	}

	for n := minN; n <= maxN; n *= 10 {
		b.Run("", func(b *testing.B) {
			b.N = int(n)
			Run(b, f)
		})
	}
}
