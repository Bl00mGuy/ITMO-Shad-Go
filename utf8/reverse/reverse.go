//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	builder := strings.Builder{}
	builder.Grow(len(input))
	for i := len(input); i > 0; {
		r, size := utf8.DecodeLastRuneInString(input[:i])
		if r == utf8.RuneError {
			builder.WriteRune('�')
			i -= 1
		} else {
			builder.WriteRune(r)
			i -= size
		}
	}
	return builder.String()
}
