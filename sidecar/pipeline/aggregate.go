package pipeline

import "math"

const signalWeight = 0.4

// AggregateScore converts triggered signals into a bounded risk score.
func AggregateScore(signals []string) float64 {
	score := float64(len(signals)) * signalWeight
	return math.Min(1.0, score)
}
