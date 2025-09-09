package generator

import (
	"strings"
)

const (
	// The character set for our base62 encoding.
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	base     = int64(len(alphabet))
)

// Encode converts a base-10 integer (our database ID) to a base-62 string.
func Encode(n int64) string {
	if n == 0 {
		return string(alphabet[0])
	}

	var sb strings.Builder
	for n > 0 {
		remainder := n % base
		sb.WriteByte(alphabet[remainder])
		n /= base
	}

	return reverse(sb.String())
}

// reverse is a helper function to reverse a string.
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
