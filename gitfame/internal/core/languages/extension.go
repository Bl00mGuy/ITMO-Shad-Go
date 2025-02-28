package languages

import (
	"strings"
)

func GetExtensions(lang string) []string {
	LoadLanguages()
	return mapOfLanguages[strings.ToLower(lang)]
}
