package utils

import "regexp"

func GetFirstMatchedString(re *regexp.Regexp, source string) (string, bool) {
	match := re.FindStringSubmatch(source)
	if len(match) == 0 {
		return "", false
	}
	return match[1], true
}
