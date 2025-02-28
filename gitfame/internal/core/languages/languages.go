package languages

import (
	"path/filepath"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func GetFileExtension(fileName string) string {
	return filepath.Ext(fileName)
}

func IsSupportedExtension(lang string, fileExtension string) bool {
	langExts := GetExtensions(lang)
	for _, ext := range langExts {
		if ext == fileExtension {
			return true
		}
	}
	return false
}

func LanguageCheck(request stats.RepoFlags, fileName string) bool {
	if len(request.Languages) == 0 {
		return false
	}

	fileExtension := GetFileExtension(fileName)

	for _, lang := range request.Languages {
		if IsSupportedExtension(lang, fileExtension) {
			return false
		}
	}

	return true
}
