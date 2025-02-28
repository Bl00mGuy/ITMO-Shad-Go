package languages

import (
	"encoding/json"
	"os"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

var mapOfLanguages = make(map[string][]string)

func loadFileContent(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func parseJSON(content []byte) ([]stats.Language, error) {
	var languages []stats.Language
	err := json.Unmarshal(content, &languages)
	return languages, err
}

func populateLanguageMap(languages []stats.Language) {
	for _, lang := range languages {
		mapOfLanguages[strings.ToLower(lang.Name)] = lang.Extensions
	}
}

func LoadLanguages() {
	if len(mapOfLanguages) > 0 {
		return
	}

	content, err := loadFileContent("../../configs/language_extensions.json")
	if err != nil {
		panic("failed to read language_extensions.json: " + err.Error())
	}

	languages, err := parseJSON(content)
	if err != nil {
		panic("failed to unmarshal JSON content: " + err.Error())
	}

	populateLanguageMap(languages)
}
