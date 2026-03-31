package pipeline

import (
	"encoding/json"
	"os"
	"strings"
)

// Scanner holds normalized patterns and can scan normalized input.
type Scanner struct {
	patterns []string
}

// NewScanner loads patterns from a JSON file.
func NewScanner(patternsFile string) (*Scanner, error) {
	bytes, err := os.ReadFile(patternsFile)
	if err != nil {
		return nil, err
	}

	var patterns []string
	if err := json.Unmarshal(bytes, &patterns); err != nil {
		return nil, err
	}

	normalized := make([]string, 0, len(patterns))
	for _, p := range patterns {
		p = strings.ToLower(strings.TrimSpace(p))
		if p != "" {
			normalized = append(normalized, p)
		}
	}

	return &Scanner{patterns: normalized}, nil
}

// Scan returns all triggered pattern signals in the input.
func (s *Scanner) Scan(normalizedInput string) []string {
	signals := make([]string, 0)
	for _, pattern := range s.patterns {
		if strings.Contains(normalizedInput, pattern) {
			signals = append(signals, pattern)
		}
	}
	return signals
}
