package utils

import "strings"

type StrOption struct {
	IgnoreCase bool
}

// Checks if the string s is in the list of values specified. By default,
// checks in a case insensitive manner
func StrIn(s string, option *StrOption, values ...string) bool {
	for _, value := range values {
		if StrEquals(s, value, option) {
			return true
		}
	}
	return false
}

// Checks if 2 strings are equivalent.
func StrEquals(s1, s2 string, option *StrOption) bool {
	option = getStrOptions(option)
	if option.IgnoreCase {
		return strings.EqualFold(s1, s2)
	}
	return s1 == s2
}

// Sets the default StrOption if option is nil. Otherwise, passes through the string
// option
func getStrOptions(option *StrOption) *StrOption {
	if option == nil {
		return &StrOption{
			IgnoreCase: true,
		}
	}
	return option
}

// Checks if string is empty or just whitespace
func StrIsEmptyOrWhitespace(s string) bool {
	return strings.TrimSpace(s) == ""
}

// StrIsIn reports whether s exists in a list of string
func StrIsIn(s string, list []string) bool {
	for _, ss := range list {
		if ss == s {
			return true
		}
	}

	return false
}

// StrIsInEqualFold reports whether s exists in a list of string, interpreted as UTF-8 strings,
// under Unicode case-folding
func StrIsInEqualFold(s string, list []string) bool {
	for _, ss := range list {
		if strings.EqualFold(s, ss) {
			return true
		}
	}
	return false
}
