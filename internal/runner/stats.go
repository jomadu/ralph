package runner

import (
	"fmt"
	"math"
	"time"
)

// IterationStats tracks iteration timing statistics using Welford's online algorithm.
type IterationStats struct {
	count int
	min   time.Duration
	max   time.Duration
	mean  float64
	m2    float64
}

// NewIterationStats creates a new statistics tracker.
func NewIterationStats() *IterationStats {
	return &IterationStats{
		min: time.Duration(math.MaxInt64),
	}
}

// Add records a new iteration duration using Welford's online algorithm.
func (s *IterationStats) Add(duration time.Duration) {
	s.count++
	
	// Update min/max
	if duration < s.min {
		s.min = duration
	}
	if duration > s.max {
		s.max = duration
	}
	
	// Welford's online algorithm for mean and variance
	x := duration.Seconds()
	delta := x - s.mean
	s.mean += delta / float64(s.count)
	delta2 := x - s.mean
	s.m2 += delta * delta2
}

// Report formats and returns the statistics summary.
func (s *IterationStats) Report() string {
	if s.count == 0 {
		return ""
	}
	
	if s.count == 1 {
		return fmt.Sprintf("INFO: 1 iteration executed; duration: %.3fs", s.min.Seconds())
	}
	
	stddev := math.Sqrt(s.m2 / float64(s.count))
	
	return fmt.Sprintf("INFO: %d iterations executed; min: %.3fs, max: %.3fs, mean: %.3fs, stddev: %.3fs",
		s.count, s.min.Seconds(), s.max.Seconds(), s.mean, stddev)
}
