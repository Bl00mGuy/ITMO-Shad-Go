package filters

import (
	"path/filepath"

	"gitlab.com/slon/shad-go/gitfame/internal/core/stats"
)

func hasExtensions(request stats.RepoFlags) bool {
	return len(request.Extensions) > 0
}

func getFileExtension(fileName string) string {
	return filepath.Ext(fileName)
}

func isExtensionAllowed(extensions []string, fileExtension string) bool {
	for _, ext := range extensions {
		if ext == fileExtension {
			return true
		}
	}
	return false
}

func ExtensionCheck(request stats.RepoFlags, fileName string) bool {
	if !hasExtensions(request) {
		return false
	}

	fileExtension := getFileExtension(fileName)

	return !isExtensionAllowed(request.Extensions, fileExtension)
}
