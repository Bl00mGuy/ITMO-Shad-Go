//go:build !solution

package speller

import "strings"

var (
	ones = []string{"", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	tens = []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
)

func Spell(n int64) string {
	var output strings.Builder

	if n <= 0 {
		if n == 0 {
			return "zero"
		}
		return "minus " + Spell(-n)
	}
	if n >= 1000000000 {
		billion := n / 1000000000
		output.WriteString(Spell(billion) + " billion ")
		n %= 1000000000
	}
	if n >= 1000000 {
		million := n / 1000000
		output.WriteString(Spell(million) + " million ")
		n %= 1000000
	}
	if n >= 1000 {
		thousand := n / 1000
		output.WriteString(Spell(thousand) + " thousand ")
		n %= 1000
	}
	if n >= 100 {
		hundred := n / 100
		output.WriteString(ones[hundred] + " hundred ")
		n %= 100
	}
	if n >= 20 {
		ten := n / 10
		output.WriteString(tens[ten])
		n %= 10
		if n > 0 {
			output.WriteString("-" + ones[n])
		}
	} else if n > 0 {
		output.WriteString(ones[n])
	}

	return strings.TrimSpace(output.String())
}
