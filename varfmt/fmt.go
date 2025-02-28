//go:build !solution

package varfmt

import (
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	builder := strings.Builder{}
	index := 0
	formatLength := len(format)
	for i := 0; i < formatLength; i++ {
		if format[i] == '{' {
			end := i + 1
			for end < formatLength && format[end] != '}' {
				end++
			}
			if end == formatLength {
				builder.WriteString(format[i:])
				break
			}
			placeholder := format[i+1 : end]
			i = end
			if placeholder == "" {
				if index < len(args) {
					switch v := args[index].(type) {
					case int:
						builder.WriteString(strconv.Itoa(v))
					case string:
						builder.WriteString(v)
					default:
						builder.WriteString("")
					}
				}
			} else {
				if value, err := strconv.Atoi(placeholder); err == nil && value < len(args) {
					switch v := args[value].(type) {
					case int:
						builder.WriteString(strconv.Itoa(v))
					case string:
						builder.WriteString(v)
					default:
						builder.WriteString("")
					}
				}
			}
			index++
		} else {
			builder.WriteByte(format[i])
		}
	}
	return builder.String()
}
