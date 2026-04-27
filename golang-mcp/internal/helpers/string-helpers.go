package helpers

import "strings"

func ExtractWindowAroundKeywords(text string, keywords []string, windowSize int) []string {
	lower := strings.ToLower(text)
	seen := make(map[string]struct{})
	var results []string

	for _, kw := range keywords {
		kwLower := strings.ToLower(kw)
		start := 0
		for {
			i := strings.Index(lower[start:], kwLower)
			if i == -1 {
				break
			}
			i += start

			from := i - windowSize
			if from < 0 {
				from = 0
			}
			to := i + len(kwLower) + windowSize
			if to > len(text) {
				to = len(text)
			}

			window := text[from:to]
			if _, dup := seen[window]; !dup {
				seen[window] = struct{}{}
				results = append(results, window)
			}

			start = i + 1
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
