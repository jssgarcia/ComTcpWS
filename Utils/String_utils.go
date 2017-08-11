package Utils

import (
	"regexp"
	"strings"
	"unicode/utf8"
	"encoding/json"
)

func PrettyPrint(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

// Check if string is empty
func IsEmpty(s string) bool {
	if len(s) == 0 {
		return true
	}
	return false
}

// Check if string is not empty
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// Check if any one of strings are empty
func IsAnyEmpty(ss ...string) bool {
	for _, s := range ss {
		if len(s) == 0 {
			return true
		}
	}
	return false
}

// Checks if none of the strings are empty
func IsNoneEmpty(ss ...string) bool {
	for _, s := range ss {
		if len(s) == 0 {
			return false
		}
	}
	return true
}

// Checks if a string is whitespace, empty ("")
func IsBlank(s string) bool {
	if len(s) == 0 {
		return true
	}
	reg := regexp.MustCompile(`^\s+$`)
	actual := reg.MatchString(s)
	if actual {
		return true
	}
	return false
}

// Checks if a string is not empty (""), not null and not whitespace only.
func IsNotBlank(s string) bool {
	return !IsBlank(s)
}

// Checks if any one of the strings are blank ("") or null and not whitespace only..
func IsAnyBlank(ss ...string) bool {
	for _, s := range ss {
		// regexp cost is probably expensive
		if IsBlank((string)(s)) {
			return true
		}
	}
	return false
}

// Checks if none of the strings are blank ("") or null and whitespace only..
func IsNoneBlank(ss ...string) bool {
	for _, s := range ss {
		if IsBlank((string)(s)) {
			return false
		}
	}
	return true
}

// Removes control characters  from both ends of this String
func Trim(str string) string {
	return strings.Trim(str, " ")
}

func isUpperASCII(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

func isLowerASCII(c byte) bool {
	return 'a' <= c && c <= 'z'
}

func toUpperASCII(c byte) byte {
	if isLowerASCII(c) {
		return c - ('a' - 'A')
	}
	return c
}

func toLowerASCII(c byte) byte {
	if isUpperASCII(c) {
		return c + 'a' - 'A'
	}
	return c
}

func keys(m map[string]struct{}) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

func shortestLength(strs []string, shortest int) int {
	for _, s := range strs {
		if candidate := utf8.RuneCountInString(s); candidate < shortest {
			shortest = candidate
		}
	}
	return shortest
}

func longestLength(strs []string) (longest int) {
	for _, s := range strs {
		if candidate := utf8.RuneCountInString(s); candidate > longest {
			longest = candidate
		}
	}
	return longest
}