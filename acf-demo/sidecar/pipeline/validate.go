package pipeline

import (
	"errors"
	"strings"
)

// ValidateInput ensures a request contains non-empty input.
func ValidateInput(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("input cannot be empty")
	}
	return nil
}
