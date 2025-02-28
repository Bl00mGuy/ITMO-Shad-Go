//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	builder := strings.Builder{}
	builder.Grow(len(input))
	isSpace := false
	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])
		i += size
		if r == utf8.RuneError {
			builder.WriteRune('�')
			continue
		}
		if unicode.IsSpace(r) {
			if !isSpace {
				builder.WriteRune(' ')
				isSpace = true
			}
		} else {
			builder.WriteRune(r)
			isSpace = false
		}
	}
	return builder.String()
}
