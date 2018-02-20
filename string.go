package utils

import "strconv"

// StringToInt64 converts a string to int64, since strconv doesn't provide this straight up.
func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

// IsParseableAsInt64 checks whether a string is parseable as int64.
func IsParseableAsInt64(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return false
	}
	return true
}
