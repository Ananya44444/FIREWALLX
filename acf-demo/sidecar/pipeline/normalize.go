package pipeline

import (
	"encoding/base64"
	"regexp"
	"strings"
	"unicode/utf8"
)

var base64Regex = regexp.MustCompile(`^[A-Za-z0-9+/=]+$`)

// NormalizeIterative applies normalization repeatedly up to maxRounds
// and exits early when the output no longer changes.
func NormalizeIterative(input string, maxRounds int) string {
	if maxRounds < 1 {
		maxRounds = 1
	}

	current := input
	for i := 0; i < maxRounds; i++ {
		next := NormalizeOnce(current)
		if next == current {
			break
		}
		current = next
	}
	return current
}

// NormalizeOnce performs one normalization pass.
func NormalizeOnce(input string) string {
	cleaned := strings.TrimSpace(removeZeroWidth(input))
	if decoded, ok := tryBase64Decode(cleaned); ok {
		return strings.TrimSpace(removeZeroWidth(decoded))
	}

	return strings.ToLower(cleaned)
}

func removeZeroWidth(input string) string {
	replacer := strings.NewReplacer(
		"\u200b", "", // zero-width space
		"\u200c", "", // zero-width non-joiner
		"\u200d", "", // zero-width joiner
		"\ufeff", "", // byte order mark / zero-width no-break space
	)
	return replacer.Replace(input)
}

func tryBase64Decode(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) < 8 || len(trimmed)%4 != 0 {
		return "", false
	}
	if !base64Regex.MatchString(trimmed) {
		return "", false
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil || !utf8.Valid(decodedBytes) {
		return "", false
	}

	decoded := strings.TrimSpace(string(decodedBytes))
	if decoded == "" {
		return "", false
	}

	return decoded, true
}
