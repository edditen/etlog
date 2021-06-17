package utils

import "strings"

func LastSubstring(s string, sep string) string {
	idx := strings.LastIndex(s, sep)
	return s[idx+1:]
}

func FirstSubstring(s string, sep string) string {
	idx := strings.LastIndex(s, sep)
	if idx < 0 {
		return s
	}
	return s[:idx]
}
