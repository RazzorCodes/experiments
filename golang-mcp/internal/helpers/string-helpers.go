package helpers

import "strings"

// maybe switch to aho-corrasik?
// edge case: multiple windows share the same identical content
// issue: doesn't match case-insessitive
func ExtractWindowAroundKeywords(
	text string,
	keywords []string,
	windowSize int) []string {
	var results []string

	for _, kw := range keywords {
		start := 0
		for {
			i := strings.Index(text[start:], kw)
			if i == -1 {
				break
			}
			i += start

			from := i - windowSize
			if from < 0 {
				from = 0
			}

			to := i + len(kw) + windowSize
			if to > len(text) {
				to = len(text)
			}

			results = append(results, text[from:to])

			start = i + 1 // continue search
		}
	}

	return results
}

func firstNLines(text string, n int) string {
	lines := strings.Split(text, "\n")
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, "\n")
}

func truncate(text string, maxChars int) string {
	if maxChars > 0 && len(text) > maxChars {
		return text[:maxChars]
	}
	return text
}
