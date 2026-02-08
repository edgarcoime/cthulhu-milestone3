package strings

const ELLIPSIS = "..."

// Util to show the first n characters of a string and add ... if the string is longer than n characters
func TruncateString(s string, n int) string {
	if len(s) > n && n > 0 {
		return s[:n] + ELLIPSIS
	}
	return s
}
