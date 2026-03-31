package pipeline

// Decision is the output action from the sidecar.
type Decision string

const (
	DecisionAllow    Decision = "ALLOW"
	DecisionBlock    Decision = "BLOCK"
	DecisionSanitize Decision = "SANITIZE"
)

// Decide chooses an action based on score thresholds.
func Decide(score float64) Decision {
	if score >= 0.8 {
		return DecisionBlock
	}
	if score >= 0.4 {
		return DecisionSanitize
	}
	return DecisionAllow
}
